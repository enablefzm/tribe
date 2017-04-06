package exploreskills

import (
	"tribe/gameroom/constvalue"
)

// 挖矿
type Mining struct {
	*BaseExploreSkill
}

func (this *Mining) GetFeature() map[string]int {
	return nil
}

func init() {
	RegExploreSkill("mining", func() IFExploreSkill {
		pt := &Mining{
			BaseExploreSkill: NewBaseExploreSkill("挖矿", 1),
		}
		pt.sType = constvalue.ACT_COLLECTION
		return pt
	})
}
