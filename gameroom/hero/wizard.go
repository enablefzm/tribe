package hero

import (
	"fmt"
)

// 巫师
type wizard struct {
	*Hero
}

func NewWizard(iptHero *Hero) *wizard {
	ob := &wizard{
		iptHero,
	}
	ob.heroType = 3
	ob.eventUpLevel = ob.UpLevel
	return ob
}

func (this *wizard) UpLevel() {
	fmt.Println("Wizard uplevel ", this.level)
}
