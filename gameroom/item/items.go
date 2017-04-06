package item

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 带有数量的物品
type Items struct {
	*Item
	// id       int // 物品的ID
	// pid      int // 对应的玩家ID为0则未分配
	how      int // 商品数量
	property map[string]int
}

// 系统生成一个全新的物品
func NewItems(itemId, how int) (*Items, error) {
	obItem, err := OBManageItem.GetCanchItem(itemId)
	if err != nil {
		return nil, err
	}
	ob := &Items{
		Item: obItem,
		how:  how,
	}
	// 初始化基本属性
	ob.property = make(map[string]int, len(ob.Item.itemProperty))
	for k, v := range ob.Item.itemProperty {
		ob.property[k] = v
	}
	// 随机生成附加属性
	result := ob.GetRndProperty()
	for k, v := range result {
		ob.property[k] = v
	}
	return ob, err
}

// 通过玩家数据库里读取数据并实例化物品
//	func NewDBItems(rs map[string]interface{}) (*Items, error) {
//		pid := vatools.SInt(rs["pid"])
//		// 反序列JSON
//		mp := make(map[string]int)
//		err := json.Unmarshal([]byte(rs["property"]), &mp)
//		if err != nil {
//			return nil, err
//		}
//		ob, err := NewMapItems(pid, mp)
//		return ob, err
//	}

func NewMapItems(mapDB map[string]int) (*Items, error) {
	obItem, err := OBManageItem.GetCanchItem(mapDB["itemID"])
	if err != nil {
		return nil, err
	}
	ob := &Items{
		Item: obItem,
		how:  mapDB["how"],
	}
	il := len(mapDB)
	if il > 3 {
		il -= 3
	}
	ob.property = make(map[string]int, il)
	for k, v := range mapDB {
		if k == "itemID" || k == "how" || k == "pid" {
			continue
		}
		ob.property[k] = v
	}
	return ob, nil
}

func (this *Items) GetHow() int {
	if this.superposition {
		return this.how
	} else {
		return 1
	}
}

func (this *Items) OperateHow(how int) error {
	if this.superposition {
		this.how += how
	} else {
		return errors.New("这个物品不能叠加")
	}
	return nil
}

func (this *Items) GetFieldInfo() map[string]interface{} {
	res := this.Item.GetFieldInfo()
	res["how"] = this.GetHow()
	delete(res, "itemProperty")
	res["property"] = this.property
	return res
}

func (this *Items) GetSaveMap() map[string]interface{} {
	mapLen := len(this.property) + 2
	saveMap := make(map[string]interface{}, mapLen)
	saveMap["itemID"] = this.itemID
	saveMap["how"] = this.how
	for k, v := range this.property {
		saveMap[k] = v
	}
	return saveMap
}

// 获取要被保存的信息
func (this *Items) GetSaveJson() string {
	saveMap := this.GetSaveMap()
	btVal, err := json.Marshal(saveMap)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(btVal)
}
