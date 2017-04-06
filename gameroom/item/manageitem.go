package item

import (
	"sync"
)

// 基础物品信息
var OBManageItem = &ManageItem{
	mpItems: make(map[int]*Item, 1000),
	lkItems: new(sync.RWMutex),
}

type ManageItem struct {
	mpItems map[int]*Item
	lkItems *sync.RWMutex
}

// 通过缓存读取基础物品信息
func (this *ManageItem) GetCanchItem(itemId int) (*Item, error) {
	this.lkItems.RLock()
	ob, ok := this.mpItems[itemId]
	this.lkItems.RUnlock()
	if ok {
		return ob, nil
	} else {
		var err error
		this.lkItems.Lock()
		ob, ok = this.mpItems[itemId]
		if !ok {
			ob, err = NewItem(itemId)
			if err == nil {
				this.mpItems[itemId] = ob
			}
		}
		this.lkItems.Unlock()
		return ob, err
	}
}

func (this *ManageItem) Name() string {
	return "item"
}
func (this *ManageItem) Max() int {
	return this.Count()
}
func (this *ManageItem) Count() int {
	return len(this.mpItems)
}
func (this *ManageItem) Release() {
	for k, _ := range this.mpItems {
		delete(this.mpItems, k)
	}
}
func (this *ManageItem) Save() error {
	return nil
}
