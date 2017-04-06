package event

import (
	"fmt"
	"tribe/gameroom/item"
	"vava6/vatools"
)

// 构造物品奖励对象
func NewObtainItems(ptObtain *Obtain, db string) *ObtainItems {
	pt := &ObtainItems{
		Obtain: ptObtain,
		itemId: vatools.SInt(db),
	}
	return pt
}

// 这是物品奖励
type ObtainItems struct {
	*Obtain
	itemId int
}

func (this *ObtainItems) ObtainDo(iQueue IFExploreQueue) {
	how := this.getRndValue()
	if how < 1 {
		// TODO..
		//	没有资源信息的处理接口
		fmt.Println("当前事件资源数量为0")
		return
	}
	ptItems, err := item.NewItems(this.itemId, how)
	if err != nil {
		// TODO...
		// 生成物品时发生错误
		fmt.Println(err.Error())
		return
	}
	// 将物品放ExploreItems
	iQueue.ExploreGetItems(ptItems)
}

func (this *ObtainItems) GetInfo() string {
	return "获得了物品"
}
