package jsondb

// 玩家探索队列的地图Field
//	Json
type FieldPlayer struct {
	ZoneId     int
	PlayerID   int
	MaxX       int16
	MaxY       int16
	ArrField   []*Field
	NowField   string
	StartField string
}

func NewFieldPlayer() *FieldPlayer {
	return &FieldPlayer{
		ArrField: make([]*Field, 0, 10),
	}
}

func (this *FieldPlayer) AddField(ptField *Field) {
	this.ArrField = append(this.ArrField, ptField)
}
