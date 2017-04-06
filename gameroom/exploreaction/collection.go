package exploreaction

import (
	"tribe/gameroom/constvalue"
)

type Collection struct {
	*Action
}

func NewCollection(actValue int) *Collection {
	ptAction := NewAction(actValue)
	return &Collection{
		ptAction,
	}
}

func (this *Collection) GetActTypeOnID() uint8 {
	return constvalue.ACT_COLLECTION
}

func (this *Collection) GetActName() string {
	return "采集"
}
