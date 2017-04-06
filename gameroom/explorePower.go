package gameroom

import (
	"vava6/vatools"
)

type ExplorePower struct {
	PowerType  uint8 // 攻击类型
	PowerValue uint8 // 当暴发值
}

func NewExplorePower() *ExplorePower {
	return &ExplorePower{}
}

// 添加爆发值
//	成功则返回 true 失败则返回 false
//	失败则会让PowerValue和PowerType值归零
func (this *ExplorePower) AddVal(pType uint8) bool {
	if this.PowerType != pType {
		this.PowerType = pType
		this.PowerValue = 0
	}
	if this.PowerValue > 5 {
		return false
	}
	rndVal := vatools.CRnd(0, 100)
	okVal := 50
	switch this.PowerValue {
	case 0:
		okVal = 80
	case 1:
		okVal = 70
	case 2:
		okVal = 55
	case 3:
		okVal = 30
	case 4:
		okVal = 10
	case 5:
		okVal = 5
	}
	if rndVal <= okVal {
		this.PowerValue++
		return true
	} else {
		return false
	}
}

func (this *ExplorePower) Reset() {
	this.PowerType = 0
	this.PowerValue = 0
}
