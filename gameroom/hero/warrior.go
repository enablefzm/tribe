package hero

import (
	"fmt"
)

// 战士
type warrior struct {
	*Hero
}

func NewWarrior(iptHero *Hero) *warrior {
	ob := &warrior{
		iptHero,
	}
	ob.heroType = H_WARRIOR
	ob.eventUpLevel = ob.UpLevel
	return ob
}

func (this *warrior) UpLevel() {
	fmt.Println("warrior uplevel ", this.level)
}
