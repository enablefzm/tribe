package exploreaction

import (
	"tribe/gameroom/constvalue"
)

type Rest struct {
	*Action
}

func NewRest(actValue int) *Insight {
	ptAction := NewAction(actValue)
	return &Insight{
		ptAction,
	}
}

func (this *Rest) GetActTypeOnID() uint8 {
	return constvalue.ACT_REST
}

func (this *Rest) GetActName() string {
	return "休息"
}
