package skills

// 畜力对象
type Power struct {
	PowerType  uint8  // 畜力类型
	PowerValue uint16 // 畜力值
}

func NewPower() *Power {
	return &Power{
		PowerType:  0,
		PowerValue: 0,
	}
}
