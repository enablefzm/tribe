package event

import (
	"errors"
	"fmt"
	"strings"
	"tribe/baseob"
	"vava6/vatools"
)

type IFEvent interface {
	GetID() int
	GetName() string
	GetEventType() uint16
	GetProbability() uint16
}

type IFSkill interface {
	GetSkillID() int
}

type Event struct {
	*baseob.BaseOB
	id          int        // 事件ID
	name        string     // 事件名称
	eventType   uint8      // 事件类型
	defense     int        // 事件防守值
	probability uint16     // 机率-数值为千分之几
	mustSkills  []IFSkill  // 必须要触发的技能
	obtains     []IFObtain // 事件奖励
}

func NewEvent(id int) (*Event, error) {
	ptEvent := &Event{
		BaseOB:     baseob.NewBaseOB("d_event", "*", "id", map[string]interface{}{"id": id}),
		id:         id,
		mustSkills: make([]IFSkill, 0),
	}
	rs, err := ptEvent.LoadDbOnIdx(id)
	if err != nil {
		return nil, err
	}
	// 构造对象
	ptEvent.name = rs["name"]
	ptEvent.eventType = vatools.SUint8(rs["eventType"])
	ptEvent.defense = vatools.SInt(rs["defense"])
	ptEvent.probability = vatools.SUint16(rs["probability"])
	// 生成奖励对象
	arr := strings.Split(rs["obtains"], ",")
	il := len(arr)
	ptEvent.obtains = make([]IFObtain, il)
	var j int
	for i := 0; i < il; i++ {
		ptObtain, err := NewIFObtain(vatools.SInt(arr[i]))
		if err != nil {
			continue
		}
		ptEvent.obtains[j] = ptObtain
		j++
	}
	if j != il {
		ptEvent.obtains = ptEvent.obtains[0:j]
	}
	// TODO...
	// 加载各项奖励和需要触发的技能
	return ptEvent, nil
}

func (this *Event) GetID() int {
	return this.id
}

func (this *Event) GetName() string {
	return this.name
}

func (this *Event) GetEventType() uint8 {
	return this.eventType
}

func (this *Event) GetProbability() uint16 {
	return this.probability
}

func (this *Event) GetMustSkill() []IFSkill {
	return this.mustSkills
}

func (this *Event) GetDefense() int {
	return this.defense
}

func (this *Event) GetObtain() (IFObtain, error) {
	// 事件奖励为平均随机
	arrLen := len(this.obtains)
	if arrLen < 1 {
		fmt.Println("该事件没有什么奖励:", arrLen)
		return nil, errors.New("NULL")
	}
	var ptObtain IFObtain
	idx := 0
	if arrLen > 1 {
		idx = vatools.CRnd(0, arrLen-1)
	}
	ptObtain = this.obtains[idx]
	return ptObtain, nil
}
