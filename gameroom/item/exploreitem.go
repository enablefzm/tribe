package item

import (
	"errors"
)

// 被探索的物品
type ExploreItems struct {
	*Items
	rndVal uint16 // 可以被发现的随机数
	maxVal int    // 最大数
	nowHow int    // 当前数量
}

// Zone里探索的物品
func NewExploreItems(itemId, itemHow int, rndVal uint16) (*ExploreItems, error) {
	obItems, err := NewItems(itemId, itemHow)
	if err != nil {
		return nil, err
	}
	ob := &ExploreItems{
		Items:  obItems,
		rndVal: rndVal,
		maxVal: itemHow,
		nowHow: itemHow,
	}
	return ob, err
}

func (this *ExploreItems) OperateHow(how int) error {
	t := this.nowHow + how
	if t < 0 {
		return errors.New("数量不够")
	}
	this.nowHow = t
	return nil
}

func (this *ExploreItems) GetHow() int {
	return this.nowHow
}

func (this *ExploreItems) RndVal() uint16 {
	return this.rndVal
}
