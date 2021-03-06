package hero

import (
	"fmt"
)

// 猎人
type hunter struct {
	*Hero
}

func NewHunter(iptHero *Hero) *hunter {
	ob := &hunter{
		iptHero,
	}
	ob.heroType = H_HUNTER
	ob.eventUpLevel = ob.UpLevel
	return ob
}

func (this *hunter) UpLevel() {
	fmt.Println("Hunter uplevel ", this.level)
}
