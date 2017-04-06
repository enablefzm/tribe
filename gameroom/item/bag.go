package item

import (
	"errors"
	"fmt"
	"sync"
	"tribe/baseob"
	"tribe/sqldb"
	"vava6/vatools"
)

type Bag struct {
	baseob.BaseOB
	id      int                  // 玩家角色所属ID
	value   uint16               // 背包大小
	mpItems map[int]*PlayerItems // 物品容器
	lk      *sync.RWMutex        // 锁
}

func NewBag(leadId int) (*Bag, error) {
	ob := &Bag{
		id: leadId,
		lk: new(sync.RWMutex),
	}
	rs, err := ob.LoadDB("u_bag", "*", "id", map[string]interface{}{"id": leadId})
	if err != nil {
		if err.Error() != baseob.ERR_NULL {
			fmt.Println("NewBag:", err.Error())
			return nil, err
		} else {
			// 新增玩家角色背包
			// 新增对象定
			ob.value = 20
			ob.mpItems = make(map[int]*PlayerItems, ob.value)
			ob.SetNew(true)
		}
	} else {
		// 加载原有数据
		ob.SetNew(false)
		// 加载背包现有大小
		ob.value = vatools.SUint16(rs["value"])
		// 加载现有物品
		rss, err := sqldb.Querys("u_items", "*", fmt.Sprint("pid=", ob.id))
		ob.mpItems = make(map[int]*PlayerItems, len(rss))
		if err == nil {
			for _, rs := range rss {
				iptPItems, err := NewPlayerItemsOnRs(rs)
				if err != nil {
					// 错误处理
					fmt.Println("加载物品出错：", err.Error())
				} else {
					ob.mpItems[iptPItems.id] = iptPItems
				}
			}
		}
	}
	return ob, nil
}

// 获取背包ID
func (this *Bag) GetID() int {
	return this.id
}

// 获取背包数量
func (this *Bag) GetValue() uint16 {
	return this.value
}

// 获取当前背包数量
func (this *Bag) LenItems() int {
	return len(this.mpItems)
}

// 获取当前背包里物品数量
func (this *Bag) GetItems() []*PlayerItems {
	resPItems := make([]*PlayerItems, 0, this.LenItems())
	for _, pItems := range this.mpItems {
		resPItems = append(resPItems, pItems)
	}
	return resPItems
}

// 增加背包格子数量
func (this *Bag) AddValue(newVal uint16) {
	this.value += newVal
	err := this.save()
	if err != nil {
		fmt.Println("保存背包出错：", err.Error())
	}
}

// 将物品放入背包
//	@parames
//		iptItems *Items 物品
//	@return
//		int (0成功, 1背包满了, 2转换物品出错, 3物品保存出错)
func (this *Bag) PutItems(iptItems *Items) (*PlayerItems, error) {
	this.lk.Lock()
	defer this.lk.Unlock()
	// 判断物品是否能够堆叠
	if iptItems.Superposition() == true {
		// 查找背包里是否有这个物品
		for _, pItems := range this.mpItems {
			if pItems.itemID == iptItems.itemID {
				err := pItems.OperateHow(iptItems.how)
				if err != nil {
					return nil, err
				}
				return pItems, nil
			}
		}
	}
	if uint16(len(this.mpItems)) >= this.value {
		return nil, errors.New("1")
	}
	// 转换成玩家物品
	pItems, err := NewPlayerItemsOnItems(this.id, iptItems)
	if err != nil {
		fmt.Println(this.id, "生成玩家物品出错：", err.Error())
		return nil, errors.New("2")
	}
	this.mpItems[pItems.id] = pItems
	return pItems, nil
}

// 保存物品
func (this *Bag) Save() error {
	if this.IsNew() {
		err := this.save()
		if err != nil {
			fmt.Println("保存背包出错：", err.Error())
		}
	}
	// 保存背包里所有数据
	this.lk.RLock()
	for id, iptItem := range this.mpItems {
		err := iptItem.Save()
		if err != nil {
			fmt.Println("背存玩家物品出错：", id, err.Error())
		}
	}
	this.lk.RUnlock()
	return nil
}

func (this *Bag) save() error {
	saveMap := make(map[string]interface{})
	saveMap["value"] = this.value
	if this.IsNew() {
		saveMap["id"] = this.id
	}
	this.SetInfo(saveMap)
	err := this.BaseOB.Save()
	return err
}

//	// 取出Items物品
//	//	@parames
//	//		itemID 物品ID
//	func (this *Bag) GetItems(itemID int) (*Items, error) {

//	}
