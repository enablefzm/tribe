package exploreskills

import (
	"errors"
	"sync"
)

const (
	ADD_EXPLORE_PRO      = "ADD_EXPLORE_PRO"      // 增加探索值效果
	ADD_EXPLORE_HUNT_PRO = "ADD_EXPLORE_HUNT_PRO" // 增加狩猎加成值
)

type IFExploreSkill interface {
	GetName() string
	GetQuality() uint8
	// 获取探索技能效果
	//	@return
	//		如果有效果则返map[效果名称]效果值
	//		否则返回nil
	GetFeature() map[string]int
}

type BaseExploreSkill struct {
	name    string // 名称
	sType   uint8  // 技能所属的动作类型
	quality uint8  // 品质
}

func NewBaseExploreSkill(vName string, vQuality uint8) *BaseExploreSkill {
	return &BaseExploreSkill{
		name:    vName,
		quality: vQuality,
	}
}

func (this *BaseExploreSkill) GetName() string {
	return this.name
}

func (this *BaseExploreSkill) GetSkillType() uint8 {
	return this.sType
}

func (this *BaseExploreSkill) GetQuality() uint8 {
	return this.quality
}

// 探索技能构造管理对象
var mpExploreSkill = make(map[string]func() IFExploreSkill, 10)
var lkReg = new(sync.RWMutex)

// 将技能构造对象注册到管理对象
func RegExploreSkill(key string, createFunc func() IFExploreSkill) {
	lkReg.Lock()
	mpExploreSkill[key] = createFunc
	lkReg.Unlock()
}

func GetExploreSkill(key string) (IFExploreSkill, error) {
	if funCreateFunc, ok := mpExploreSkill[key]; !ok {
		return nil, errors.New("-1")
	} else {
		return funCreateFunc(), nil
	}
}
