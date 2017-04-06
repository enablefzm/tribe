package gameroom

import (
	"vava6/vatools"
)

// Zone里的格子数
type FieldZone struct {
	*FieldDb
	minHow uint16
	maxHow uint16
}

func NewFieldZone(fieldId int, minHow, maxHow uint16) (*FieldZone, error) {
	ptField, err := NewFieldDb(fieldId)
	if err != nil {
		return nil, err
	}
	ptFieldZone := &FieldZone{
		FieldDb: ptField,
		minHow:  minHow,
		maxHow:  maxHow,
	}
	return ptFieldZone, nil
}

// 获得随机数量范围
func (this *FieldZone) GetCreateHow() uint16 {
	if this.minHow == this.maxHow {
		return this.minHow
	}
	return uint16(vatools.CRnd(int(this.minHow), int(this.maxHow)))
}
