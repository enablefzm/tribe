package gameroom

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	"tribe/baseob"
	"tribe/gameroom/constvalue"
	"tribe/gameroom/event"
	"tribe/gameroom/exploreaction"
	"tribe/gameroom/hero"
	"tribe/gameroom/item"
	"tribe/sqldb"
	"vava6/vatools"
)

var OBManageExplore = &ManageExplore{
	lk:        new(sync.RWMutex),
	mpExplore: make(map[int]*ExploreQueue, 1000),
}

type ManageExplore struct {
	lk        *sync.RWMutex
	mpExplore map[int]*ExploreQueue
}

func (this *ManageExplore) GetExplore(id int) (*ExploreQueue, error) {
	this.lk.RLock()
	obExplore, ok := this.mpExplore[id]
	this.lk.RUnlock()
	if ok {
		return obExplore, nil
	}
	// 没有
	this.lk.Lock()
	var err error
	obExplore, ok = this.mpExplore[id]
	if !ok {
		// 从数据库里加载
		obExplore, err = NewExplore(id)
		if err == nil {
			this.mpExplore[id] = obExplore
		}
	}
	this.lk.Unlock()
	return obExplore, err
}

// 注册当前没有在管理对象池里的探索队列
//	@parames
//		*ExploreQueue
//	@return
//		*ExploreQueue
//		error
//			1 代表要被注册对象ID为0是个未被保存的对象
//			2 当前管理对象池里有这个ID存在
func (this *ManageExplore) RegExplore(obQue *ExploreQueue) (*ExploreQueue, error) {
	if obQue.id < 1 {
		return nil, errors.New("1")
	}
	this.lk.Lock()
	defer this.lk.Unlock()
	_, ok := this.mpExplore[obQue.id]
	if ok {
		return nil, errors.New("2")
	}
	this.mpExplore[obQue.id] = obQue
	return obQue, nil
}

func (this *ManageExplore) Name() string {
	return "explorequeue"
}

func (this *ManageExplore) Max() int {
	return 1000
}

func (this *ManageExplore) Count() int {
	return len(this.mpExplore)
}

// Cache接口实例
//		保存当前缓存并清空
func (this *ManageExplore) Release() {
	total := 0
	for k, ptQue := range this.mpExplore {
		// 只有队列状态为0空闲时才会被清出内存
		if ptQue.state == 0 {
			ptQue.Save()
			delete(this.mpExplore, k)
			total++
		} else {
			ptQue.Save()
		}
	}
	fmt.Println("释放出", total, "个探索队列")
	return
}

func (this *ManageExplore) Save() error {
	return nil
}

// 玩家探索队列组管理
type ExploreQueues struct {
	leadID int
	arrIDS []int
	maxQue uint8
}

// 探索队信息
func NewExploreQueues(leadID int) *ExploreQueues {
	obQueues := &ExploreQueues{
		leadID: leadID,
		arrIDS: make([]int, 0, 2),
		maxQue: 1,
	}
	// 加载数据库里的信息
	rss, err := sqldb.Querys("u_explorequeue", "id, leadID", fmt.Sprint("leadID=", leadID))
	if err == nil {
		for _, rs := range rss {
			obQueues.AddExploreID(vatools.SInt(rs["id"]))
		}
	}
	return obQueues
}

func (this *ExploreQueues) AddExploreID(id int) {
	this.arrIDS = append(this.arrIDS, id)
}

func (this *ExploreQueues) GetExploreQueues() []*ExploreQueue {
	arrLen := len(this.arrIDS)
	res := make([]*ExploreQueue, 0, arrLen)
	for i := 0; i < arrLen; i++ {
		obQueue, err := OBManageExplore.GetExplore(this.arrIDS[i])
		if err != nil {
			continue
		}
		res = append(res, obQueue)
	}
	return res
}

// 获取默认第一条探队列，如果当前探索队列为空时则创建一个
//	@return
//		*ExploreQueue
func (this *ExploreQueues) GetDefaultExploreQueue() (*ExploreQueue, error) {
	if len(this.arrIDS) < 1 {
		obQue, err := this.CreateNewExplore()
		if err != nil {
			return nil, err
		}
		return obQue, nil
	} else {
		obQue, err := OBManageExplore.GetExplore(this.arrIDS[0])
		if err != nil {
			return nil, err
		}
		return obQue, err
	}
}

// 创建新的探索队
//	@return
//		*ExploreQueue, error
func (this *ExploreQueues) CreateNewExplore() (*ExploreQueue, error) {
	arrLen := uint8(len(this.arrIDS))
	if arrLen >= this.maxQue {
		return nil, errors.New("1")
	}
	obQue := &ExploreQueue{
		leadID:     this.leadID,
		startTime:  time.Now().Unix(),
		food:       1000,
		zoneFields: NewFieldPlayer(0),
		arrLogs:    make([]string, 0, 50),
	}
	// 初始化已放现道具列表
	obQue.initItems("")
	obQue.SetDBInfo("u_explorequeue", "*", "id", map[string]interface{}{})
	obQue.SetNew(true)
	obQue.Save()
	// 加入当前管理队列组
	this.AddExploreID(obQue.id)
	// 注册到管理对象里
	_, _ = OBManageExplore.RegExplore(obQue)
	return obQue, nil
}

// 通过DB里构造玩家数据信息
//	@parames
//		id		探索队ID
//	@return
//		*ExploreQueue
//		error
func NewExplore(id int) (*ExploreQueue, error) {
	rss, err := sqldb.Querys("u_explorequeue", "*", fmt.Sprint("id=", id))
	if err != nil {
		return nil, err
	}
	if len(rss) < 1 {
		return nil, errors.New(baseob.ERR_NULL)
	}
	ob := NewExploreQueueOnRS(rss[0])
	return ob, nil
}

// 通过DB->RS构造玩家ExploreQueue
//	@parames
//		rs	map[string]string
//	@return
//		*ExploreQueue
func NewExploreQueueOnRS(rs map[string]string) *ExploreQueue {
	ob := &ExploreQueue{
		actHunt:       exploreaction.NewHunt(1000),
		actTreasure:   exploreaction.NewTreasure(100),
		actCollection: exploreaction.NewCollection(20),
		actInsight:    exploreaction.NewInsight(500),
		actRest:       exploreaction.NewRest(0),
		expPower:      NewExplorePower(),
		arrLogs:       make([]string, 0, 50),
	}
	ob.id = vatools.SInt(rs["id"])
	ob.leadID = vatools.SInt(rs["leadID"])
	ob.food = vatools.SInt(rs["food"])
	ob.zoneID = vatools.SInt(rs["zoneID"])
	ob.startTime = vatools.STime(rs["startTime"]).Unix()
	ob.findHow = vatools.SUint(rs["findHow"])
	// 构造已查找到的物品信息
	ob.initItems(rs["items"])
	// 构造玩家地图
	ob.initZoneFields(rs["zoneFields"])
	// 设定DB信息
	ob.SetDBInfo("u_explorequeue", "*", "id", map[string]interface{}{"id": ob.id, "leadID": ob.leadID})
	ob.SetNew(false)
	return ob
}

// 玩家探索队列
type ExploreQueue struct {
	baseob.BaseOB
	id            int
	leadID        int
	zoneID        int
	food          int
	findHow       uint                   // 已探索次数
	physical      uint16                 // 体能值
	state         uint8                  // 探索队状态（0-空闲、1-正在探索中）
	eState        uint16                 // 探索中的状态（0：休息、1：行走、2：探索，3：采集、4：挖矿，）
	startTime     int64                  // 探索开始时间
	maxHero       uint8                  // 队列英雄上限数值
	exploreVal    uint                   // 探索值
	arrItems      []*item.Items          // 掉落物品信息
	arrHeroID     []int                  // 队列玩家所属的英雄ID
	zoneFields    *FieldPlayer           // 玩家探索队列的地图Field对象
	expHunt       uint                   // 探索队守猎能力值
	expTreasure   uint                   // 探索队寻宝能力值
	expCollection uint                   // 探索队探索能力值
	expInsight    uint                   // 探索队洞察能力值
	arrLogs       []string               // 玩家的日列表
	actRest       exploreaction.IFAction // 休息
	actHunt       exploreaction.IFAction // 狩猎动作
	actTreasure   exploreaction.IFAction // 寻宝动作
	actCollection exploreaction.IFAction // 采集动作
	actInsight    exploreaction.IFAction // 洞察动作
	expPower      *ExplorePower          // 畜力值
	ptProportion  *vatools.Proportion    // 占比对象
	nowAct        exploreaction.IFAction // 当前正在操作的动作
}

func (this *ExploreQueue) AddItems(items *item.Items) {
	if items.Superposition() {
		for _, ob := range this.arrItems {
			if ob.ItemID() == items.ItemID() {
				ob.OperateHow(items.GetHow())
				return
			}
		}
	}
	this.arrItems = append(this.arrItems, items)
}

// 初始化已查询到商品信息
//	格式： itemID:how,itemID:how
func (this *ExploreQueue) initItems(strJson string) {
	this.arrItems = make([]*item.Items, 0, 20)
	// 反序列JSON
	arrmp := make([]map[string]int, 0, 20)
	err := json.Unmarshal([]byte(strJson), &arrmp)
	if err != nil {
		fmt.Println("反序列探索队物品的JSON出错:", err.Error())
		return
	}
	il := len(arrmp)
	for i := 0; i < il; i++ {
		obItems, err := item.NewMapItems(arrmp[i])
		if err != nil {
			fmt.Println("实例化探索队控索到物品出错:", err.Error())
			fmt.Println(arrmp[i])
			continue
		}
		// 添加到当前数据库
		this.arrItems = append(this.arrItems, obItems)
		// 加载成功
		fmt.Println("加载成功 -", obItems.ItemName())
	}
}

// 将json的地图信息转换为地图对象
func (this *ExploreQueue) initZoneFields(strJson string) {
	// 转换为对象
	this.zoneFields = NewFieldPlayerOnJson(strJson)
}

// 记录玩家日志
func (this *ExploreQueue) Log(logInfo ...interface{}) {
	logStr := fmt.Sprint(logInfo)
	// 清除之前的日志
	if len(this.arrLogs) >= 50 {
		this.arrLogs = append(this.arrLogs[1:], logStr)
	} else {
		this.arrLogs = append(this.arrLogs, logStr)
	}
	// LOG
	fmt.Println(logStr)
}

// 探索队执行探索
func (this *ExploreQueue) DoExplore() {
	// 先判断本身是否在探索状态
	if this.CheckIsQuit() {
		return
	}
	// 获取对象所属的Zone
	ptZone, err := OBManageZone.GetCanchZone(this.zoneID)
	// 获取不到Zone对象，设定对象为退出状态
	if err != nil {
		this.SetStateQuit()
		return
	}
	// 获取相应的格子，没有找到相应的Field格子对象退出
	nowField, err := this.zoneFields.GetNowFieldErr()
	if err != nil {
		this.SetStateQuit()
		return
	}
	ptFieldZone, err := ptZone.getFieldDb(nowField.fieldID)
	if err != nil {
		this.SetStateQuit()
		return
	}
	// 有动作完成上一次动作
	if this.nowAct != nil {
		// 有动作则执行结果
		et, ok := ptFieldZone.GetEvent(this.nowAct.GetActTypeOnID())
		if !ok {
			this.doMoveNextField()
			return
		}
		rndVal := this.getFightEventRndValue(this.nowAct, et)
		if rndVal <= et.GetProbability() {
			obtain, err := et.GetProbability()
			if err != nil {
				this.Log(queueAct.GetActName(), "时", et.GetName(), "但是什么都没有得到")
			} else {
				ptQue.Log(fmt.Sprint(queueAct.GetActName(), "时", et.GetName(), "获得", obtain.GetInfo()))
			}
		} else {
			if ok := this.GetExpPower().AddVal(this.nowAct.GetActTypeOnID()); !ok {
				this.doMoveNextField()
			}
		}
	}

	// 获取要被消耗的食物计算方法，未定
	//	TODO...
	//	needFood := ptZone.zoneFood
	//	if this.food < needFood {
	//		this.SetStateQuit()
	//		return
	//	}
	//	// 没有动作选择动作
	//	var ok bool
	//	if this.nowAct, ok = this.GetEventOnFieldDb(ptFieldZone.FieldDb); !ok {
	//		this.doMoveNextField()
	//	}
}

// 通过当前地形获得动作
func (this *ExploreQueue) CreateAct() {
	if this.CheckIsQuit() {
		return
	}
	ptZone, err := OBManageZone.GetCanchZone(this.zoneID)
	if err != nil {
		this.SetStateQuit()
		return
	}
	// Debug
	nowField, err := this.zoneFields.GetNowFieldErr()

}

// 动作能力值和事件防守值对抗，对抗后获得事件对象随机值，随机值决定是否可以获得事件奖励
//	@parames
//		exploreaction.IFAction		// 动作对象
//		*event.Event				// 事件对象
//	@return
//		uint16						// 对抗结果值用于是否可以获得事件奖励
func (this *ExploreQueue) getFightEventRndValue(queueAct exploreaction.IFAction, et *event.Event) uint16 {
	blnAdd := true
	fightValue := queueAct.GetActValue() - et.GetDefense()
	if fightValue < 0 {
		fightValue *= -1
		blnAdd = false
	}
	// 获取可以获得增加成功率的值
	addRnd := uint16(float64(fightValue) / float64(et.GetDefense()) * float64(et.GetProbability()))
	maxAddRnd := et.GetProbability() * 3
	if addRnd > maxAddRnd {
		addRnd = maxAddRnd
	}
	// 获取随机值
	rndVal := uint16(vatools.CRnd(1, 1000))
	if blnAdd == true {
		if rndVal > addRnd {
			rndVal -= addRnd
		} else {
			rndVal = 0
		}
	} else {
		rndVal += addRnd
	}
	return rndVal
}

// 探索队移动到下一格如果下一格不存则标记为退出
//	被标记为退出的探索队会被Zone的清理对象移出内存
func (this *ExploreQueue) doMoveNextField() {
	if err := this.MoveNextField(); err != nil {
		this.SetStateQuit()
	}
}

// 玩家探索队在当前Zone里移动到下一格
func (this *ExploreQueue) MoveNextField() error {
	this.expPower.Reset()
	nextField, err := this.zoneFields.NextField()
	if err != nil {
		return err
	}
	nowField := this.zoneFields.GetNowField()
	nowField.after++
	this.zoneFields.nowField = nextField
	this.Log("移动到下一格：", this.zoneFields.nowField.fieldID, this.zoneFields.nowField.pointX, this.zoneFields.nowField.pointY)
	return nil
}

// 设定探索队离开这个队列
func (this *ExploreQueue) SetStateQuit() {
	this.state = 0
}

// 判断是不是已经不在探索状态
func (this *ExploreQueue) CheckIsQuit() bool {
	if this.state == 0 {
		return true
	} else {
		return false
	}
}

// 获取当前随机比例权重对象
func (this *ExploreQueue) GetProportion() *vatools.Proportion {
	if this.ptProportion == nil {
		mp := make(map[string]int, 4)
		mp["E_HUNT"] = int(this.expHunt)
		mp["E_TREASURE"] = int(this.expTreasure)
		mp["E_COLLECTION"] = int(this.expCollection)
		mp["E_INSIGHT"] = int(this.expInsight)
		this.ptProportion = vatools.NewProportion(1000, mp)
	}
	return this.ptProportion
}

// 加入新的英雄
//	@parames
//		*Hero	英雄对象
//	@return
//		error	是否成功
func (this *ExploreQueue) JoinHero(obHero *hero.Hero) error {
	if uint8(len(this.arrHeroID)) >= this.maxHero {
		return errors.New("1")
	}
	this.arrHeroID = append(this.arrHeroID, obHero.GetID())
	// 改变探索值
	// Todo....
	return nil
}

// 获取有多少个英雄信息
// 	@return
//		int 	探索队里的英雄数量
func (this *ExploreQueue) LenHero() int {
	return len(this.arrHeroID)
}

// 获取可以执行的动作列表
func (this *ExploreQueue) GetEventOnFieldDb(ptField *FieldDb) (exploreaction.IFAction, bool) {
	if this.expPower.PowerValue > 0 {
		// 有记录上一次的动作暴发直接返回上一次的动作
		fmt.Println("返回上一次动作:", this.expPower.PowerType)
		return this.IdxAction(this.expPower.PowerType), true
	} else {
		acts := ptField.GetActs()
		ilAct := len(acts)
		if ilAct < 1 {
			return nil, false
		}
		mp := make(map[string]int, ilAct)
		for _, act := range acts {
			// 获得键值的字符串Key 并获得相对应的动作能力值
			mp[constvalue.GetActionKey(act)] = this.IdxAction(act).GetActValue()
		}
		// 生成分配器
		pPro := vatools.NewProportion(1000, mp)
		// 通过分配器返回概率性的动作Key，能力值越高的动作越高机率返回
		actKey := pPro.GetRndKey()
		// 通过动作字符串的Key返回相应的动作
		fmt.Println("获取新的动作：", this.KeyAction(actKey).GetActTypeOnID())
		return this.KeyAction(actKey), true
	}
}

// 获取探索队列的LeadID
// 	@return
//		int		队列所属的角色ID
func (this *ExploreQueue) GetLeadID() int {
	return this.leadID
}

func (this *ExploreQueue) GetInfo() string {
	return fmt.Sprint("ID:", this.id, " LeadID:", this.leadID, " Time:", this.startTime)
}

func (this *ExploreQueue) GetMapInfo() map[string]interface{} {
	// 读取已找到的物品信息
	mpItemsInfo := make([]map[string]interface{}, 0, 10)
	for _, ptItems := range this.arrItems {
		tInfo := ptItems.GetFieldInfo()
		mpItemsInfo = append(mpItemsInfo, tInfo)
	}
	return map[string]interface{}{
		"id":         this.id,
		"leadID":     this.leadID,
		"food":       this.food,
		"lenHero":    this.LenHero(),
		"itemInfo":   mpItemsInfo,
		"zoneFields": this.zoneFields.GetFieldInfo(),
	}
}

// 通过索引获取相应的动作
func (this *ExploreQueue) IdxAction(idx uint8) exploreaction.IFAction {
	switch idx {
	case constvalue.ACT_HUNT:
		return this.actHunt
	case constvalue.ACT_TREASURE:
		return this.actTreasure
	case constvalue.ACT_COLLECTION:
		return this.actCollection
	case constvalue.ACT_INSIGHT:
		return this.actInsight
	default:
		return this.actRest
	}
}

// 通过索引Key获取动作
func (this *ExploreQueue) KeyAction(idxKey string) exploreaction.IFAction {
	switch idxKey {
	case constvalue.STR_HUNT:
		return this.actHunt
	case constvalue.STR_TREASURE:
		return this.actTreasure
	case constvalue.STR_COLLECTION:
		return this.actCollection
	case constvalue.STR_INSIGHT:
		return this.actInsight
	default:
		return this.actRest
	}
}

// 获取LOG数组
func (this *ExploreQueue) GetLogInfos() []string {
	return this.arrLogs
}

// 改变当前状态
func (this *ExploreQueue) ChangeState(state uint8) {
	this.state = state
}

// 标记为新的Zone标识
func (this *ExploreQueue) JoinZone(zoneId int) {
	this.zoneID = zoneId
	this.expPower.Reset()
}

// 获得当前探索队可以执行的动作
//	@return
//		相应的动作
func (this *ExploreQueue) GetAction() exploreaction.IFAction {
	switch this.getActionIdx() {
	case constvalue.ACT_HUNT:
		return this.actHunt
	case constvalue.ACT_TREASURE:
		return this.actTreasure
	case constvalue.ACT_COLLECTION:
		return this.actCollection
	case constvalue.ACT_INSIGHT:
		return this.actInsight
	default:
		return this.actHunt
	}
}

func (this *ExploreQueue) getActionIdx() uint8 {
	// 判断当前能量值里是不是为0
	if this.expPower.PowerValue < 1 {
		// 为0则重新获取动作
		sKey := this.ptProportion.GetRndKey()
		return constvalue.GetActionIdx(sKey)
	} else {
		// 不为0则返回当前动作类型
		return this.expPower.PowerType
	}
}

func (this *ExploreQueue) GetExpPower() *ExplorePower {
	return this.expPower
}

func (this *ExploreQueue) initItemsToJson() string {
	arrMap := make([]map[string]interface{}, 0, len(this.arrItems))
	for _, item := range this.arrItems {
		str := item.GetSaveMap()
		fmt.Println(str)
		arrMap = append(arrMap, str)
	}
	btVal, err := json.Marshal(arrMap)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	result := string(btVal)
	fmt.Println("保存的信息JSON", result)
	return result
}

// 探索中将物品放入探索队里
//	也是事件中Event IFExploreQueue接口的实例
//	@parames
//		ptItems item.Items 	物品对象指针
//	@return
//		void
func (this *ExploreQueue) ExprloreGetItems(ptItems *item.Items) {
	// 记录物品
	fmt.Println("探索获得物品：", ptItems.ItemName(), ptItems.GetHow())
}

func (this *ExploreQueue) Save() error {
	saveMap := make(map[string]interface{})
	saveMap["food"] = this.food
	saveMap["zoneID"] = this.zoneID
	saveMap["startTime"] = vatools.GetTimeString(this.startTime)
	saveMap["findHow"] = this.findHow
	saveMap["physical"] = this.physical
	saveMap["state"] = this.state
	// 保存现在的物品
	saveMap["items"] = this.initItemsToJson()
	saveMap["zoneFields"] = this.zoneFields.GetJsonSave()
	this.SetInfo(saveMap)
	blnNew := this.IsNew()
	if blnNew {
		saveMap["leadID"] = this.leadID
	}
	err := this.BaseOB.Save()
	if err == nil && blnNew {
		this.id = int(this.GetLastAutoID())
		this.SetDBInfo("u_explorequeue", "*", "id", map[string]interface{}{"id": this.id, "leadID": this.leadID})
	}
	return err
}
