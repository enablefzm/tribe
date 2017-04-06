package exploreaction

import (
	"tribe/gameroom/constvalue"
)

type Treasure struct {
	*Action
}

func NewTreasure(actValue int) *Treasure {
	ptAction := NewAction(actValue)
	return &Treasure{
		ptAction,
	}
}

func (this *Treasure) GetActTypeOnID() uint8 {
	return constvalue.ACT_TREASURE
}

func (this *Treasure) GetActName() string {
	return "寻宝"
}
