package hero

import (
	"vava6/vatools"
)

// 通过酒馆生成英雄
func CreatePub() IFHero {
	// 随机英雄类型
	rndHeroType := vatools.CRnd(1, 4)
	// 获取品质
	rndHeroQuality := vatools.CRnd(1, 1000)
	heroQuality := 1
	switch {
	case rndHeroQuality > 500 && rndHeroQuality < 700:
		heroQuality = 2
	case rndHeroQuality >= 750 && rndHeroQuality < 900:
		heroQuality = 3
	case rndHeroQuality >= 900:
		heroQuality = 4
	default:
		heroQuality = 1
	}
	return NewCreateHero(uint8(rndHeroType), uint8(heroQuality))
}
