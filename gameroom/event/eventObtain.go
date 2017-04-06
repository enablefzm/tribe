package event

import (
	"tribe/gameroom/item"
	"vava6/vatools"
)

// 物品为获得对象接口
type IFObtain interface {
	ObtainDo(IFExploreQueue)
	GetInfo() string
}

type IFExploreQueue interface {
	ExploreGetItems(*item.Items) // 玩家探索获得物品的
}

// 事件获得的奖励
type Obtain struct {
	id          int
	minGet      uint8  // 可以获得数量
	maxGet      uint8  // 最大可以获得的数量
	value       int    // 当前数量
	maxValue    int    // 最大能够增长的数量
	isBoundless bool   // 是否是无限资源
	stepUp      uint16 // 一次增长多少个
	stepTime    uint16 // 多久增长一次以秒为单位 60s 一分钟增长一次
	nowTime     uint   // 记录增长时间
}

func (this *Obtain) GetID() int {
	return this.id
}

func (this *Obtain) MinGet() uint8 {
	return this.minGet
}

func (this *Obtain) MaxGet() uint8 {
	return this.maxGet
}

func (this *Obtain) GetValue() int {
	return this.value
}

func (this *Obtain) IsBoundless() bool {
	return this.isBoundless
}

// 获取前可以得到的数量
func (this *Obtain) GetRndValue() uint8 {
	if !this.isBoundless {
		if this.value < 1 {
			return 0
		}
		res := this.getRndValue()
		if res > this.value {
			res = this.value
		}
		this.value -= res
		return uint8(res)
	} else {
		return uint8(this.getRndValue())
	}
}

func (this *Obtain) getRndValue() int {
	return vatools.CRnd(int(this.minGet), int(this.maxGet))
}
