package exploreaction

// "tribe/gameroom/exploreskills"

type IFSkill interface {
	GetSkillID() uint16
	GetSkillName() string
	GetSkillPower() uint
}

type IFAction interface {
	GetActValue() int
	GetExpSkills() []IFSkill
	AddExpSkill(IFSkill)
	SetActValue(int)
	OperateValue(int)
	GetActTypeOnID() uint8
	GetActName() string
	// GetActTypeOnString() string
}

// 探索动作基类
type Action struct {
	actValue  int       // 探索动作值
	expSkills []IFSkill // 当前动作绑定的探索技能
}

func NewAction(actValue int) *Action {
	return &Action{
		actValue:  actValue,
		expSkills: make([]IFSkill, 0, 5),
	}
}

func (this *Action) GetActValue() int {
	return this.actValue
}

func (this *Action) GetExpSkills() []IFSkill {
	return this.expSkills
}

func (this *Action) AddExpSkill(skill IFSkill) {
	this.expSkills = append(this.expSkills, skill)
}

func (this *Action) SetActValue(setVal int) {
	this.actValue = setVal
}

func (this *Action) OperateValue(opVal int) {
	this.actValue += opVal
}
