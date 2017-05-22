package gameroom

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"tribe/baseob"
	"tribe/gameroom/event"
	"tribe/gameroom/item"
	"tribe/gameroom/npc"
	"tribe/sqldb"
	"vava6/vatools"
)

var OBManageZone = &ManageZone{
	mpZones: make(map[int]*Zone, 1000),
	lk:      new(sync.RWMutex),
}

type ManageZone struct {
	mpZones map[int]*Zone
	lk      *sync.RWMutex
}

func (this *ManageZone) GetCanchZone(zoneId int) (*Zone, error) {
	this.lk.RLock()
	obZone, ok := this.mpZones[zoneId]
	this.lk.RUnlock()
	// 读取到返回正确的obZone对象
	if ok {
		return obZone, nil
	} else {
		// 没有从数据库里获取
		this.lk.Lock()
		var err error
		obZone, ok = this.mpZones[zoneId]
		if !ok {
			obZone, err = NewZone(zoneId)
			if err == nil {
				this.mpZones[zoneId] = obZone
			}
		}
		this.lk.Unlock()
		return obZone, err
	}
}

// 加载初始化所有的Zone信息
func (this *ManageZone) InitZones() {
	// 通过DB读取所有信息
	rss, err := sqldb.Querys("d_zone", "zoneID", "")
	if err != nil {
		fmt.Println("加载ZONE出错：", err.Error())
		return
	}
	// 遍历
	for _, rs := range rss {
		_, _ = this.GetCanchZone(vatools.SInt(rs["zoneID"]))
	}
	fmt.Println("加载ZONE完成", len(this.mpZones))
}

func (this *ManageZone) Name() string {
	return "zone"
}

func (this *ManageZone) Max() int {
	return 1000
}

func (this *ManageZone) Count() int {
	return len(this.mpZones)
}

func (this *ManageZone) Release() {
	return
}

func (this *ManageZone) Save() error {
	return nil
}

// Zone区域
//	玩家探索的区域
type Zone struct {
	baseob.BaseOB
	zoneID        int                      // zoneID
	zoneName      string                   // zone名称
	zoneLevel     int                      // 等级
	zoneType      int8                     // zone类型 0-PVE 1-PVP
	zoneArea      int                      // zone面积
	zoneFood      int                      // zone消耗的食物的基数
	zoneImg       string                   // 图标
	zoneMsg       string                   // 描述
	zoneInfoItems []*item.Item             // 这个区域可以掉落的物品信息
	zoneInfoNpcs  []*npc.Monster           // 这个区域可以被发现的怪物
	zoneEvents    []*event.Event           // 事件列表共公事件
	mpExpque      map[int]bool             // 当前Zone探索中的探索队伍
	mpField       map[int]*FieldZone       // 格子对象
	lk            *sync.RWMutex            // mpExpque锁
	chExpQue      chan *ExploreQueue       // 添加到当前探索队列
	tmpHowField   uint16                   // 临时存放的可以生成的格子数量
	tmpInfoItems  []map[string]interface{} // 临时存放的物品信息
	tmpInfoNpcs   []map[string]interface{} // 临时存放的NPC信息
}

func NewZone(zoneId int) (*Zone, error) {
	ob := &Zone{
		zoneID:        zoneId,
		lk:            new(sync.RWMutex),
		mpField:       make(map[int]*FieldZone, 10),
		chExpQue:      make(chan *ExploreQueue, 1000),
		zoneInfoItems: make([]*item.Item, 0, 4),
		zoneInfoNpcs:  make([]*npc.Monster, 0, 4),
	}
	rs, err := ob.LoadDB("d_zone", "*", "zoneID", map[string]interface{}{"zoneID": zoneId})
	if err != nil {
		return nil, err
	}
	ob.zoneName = rs["zoneName"]
	ob.zoneLevel = vatools.SInt(rs["zoneLevel"])
	ob.zoneImg = rs["zoneImg"]
	ob.zoneMsg = rs["zoneMsg"]
	ob.zoneFood = vatools.SInt(rs["zoneFood"])
	ob.mpExpque = make(map[int]bool, 1000)
	// 加载地形
	ob._initField(rs["fieldTypeHow"])
	// 初始化可以控索到物品的对象
	ob._initInfoItems(rs["zoneInfoItems"])
	ob._initInfoNpc(rs["zoneInfoNpc"])
	// 初始华探索队列
	ob.initExploreQueue()
	fmt.Println("ZONE_DB加载", ob.zoneName, " 当前共有", len(ob.mpExpque), "探索队")
	// 协程执行探索通道监听
	go func() {
		for {
			ptQue := <-ob.chExpQue
			ob.doExplore(ptQue)
		}
	}()
	go func() {
		for {
			// 时钟10s执行一次
			time.Sleep(time.Second * time.Duration(10))
			// 执行一次运算
			ob.doAllExplore()
		}
	}()
	return ob, nil
}

func (this *Zone) _initInfoItems(str string) {
	arrs := strings.Split(str, ",")
	for _, v := range arrs {
		itemID := vatools.SInt(v)
		ptItem, err := item.OBManageItem.GetCanchItem(itemID)
		if err != nil {
			continue
		}
		this.zoneInfoItems = append(this.zoneInfoItems, ptItem)
	}
}

func (this *Zone) _initInfoNpc(str string) {
	arrs := strings.Split(str, ",")
	for _, v := range arrs {
		npcID := vatools.SInt(v)
		ptMonster, err := npc.NewMonster(npcID)
		if err != nil {
			continue
		}
		this.zoneInfoNpcs = append(this.zoneInfoNpcs, ptMonster)
	}
}

// 初始化Zone里的格子数
func (this *Zone) _initField(str string) {
	arrs := strings.Split(str, ",")
	for _, v := range arrs {
		arr := strings.Split(v, "=")
		if len(arr) != 2 {
			continue
		}
		// 获取FieldDB ID
		fieldID := vatools.SInt(arr[0])
		// 格子的数量
		arrHow := strings.Split(arr[1], "-")
		minHow := vatools.SUint16(arrHow[0])
		maxHow := minHow
		if len(arrHow) == 2 {
			maxHow = vatools.SUint16(arrHow[1])
		}
		if tmpFieldZone, err := NewFieldZone(fieldID, minHow, maxHow); err == nil {
			if _, ok := this.mpField[fieldID]; !ok {
				this.mpField[fieldID] = tmpFieldZone
			}
		}
	}
}

func (this *Zone) _getTmpInfoItems() []map[string]interface{} {
	if this.tmpInfoItems == nil {
		this.tmpInfoItems = make([]map[string]interface{}, len(this.zoneInfoItems))
		for k, v := range this.zoneInfoItems {
			this.tmpInfoItems[k] = v.GetFieldInfo()
		}
	}
	return this.tmpInfoItems
}

func (this *Zone) _getTmpInfoNpcs() []map[string]interface{} {
	if this.tmpInfoNpcs == nil {
		this.tmpInfoNpcs = make([]map[string]interface{}, len(this.zoneInfoNpcs))
		for k, v := range this.zoneInfoNpcs {
			this.tmpInfoNpcs[k] = v.GetFieldInfo()
		}
	}
	return this.tmpInfoNpcs
}

// 获取Zone名称
//	@return
//		string Zone名称
func (this *Zone) GetName() string {
	return this.zoneName
}

// 获取ZoneID
//	@return
//		int ZoneID
func (this *Zone) GetZoneID() int {
	return this.zoneID
}

// 获取Zone的描述
//	@return
//		string Zone的描述信息
func (this *Zone) GetZoneMsg() string {
	return this.zoneMsg
}

// 获取Zone的map类型的描述信息，方便调用于json输出
//	@return
//		map[string]interface{} Zone的详细信息
func (this *Zone) GetInfo() map[string]interface{} {
	res := make(map[string]interface{})
	res["zoneId"] = this.zoneID
	res["zoneName"] = this.zoneName
	res["zoneMsg"] = this.zoneMsg
	res["zoneFieldHow"] = this.GetHowField()
	res["zoneLevel"] = this.zoneLevel
	res["zoneFood"] = this.zoneFood
	res["zoneImg"] = this.zoneImg
	res["zoneInfoItems"] = this._getTmpInfoItems()
	res["zoneInfoNpcs"] = this._getTmpInfoNpcs()
	return res
}

// 加入探索队伍
//	@parames
//		ptQue	*ExploreQueue
//	@return
//		error
func (this *Zone) JoinExploreQueue(ptQue *ExploreQueue) error {
	if ptQue.food < 1 {
		return errors.New("食物不够")
	}
	this.lk.Lock()
	var err error
	if _, ok := this.mpExpque[ptQue.id]; ok {
		err = errors.New("已存在这个探索队列")
	} else {
		this.mpExpque[ptQue.id] = true
	}
	this.lk.Unlock()
	if err == nil {
		ptQue.ChangeState(1)
		ptQue.JoinZone(this.zoneID)
		// 分配Field给指定的ExploreQueue
		ptQue.zoneFields = this.CreatePlayerField()
		ptQue.zoneFields.playerID = ptQue.leadID
	}
	return err
}

// 初始化这个Zone里所有的探索队
func (this *Zone) initExploreQueue() {
	// 通DB读取所有的属于这个区域的探索队
	rss, err := sqldb.Querys("u_explorequeue", "id", fmt.Sprint("zoneID=", this.zoneID, " AND state=1"))
	if err != nil {
		return
	}
	// 保存到现有队列
	for _, rs := range rss {
		this.mpExpque[vatools.SInt(rs["id"])] = true
	}
}

func (this *Zone) doAllExplore() {
	this.lk.RLock()
	for queID, _ := range this.mpExpque {
		ptQue, err := OBManageExplore.GetExplore(queID)
		if err != nil {
			// 读取指定探索队伍发生错误
			//	...
			continue
		}
		// 不使用Go程运行
		// this.doExplore(ptQue)
		// 使用Go程运行
		go this.doExplore(ptQue)
	}
	this.lk.RUnlock()
}

// 执行探索
//	@parames
//		ptQue	*ExploreQueue	要被执行玩家探索队
func (this *Zone) doExplore(ptQue *ExploreQueue) {
	// 让探索队执行探索
	ptQue.DoExplore()
}

// 获取当前Zone里指定的DBField
func (this *Zone) getFieldDb(id int) (*FieldZone, error) {
	dbField, ok := this.mpField[id]
	if !ok {
		return nil, errors.New("NO FieldDB")
	}
	return dbField, nil
}

// 根据Zone地形生成玩家的场景地图
//	@return
//		*FieldPlayer
func (this *Zone) CreatePlayerField() *FieldPlayer {
	// 获得场景块数
	howField := this.GetHowField()
	// 获得二维
	var x, y int16 = 0, 0
	switch {
	case howField <= 20:
		x, y = 4, 6
	case howField <= 50:
		x, y = 6, 9
	case howField <= 100:
		x, y = 10, 15
	default:
		x, y = 20, 30
	}
	playerField := NewFieldPlayer(this.zoneID)
	playerField.zoneId = this.zoneID
	playerField.maxX = x
	playerField.maxY = y
	// 获得地形数量
	tMap := this.getCreateField()
	// 比例分布器
	tPro := make(map[string]int, len(tMap))
	for k, v := range tMap {
		tPro[strconv.Itoa(int(k))] = int(v)
	}
	ptPro := vatools.NewCanGetProportion(1000, tPro)
	for {
		k, err := ptPro.GetValue()
		if err != nil {
			break
		}
		// 分配格子，如果以后还可以这上面基础分配其它事件
		// TODO...
		field, ok := this.mpField[vatools.SInt(k)]
		if !ok {
			continue
		}
		// 生成格子对象
		tField := NewField(field.GetID())
		// 生成图片
		tField.img = field.GetFieldImage()
		// 添加到玩家的格子管理对象
		playerField.AddField(tField)
	}
	return playerField
}

func (this *Zone) getCreateField() map[int]uint16 {
	res := make(map[int]uint16, len(this.mpField))
	for _, v := range this.mpField {
		res[v.GetID()] = v.GetCreateHow()
	}
	return res
}

// 获取当前Zone拥有的Field数量
func (this *Zone) GetHowField() uint16 {
	if this.tmpHowField < 1 {
		// 计算数量
		for _, ptField := range this.mpField {
			this.tmpHowField += ptField.maxHow
		}
	}
	return this.tmpHowField
}
