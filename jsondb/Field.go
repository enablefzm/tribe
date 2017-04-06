package jsondb

import (
	"fmt"
	"reflect"
)

type Field struct {
	FieldID int         // 格子对应的FieldID
	Img     *FieldImage // 格子对应的图片
	PointX  int16       // 格子位于Zone的坐标X
	PointY  int16       // 格子位于Zone的坐标Y
	Musts   []int16     // 格子必须要触发的事件类型
	After   uint16      // 玩家角色经过格子的次数
}

func NewField() *Field {
	return &Field{
		Musts: make([]int16, 0, 3),
	}
}

func (this *Field) AddMusts(event int16) {
	this.Musts = append(this.Musts, event)
}

// ============================================================================
// TEST reflect
//	只用于测试
func JsdbReflect(source interface{}) {
	refType := reflect.TypeOf(source)
	fmt.Println("type:", refType, " kind:", refType.Kind())
	var s reflect.Value
	if refType.Kind() == reflect.Ptr {
		s = reflect.ValueOf(source).Elem()
	} else {
		s = reflect.ValueOf(source)
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Println(i, "Name:", s.Type().Field(i).Name, " Value:", f)
	}
}

func GetReflectValue(source interface{}) reflect.Value {
	refType := reflect.TypeOf(source)
	var s reflect.Value
	if refType.Kind() == reflect.Ptr {
		s = reflect.ValueOf(source).Elem()
	} else {
		s = reflect.ValueOf(source)
	}
	return s
}

func NewFieldOnObject(source interface{}) *FieldPlayer {
	s := GetReflectValue(source)
	ob := NewFieldPlayer()
	ob.PlayerID = int(s.FieldByName("playerID").Int())
	ob.ZoneId = int(s.FieldByName("zoneId").Int())
	ob.StartField = s.FieldByName("startField").String()
	//	ob.MaxX = s.FieldByName("maxX").String()
	return ob
}

// ============================================================================
