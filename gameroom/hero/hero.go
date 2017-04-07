package hero

import (
	"fmt"
	"tribe/baseob"
	"tribe/gameroom/item"
	"vava6/vatools"
)

// 英雄分类
//	1 - 战士    主力
//  2 - 猎人    主敏
//	3 - 巫师    主智
//  4 - 图腾师  主耐
//	英雄品质
//		等级   初始   升级点数
//		 1     10     4 (2-2)
//       2     12     5 (2-3)
//       3     14     7 (3-4)
//       4     16     9 (4-5)
type Hero struct {
	baseob.BaseOB
	id           int               // 英雄ID
	leadID       int               // 所属角色ID
	name         string            // 英雄名称
	heroType     uint8             // 英雄类型
	epic         bool              // 是否是史实英雄
	quality      uint8             // 英雄品质
	exp          int               // 经验值
	level        int               // 等级
	hp           int               // 生命值
	minAtt       int               // 攻击力
	maxAtt       int               // 攻击力
	speed        int               // 攻速
	power        int               // 力量
	stamina      int               // 体力
	agile        int               // 敏捷
	iq           int               // 智力
	def          int               // 防御力
	hit          int               // 命中
	crit         int               // 暴击
	dodge        int               // 闪避
	pHead        *item.PlayerItems // 头部装备
	pChest       *item.PlayerItems // 胸部装备
	pLeg         *item.PlayerItems // 腿部装备
	pFoot        *item.PlayerItems // 脚部装备
	pTrinket     *item.PlayerItems // 饰品
	upPoint      int16             // 成长点数
	mainPoint    int16             // 主成长点数
	eventUpLevel func()            // 升级事件
}

// 通过数据库构造对象
func NewHero(id int) (*Hero, error) {
	ob := &Hero{}
	rs, err := ob.LoadDB("u_hero", "*", "id", map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	ob.id = id
	ob.leadID = vatools.SInt(rs["leadID"])
	ob.name = rs["name"]
	ob.heroType = vatools.SUint8(rs["heroType"])
	// ob.att = vatools.SInt(rs["att"])
	ob.power = vatools.SInt(rs["power"])
	return ob, nil
}

// 通过英雄类型和品质创建英雄
func NewCreateHero(heroType, quality uint8) IFHero {
	obHero := &Hero{
		name:     vatools.OBCreateName.GetName(),
		heroType: heroType,
		quality:  quality,
	}
	i := 0
	switch quality {
	case 1:
		obHero.upPoint = 2
		obHero.mainPoint = 2
		i = 10
	case 2:
		obHero.upPoint = 2
		obHero.mainPoint = 3
		i = 12
	case 3:
		obHero.upPoint = 3
		obHero.mainPoint = 4
		i = 14
	default:
		obHero.upPoint = 4
		obHero.mainPoint = 5
		i = 16
	}
	// 分配初始点数
	//	power        int               // 力量
	//	stamina      int               // 体力
	//	agile        int               // 敏捷
	//	iq           int               // 智力
	obHero.level = 1
	for j := 0; j < i; j++ {
		rndVal := vatools.CRnd(1, 4)
		switch rndVal {
		case 1:
			obHero.power += 1
		case 2:
			obHero.stamina += 1
		case 3:
			obHero.agile += 1
		case 4:
			obHero.iq += 1
		}
	}
	// 生成相应类型的英雄对象
	return NewIFHero(obHero)
}

func (this *Hero) GetID() int {
	return this.id
}

// 获取英雄名称
func (this *Hero) GetName() string {
	return this.name
}

// 获得攻击力
func (this *Hero) GetAtt() (int, int) {
	return this.minAtt, this.maxAtt
}

// 获得攻击速度
func (this *Hero) GetSpeed() int {
	return this.speed
}

// 获得力量
func (this *Hero) GetPower() int {
	return this.power
}

// 获得体力
func (this *Hero) GetStamina() int {
	return this.stamina
}

// 获得敏捷
func (this *Hero) GetAgile() int {
	return this.agile
}

// 获得智力
func (this *Hero) GetIq() int {
	return this.iq
}

// 获得防御力
func (this *Hero) GetDef() int {
	return this.def
}

// 获得命中
func (this *Hero) GetHit() int {
	return this.hit
}

// 获得暴击
func (this *Hero) GetCrit() int {
	return this.crit
}

// 获得闪避
func (this *Hero) GetDodge() int {
	return this.dodge
}

// 英雄所属的角色ID
func (this *Hero) GetLeadID() int {
	return this.leadID
}

func (this *Hero) SetLeadID(leadID int) {
	this.leadID = leadID
}

func (this *Hero) GetInfo() map[string]interface{} {
	res := map[string]interface{}{
		"id":       this.id,
		"leadID":   this.leadID,
		"name":     this.name,
		"heroType": this.heroType,
		"power":    this.power,
	}
	return res
}

// 一分钟需要消耗的食物量
func (this *Hero) NeedFood() int {
	return 2
}

func (this *Hero) IsEpic() bool {
	return this.epic
}

// 向角色添加经验
func (this *Hero) AddExp(val int) {
	this.exp += val
	for {
		if this.exp < this.GetNextExp() {
			break
		}
		// 成功一级
		this.exp -= this.GetNextExp()
		this.UpLevel()
	}
}

// 等级升级
func (this *Hero) UpLevel() {
	// 成长
	if this.eventUpLevel != nil {
		this.eventUpLevel()
	} else {
		fmt.Println("BASE UP Level ", this.level)
	}
}

// 获取下一级需要等级
func (this *Hero) GetNextExp() int {
	nextExp := this.level * (this.level + 5) * 10
	return nextExp
}

func (this *Hero) GetMapInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["id"] = this.id
	info["leadID"] = this.leadID
	info["name"] = this.name
	info["heroType"] = this.heroType
	info["quality"] = this.quality
	info["minAtt"] = this.minAtt
	info["maxAtt"] = this.maxAtt
	info["speed"] = this.speed
	info["power"] = this.power
	info["stamina"] = this.stamina
	info["agile"] = this.agile
	info["iq"] = this.iq
	info["def"] = this.def
	info["hit"] = this.hit
	info["crit"] = this.crit
	info["dodge"] = this.dodge
	info["upPoint"] = this.upPoint
	info["mainPoint"] = this.mainPoint
	return info
}

func (this *Hero) Save() error {
	this.BaseOB.SetDBInfo("u_hero", "*", "id", map[string]interface{}{"id": this.id})
	saveMap := this.GetMapInfo()
	delete(saveMap, "id")
	if this.id == 0 {
		this.SetNew(true)
	} else {
		delete(saveMap, "leadID")
		delete(saveMap, "heroType")
		delete(saveMap, "upPoint")
		delete(saveMap, "mainPoint")
	}
	if this.IsEpic() {
		saveMap["epic"] = 1
	} else {
		saveMap["epic"] = 0
	}
	// 保存装备ID
	if this.pChest != nil {
		saveMap["pChest"] = this.pChest.GetID()
	} else {
		saveMap["pChest"] = nil
	}
	if this.pFoot != nil {
		saveMap["pFoot"] = this.pFoot.GetID()
	} else {
		saveMap["pFoot"] = nil
	}
	if this.pHead != nil {
		saveMap["pHead"] = this.pHead.GetID()
	} else {
		saveMap["pHead"] = nil
	}
	if this.pTrinket != nil {
		saveMap["pTrinket"] = this.pTrinket.GetID()
	} else {
		saveMap["pTrinket"] = nil
	}
	this.BaseOB.SetInfo(saveMap)
	err := this.BaseOB.Save()
	if this.id == 0 {
		this.id = int(this.GetLastAutoID())
	}
	return err
}
