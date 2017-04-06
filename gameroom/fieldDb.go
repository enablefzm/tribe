package gameroom

import (
	"strings"
	"tribe/baseob"
	"tribe/gameroom/event"
	"vava6/vatools"
)

// Zone里的Field对象
type FieldDb struct {
	*baseob.BaseOB
	id       int                      // 格子ID
	name     string                   // 格子名称
	fType    uint8                    // 格子类型
	defense  int                      // 格子防守值
	mpEvents map[uint8][]*event.Event // 存放各种事件、事件分类uint8
	images   []*FieldImage            // 格子的图片列表
	acts     []uint8                  // 动作列表
}

// 通过ID获取指定的FieldDb对象
func NewFieldDb(id int) (*FieldDb, error) {
	ob := &FieldDb{
		BaseOB:   baseob.NewBaseOB("d_field", "*", "id", map[string]interface{}{"id": id}),
		id:       id,
		images:   make([]*FieldImage, 0, 1),
		mpEvents: make(map[uint8][]*event.Event, 4),
	}
	// 加载数据
	rs, err := ob.LoadDbOnIdx(id)
	if err != nil {
		return nil, err
	}
	ob._init(rs)
	return ob, nil
}

// 加载数据
func (this *FieldDb) _init(rs map[string]string) {
	this.name = rs["name"]
	this.fType = vatools.SUint8(rs["fType"])
	// 构造事件列表
	// TODO...
	this._initEvent(rs["events"])
	// 构造图片列表
	this._initImages(rs["imgs"])

}

// 加载图片
//	str 格式：<图片ID>:[物件1],[物件2];<图片ID> ......
func (this *FieldDb) _initImages(str string) {
	arrs := strings.Split(str, ";")
	for _, v := range arrs {
		this.images = append(this.images, NewFieldImage(v))
	}
}

func (this *FieldDb) _initEvent(str string) {
	// 添加到指定的事件库里
	arr := strings.Split(str, ",")
	for _, v := range arr {
		ptEvent, err := event.NewEvent(vatools.SInt(v))
		if err != nil {
			continue
		}
		this._addEvent(ptEvent)
	}
}

func (this *FieldDb) _addEvent(ptEvent *event.Event) {
	eType := ptEvent.GetEventType()
	if _, ok := this.mpEvents[eType]; !ok {
		this.mpEvents[eType] = make([]*event.Event, 0, 2)
	}
	this.mpEvents[eType] = append(this.mpEvents[eType], ptEvent)
}

func (this *FieldDb) GetID() int {
	return this.id
}

// 随机获取需要加载的图片
func (this *FieldDb) GetFieldImage() *FieldImage {
	il := len(this.images)
	if il < 1 {
		return NewFieldImage("")
	}
	idx := 0
	if il > 1 {
		idx = vatools.CRnd(0, il-1)
	}
	return this.images[idx]
}

func (this *FieldDb) GetFieldInfo() map[string]interface{} {
	info := make(map[string]interface{}, 5)
	info["id"] = this.id
	info["name"] = this.name
	info["fType"] = this.fType
	return info
}

// 获取可以执行的动作
func (this *FieldDb) GetActs() []uint8 {
	if this.acts == nil {
		this.acts = make([]uint8, len(this.mpEvents))
		i := 0
		for act, _ := range this.mpEvents {
			this.acts[i] = act
			i++
		}
	}
	return this.acts
}

// 随机获取指定类型事件，从指定事件类型中随机抽取一个事件
//	@parames
//		actType uint8
//  @return
//		*event.Event, bool
func (this *FieldDb) GetEvent(actType uint8) (*event.Event, bool) {
	arrEvent, ok := this.mpEvents[actType]
	if !ok {
		return nil, false
	}
	il := len(arrEvent)
	if il < 1 {
		return nil, false
	}
	if il == 1 {
		return arrEvent[0], ok
	}
	idx := vatools.CRnd(0, il-1)
	return arrEvent[idx], ok
}
