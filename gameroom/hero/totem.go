package hero

import (
	"fmt"
)

// 图腾师
type totem struct {
	*Hero
}

func NewTotem(iptHero *Hero) *totem {
	ob := &totem{
		iptHero,
	}
	ob.heroType = 4
	ob.eventUpLevel = ob.UpLevel
	return ob
}

func (this *totem) UpLevel() {
	fmt.Println("Totem uplevel ", this.level)
}
