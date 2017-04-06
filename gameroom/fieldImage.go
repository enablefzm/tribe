package gameroom

import (
	"strings"
	"tribe/jsondb"
)

// 格子的图片
type FieldImage struct {
	img   string   // 底部图片
	thing []string // 物件
}

// <图片ID>:[物件1],[物件2]
func NewFieldImage(str string) *FieldImage {
	arrs := strings.Split(str, ":")
	pOb := &FieldImage{
		img:   arrs[0],
		thing: make([]string, 0),
	}
	if len(arrs) > 1 {
		arr := strings.Split(arrs[1], ",")
		for _, s := range arr {
			if len(s) > 1 {
				pOb.thing = append(pOb.thing, s)
			}
		}
	}
	return pOb
}

func NewFieldImageOnJson(db *jsondb.FieldImage) *FieldImage {
	pt := &FieldImage{
		img:   db.Img,
		thing: db.Thing,
	}
	return pt
}

func (this *FieldImage) GetFieldInfo() map[string]interface{} {
	return map[string]interface{}{
		"img":   this.img,
		"thing": this.thing,
	}
}
