package exploreaction

import (
	"tribe/gameroom/constvalue"
)

// 狩猎动作
type Hunt struct {
	*Action
}

func NewHunt(actValue int) *Hunt {
	ptAction := NewAction(actValue)
	return &Hunt{
		ptAction,
	}
}

func (this *Hunt) GetActTypeOnID() uint8 {
	return constvalue.ACT_HUNT
}

func (this *Hunt) GetActName() string {
	return "狩猎"
}
