package gameroom

import (
	"tribe/jsondb"
)

// 玩家探索队的地图场景
type Field struct {
	fieldID int
	img     *FieldImage // 格子图片
	pointX  int16       // 坐标X
	pointY  int16       // 坐标Y
	musts   []int16     // 必定会被触发的事件（发现物品、怪物、事件、未知？）
	after   uint16      // 经过次数
}

func NewField(fieldId int) *Field {
	return &Field{
		fieldID: fieldId,
		musts:   make([]int16, 0),
	}
}

func (this *Field) GetMusts() []int16 {
	return this.musts
}

func (this *Field) GetFieldInfo() map[string]interface{} {
	return map[string]interface{}{
		"img":    this.img.GetFieldInfo(),
		"pointX": this.pointX,
		"pointY": this.pointY,
		"musts":  this.musts,
	}
}

func (this *Field) GetJsonDB() *jsondb.Field {
	res := jsondb.NewField()
	res.FieldID = this.fieldID
	res.Img = jsondb.NewFieldImage(this.img.img, this.img.thing)
	res.PointX = this.pointX
	res.PointY = this.pointY
	res.Musts = this.musts
	return res
}
