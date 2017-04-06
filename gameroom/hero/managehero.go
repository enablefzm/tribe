package hero

import (
	"sync"
	"vava6/vatools"
)

var OBManageHero = &ManageHero{
	lk:     new(sync.RWMutex),
	mpHero: make(map[int]*Hero, 1000),
}

type ManageHero struct {
	lk     *sync.RWMutex
	mpHero map[int]*Hero
}

// 通过缓存获取英雄信息如果缓存里没有数据则通过数据库里加载
//	@parames
//		id	int
//	@return
//		*Hero
//		error
func (this *ManageHero) GetCacheHero(id int) (*Hero, error) {
	this.lk.RLock()
	obHero, ok := this.mpHero[id]
	this.lk.RUnlock()
	if ok {
		return obHero, nil
	}
	// 写锁
	this.lk.Lock()
	var err error
	obHero, ok = this.mpHero[id]
	if !ok {
		// 创建新对象
		obHero, err = NewHero(id)
		if err == nil {
			this.mpHero[id] = obHero
		}
	}
	this.lk.Unlock()
	return obHero, err
}

// 随机创建一个英雄
//	@parames
//		parames	map[string]string 动态参数
func (this *ManageHero) CreateHero(parames map[string]string) iHero {
	// 随机英雄类型
	// heroType := vatools.CRnd(1, 4)
	// 随机英雄品质

	ob := NewWarrior(&Hero{})
	ob.name = vatools.OBCreateName.GetName()
	// ob.att = vatools.CRnd(3, 10)
	ob.power = vatools.CRnd(1, 9)
	ob.leadID = 1
	ob.minAtt = 2
	ob.maxAtt = 3
	ob.speed = 10
	return ob
}
