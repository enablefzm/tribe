package gameroom

import (
	"errors"
	"fmt"
	"sync"
)

// 角色管理对象
var OBManageLead = &ManageLead{
	mapLeads: make(map[int]*Lead, 100),
	lk:       new(sync.RWMutex),
}

type ManageLead struct {
	mapLeads map[int]*Lead
	lk       *sync.RWMutex
}

// 添加角色到管理池
//	@parames
//		ld	*Lead	角色对象指针
func (this *ManageLead) AddLead(ld *Lead) *Lead {
	res := ld
	this.lk.RLock()
	if _, ok := this.mapLeads[ld.id]; !ok {
		this.lk.RUnlock()
		this.lk.Lock()
		if _, ok = this.mapLeads[ld.id]; !ok {
			this.mapLeads[ld.id] = ld
		} else {
			res = this.mapLeads[ld.id]
		}
		this.lk.Unlock()
	} else {
		this.lk.RUnlock()
	}
	return res
}

// 获取指定ID的角色，如果池里没有则从数据库里抓取，如果数据库里没有则创建一个
//	@parames
//		id	int		角色ID
//	@return
//		*Lead
func (this *ManageLead) GetLead(id int) *Lead {
	this.lk.RLock()
	ld, ok := this.mapLeads[id]
	this.lk.RUnlock()
	if !ok {
		this.lk.Lock()
		ld, ok = this.mapLeads[id]
		if !ok {
			ld = NewLead(id)
			this.mapLeads[ld.id] = ld
		}
		this.lk.Unlock()
	}
	return ld
}

func (this *ManageLead) Remove(id int) error {
	this.lk.RLock()
	ld, ok := this.mapLeads[id]
	this.lk.RUnlock()
	if !ok {
		return errors.New("指定的移除Lead对象不存在")
	}
	this.lk.Lock()
	ld.Save()
	delete(this.mapLeads, id)
	return nil
}

// 能过缓存里获取Lead并且不会主动创建
func (this *ManageLead) GetLeadNoCreate(id int) (*Lead, error) {
	var ld *Lead
	this.lk.RLock()
	ld, ok := this.mapLeads[id]
	this.lk.RUnlock()
	if !ok {
		tLd, err := newLeadInDB(id)
		if err != nil {
			return nil, err
		}
		// 加入缓存
		ld = this.AddLead(tLd)
	}
	return ld, nil
}

// 通过字符串
/*
func (this *ManageLead) GetLeadNoCreateInUID(sUid string) (*Lead, error) {
	return nil, nil
}
*/

// 清空所有对象
func (this *ManageLead) Clear() {
	this.lk.Lock()
	for k, _ := range this.mapLeads {
		this.mapLeads[k] = nil
		delete(this.mapLeads, k)
		fmt.Println("清空", k)
	}
	this.lk.Unlock()
}

func (this *ManageLead) Name() string {
	return "lead"
}
func (this *ManageLead) Max() int {
	return 1000
}
func (this *ManageLead) Count() int {
	return len(this.mapLeads)
}
func (this *ManageLead) Release() {
	for k, iptLead := range this.mapLeads {
		if !iptLead.IsOnline() {
			// iptLead.Save()
			iptLead.SaveAll()
			delete(this.mapLeads, k)
		}
	}
}
func (this *ManageLead) Save() error {
	return nil
}
