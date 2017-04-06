package exploreaction

import (
	"tribe/gameroom/constvalue"
)

type Insight struct {
	*Action
}

func NewInsight(actValue int) *Insight {
	ptAction := NewAction(actValue)
	return &Insight{
		ptAction,
	}
}

func (this *Insight) GetActTypeOnID() uint8 {
	return constvalue.ACT_INSIGHT
}

func (this *Insight) GetActName() string {
	return "洞察"
}
