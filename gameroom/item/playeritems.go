package item

import (
	"encoding/json"
	"fmt"
	"tribe/baseob"
	"vava6/vatools"
)

// 玩家背包里的物品
type PlayerItems struct {
	baseob.BaseOB
	*Items
	id     int  // 物品ID
	pid    int  // 玩家所属ID
	isEdit bool // 是否被编辑过
}

// 通过数据库里的数据构造对象
func NewPlayerItemsOnRs(rs map[string]string) (*PlayerItems, error) {
	id := vatools.SInt(rs["id"])
	pid := vatools.SInt(rs["pid"])
	// 构造Items
	mp := make(map[string]int)
	err := json.Unmarshal([]byte(rs["property"]), &mp)
	if err != nil {
		return nil, err
	}
	obItems, err := NewMapItems(mp)
	if err != nil {
		return nil, err
	}
	ob := &PlayerItems{
		Items:  obItems,
		id:     id,
		pid:    pid,
		isEdit: false,
	}
	ob.SetDBInfo("u_items", "*", "id", map[string]interface{}{"id": ob.id})
	ob.SetNew(false)
	return ob, nil
}

// 通过Items转为玩家PlayerItems
func NewPlayerItemsOnItems(pid int, iptItems *Items) (*PlayerItems, error) {
	ob := &PlayerItems{
		Items:  iptItems,
		pid:    pid,
		isEdit: true,
	}
	ob.SetDBInfo("u_items", "*", "id", map[string]interface{}{"id": 0})
	ob.SetNew(true)
	// 直接保存
	err := ob.Save()
	if err != nil {
		return nil, err
	}
	ob.id = int(ob.GetLastAutoID())
	ob.SetKey(map[string]interface{}{"id": ob.id})
	return ob, nil
}

func (this *PlayerItems) GetID() int {
	return this.id
}

func (this *PlayerItems) GetPID() int {
	return this.pid
}

func (this *PlayerItems) SetEdit() {
	if this.isEdit != true {
		this.isEdit = true
	}
}

func (this *PlayerItems) GetFieldInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["id"] = this.id
	info["pid"] = this.pid
	info["items"] = this.Items.GetFieldInfo()
	return info
}

func (this *PlayerItems) OperateHow(how int) error {
	err := this.Items.OperateHow(how)
	if err == nil {
		this.SetEdit()
	}
	return err
}

func (this *PlayerItems) Save() error {
	if this.isEdit != true {
		return nil
	}
	saveMap := make(map[string]interface{})
	saveMap["property"] = this.GetSaveJson()
	if this.IsNew() {
		saveMap["pid"] = this.pid
		saveMap["itemId"] = this.itemID
	}
	this.SetInfo(saveMap)
	err := this.BaseOB.Save()
	fmt.Println(this.itemName, " 保存")
	if err == nil {
		this.isEdit = false
	}
	return err
}
