package gameroom

import (
	"errors"
	"fmt"
	"sort"
	"tribe/jsondb"
	"vava6/vatools"
)

const (
	PLAIN    = iota + 1 // 平原 1
	HILLS               // 山丘 2
	JUNGLE              // 丛林 3
	RIVERWAY            // 河道 4
	CAVE                // 洞穴 5
)

const (
	Z_UP        = "up"
	Z_DOWN      = "down"
	Z_LEFTUP    = "leftUp"
	Z_LEFTDOWN  = "leftDown"
	Z_RIGHTUP   = "rightUp"
	Z_RIGHTDOWN = "rightDown"
)

// 玩家的探索队的地图
type FieldPlayer struct {
	zoneId     int   // 对应的zoneID
	playerID   int   // 玩家ID
	maxX       int16 // 尺寸X
	maxY       int16 // 尺寸Y
	fields     map[string]*Field
	nowField   *Field
	startField *Field
}

func NewFieldPlayer(zoneId int) *FieldPlayer {
	return &FieldPlayer{
		zoneId: zoneId,
		maxX:   5,
		maxY:   5,
		fields: make(map[string]*Field, 10),
	}
}

// 从数据JSON实例化格子对象
func NewFieldPlayerOnJson(strJson string) *FieldPlayer {
	// 将json转为DB
	var dbPlayerField jsondb.FieldPlayer
	jsondb.JsdbUnJson(strJson, &dbPlayerField)
	obFieldPlayer := NewFieldPlayer(0)
	obFieldPlayer.zoneId = dbPlayerField.ZoneId
	obFieldPlayer.playerID = dbPlayerField.PlayerID
	obFieldPlayer.maxX = dbPlayerField.MaxX
	obFieldPlayer.maxY = dbPlayerField.MaxY
	for _, v := range dbPlayerField.ArrField {
		tField := NewField(v.FieldID)
		tField.img = NewFieldImageOnJson(v.Img)
		tField.pointX = v.PointX
		tField.pointY = v.PointY
		tField.musts = v.Musts
		tField.after = v.After
		obFieldPlayer.fields[obFieldPlayer.GetXYKey(tField.pointX, tField.pointY)] = tField
	}
	// 获得入口位置
	obFieldPlayer.startField = obFieldPlayer.fields[dbPlayerField.StartField]
	// 获得当前位置
	obFieldPlayer.nowField = obFieldPlayer.fields[dbPlayerField.NowField]
	return obFieldPlayer
}

// 添加场景
func (this *FieldPlayer) AddField(field *Field) {
	var nx, ny int16
	var x, y int16
	if this.nowField == nil {
		// 随机获取初生地址
		x = 1
		y = int16(vatools.CRnd(0, 3))
		this.startField = field
	} else {
		nx, ny = this.nowField.pointX, this.nowField.pointY
		// 获取下当前field周边位置
		x, y = this.GetRndNextNilXY(nx, ny)
		if x == -1 || y == -1 {
			// 周片没有可以放入的位置
			// 从现有的随机挑一个路口
			fmt.Println("当前位置", nx, ny, " 目标位置", x, y)
			// 周边位置
			fmt.Println("周片位置：", this.GetFieldRange(nx, ny))
			for _, t := range this.fields {
				x, y = this.GetRndNextNilXY(t.pointX, t.pointY)
				if x != -1 && y != -1 {
					break
				}
			}
		}
		if x == -1 || y == -1 {
			fmt.Println("周片都没有位置放了", nx, ny)
			// 没有空间增加了...
			return
		}
	}
	// 增加成功
	field.pointX = x
	field.pointY = y
	this.fields[this.GetXYKey(x, y)] = field
	this.nowField = field
	// fmt.Println("添加新场景：", field.fieldType, " 坐标x", x, " y", y)
	return
}

// 获得PlayerFields下属的所有Field
func (this *FieldPlayer) GetFields() map[string]*Field {
	return this.fields
}

// 获取当前格
func (this *FieldPlayer) GetNowField() *Field {
	// 如果单前格子为空则随机挑一个
	if this.nowField == nil {
		for _, v := range this.fields {
			this.nowField = v
			break
		}
	}
	return this.nowField
}

func (this *FieldPlayer) GetNowFieldErr() (*Field, error) {
	nowField := this.GetNowField()
	if nowField == nil {
		return nil, errors.New("NULL")
	} else {
		return nowField, nil
	}
}

// 移动到下一格
func (this *FieldPlayer) NextField() (*Field, error) {
	nowField := this.GetNowField()
	if nowField == nil {
		return nil, errors.New("NULL")
	}
	return this.GetNextField(nowField)
}

// 获取旁边格子
func (this *FieldPlayer) GetNextField(nowField *Field) (*Field, error) {
	// 获取旁边格子数
	x, y := nowField.pointX, nowField.pointY
	mpRange := this.GetFieldRange(x, y)
	il := len(mpRange)
	if il < 1 {
		return nil, errors.New("NO ROUTE")
	}
	fields := make([]*Field, 0, 6)
	for _, v := range mpRange {
		pt, ok := this.fields[this.GetXYKey(v[0], v[1])]
		if ok {
			fields = append(fields, pt)
		}
	}
	// 排序
	sort.Slice(fields, func(i, j int) bool { return fields[i].after < fields[j].after })
	// 判断是否要走哪条路
	rndVal := vatools.CRnd(0, 100)
	if rndVal < 60 {
		return fields[0], nil
	}
	rndKey := vatools.CRnd(1, len(fields)) - 1
	return fields[rndKey], nil
}

func (this *FieldPlayer) GetFieldInfo() map[string]interface{} {
	res := make(map[string]interface{}, 10)
	res["zoneId"] = this.zoneId
	res["maxX"] = this.maxX
	res["maxY"] = this.maxY
	res["count"] = len(this.fields)
	// res["fields"] = make(map[string]interface{})
	info := make(map[string]interface{})
	for k, v := range this.fields {
		info[k] = v.GetFieldInfo()
	}
	res["fields"] = info
	return res
}

// 	获得位置
//		up
//		down
//		leftUp
//		leftDown
//		rightUp
//		rightDown
func (this *FieldPlayer) GetFieldRange(nx, ny int16) map[string][2]int16 {
	cRange := make(map[string][2]int16, 6)
	// 计算出旁边所有范围
	// 获得计算当前是奇数还是偶数
	var isQi bool
	if nx%2 > 0 {
		isQi = true
	} else {
		isQi = false
	}
	// 获得上下左右位置
	up := ny - 1
	down := ny + 1
	left := nx - 1
	right := nx + 1

	// 获取上
	if up >= 0 && nx >= 0 {
		cRange[Z_UP] = [2]int16{nx, up}
	}
	// 获取下
	if down <= this.maxY && nx >= 0 {
		cRange[Z_DOWN] = [2]int16{nx, down}
	}

	if isQi == true {
		// 计算左上
		if left >= 0 {
			if ny >= 0 {
				cRange[Z_LEFTUP] = [2]int16{left, ny}
			}
			// 计算左下
			if down <= this.maxY {
				cRange[Z_LEFTDOWN] = [2]int16{left, down}
			}
		}
		// 计算右上
		if right <= this.maxX {
			if ny >= 0 {
				cRange[Z_RIGHTUP] = [2]int16{right, ny}
			}
			if down <= this.maxY {
				cRange[Z_RIGHTDOWN] = [2]int16{right, down}
			}
		}
	} else {
		// 偶数
		if left >= 0 {
			// 左上
			if up >= 0 {
				cRange[Z_LEFTUP] = [2]int16{left, up}
			}
			// 左下
			if ny <= this.maxY {
				cRange[Z_LEFTDOWN] = [2]int16{left, ny}
			}
		}
		// 计算右上下
		if right <= this.maxX {
			if up >= 0 {
				cRange[Z_RIGHTUP] = [2]int16{right, up}
			}
			if ny <= this.maxY {
				cRange[Z_RIGHTDOWN] = [2]int16{right, ny}
			}
		}
	}
	return cRange
}

// 获得下一个随机目标位置
//  如果下一个目录没有可以放入Field的空位则返 x = -1, y = -1
//	并且更多机率是向右下方向偏移
func (this *FieldPlayer) GetRndNextNilXY(nx, ny int16) (x, y int16) {
	x, y = -1, -1
	// 获取周边位置
	mpFieldXY := this.GetFieldRange(nx, ny)
	il := len(mpFieldXY)
	if il == 0 {
		return
	}
	arrRight := [3]string{Z_RIGHTUP, Z_RIGHTDOWN, Z_DOWN}
	arrLeft := [3]string{Z_UP, Z_LEFTDOWN, Z_LEFTUP}
	arrDir := [2][3]string{}
	rnd := vatools.CRnd(1, 100)
	if rnd < 80 {
		arrDir[0] = arrRight
		arrDir[1] = arrLeft
	} else {
		arrDir[0] = arrLeft
		arrDir[1] = arrRight
	}
	for i := 0; i < 2; i++ {
		j := vatools.CRnd(0, 2)
		for l := 0; l < 3; l++ {
			v, ok := mpFieldXY[arrDir[i][j]]
			if ok {
				if !this.checkFieldIsNil(v[0], v[1]) {
					x, y = v[0], v[1]
					return
				}
			}
			j++
			if j > 2 {
				j = 0
			}
		}
	}
	return
}

// 判断当前Fields是否有这个位置的Field
func (this *FieldPlayer) checkFieldIsNil(nx, ny int16) bool {
	_, ok := this.fields[this.GetXYKey(nx, ny)]
	return ok
}

func (this *FieldPlayer) GetXYKey(nx, ny int16) string {
	return fmt.Sprint("P_", nx, "_", ny)
}

func (this *FieldPlayer) GetSave() string {
	ptJsonFieldPlayer := jsondb.NewFieldPlayer()
	ptJsonFieldPlayer.ZoneId = this.zoneId
	ptJsonFieldPlayer.PlayerID = this.playerID
	ptJsonFieldPlayer.MaxX = this.maxX
	ptJsonFieldPlayer.MaxY = this.maxY
	if this.startField != nil {
		ptJsonFieldPlayer.StartField = fmt.Sprint("P_", this.startField.pointX, "_", this.startField.pointY)
	}
	if this.nowField != nil {
		ptJsonFieldPlayer.NowField = fmt.Sprint("P_", this.nowField.pointX, "_", this.nowField.pointY)
	}
	for _, ptField := range this.fields {
		ptJsonFieldPlayer.AddField(ptField.GetJsonDB())
	}
	strJson, err := jsondb.JsdbJson(ptJsonFieldPlayer)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(strJson)
	}
	return strJson
}

func (this *FieldPlayer) GetJsonSave() string {
	return this.GetSave()
}
