package hero

//	1 - 战士    主力
//  2 - 猎人    主敏
//	3 - 巫师    主智
//  4 - 图腾师  主耐
const (
	H_WARRIOR uint8 = iota + 1
	H_HUNTER
	H_WIZARD
	H_TOTEM
)

type IFHero interface {
	AddExp(val int)
	GetAgile() int
	GetAtt() (int, int)
	GetCrit() int
	GetDef() int
	GetDodge() int
	GetHit() int
	GetPower() int
	GetSpeed() int
	GetStamina() int
	GetIq() int
	IsEpic() bool
	GetID() int
	GetInfo() map[string]interface{}
	GetMapInfo() map[string]interface{}
	GetName() string
	GetNextExp() int
	NeedFood() int
	Save() error
	GetLeadID() int
	SetLeadID(leadId int)
	UpLevel()
}

func NewIFHero(ptHero *Hero) IFHero {
	switch ptHero.heroType {
	case H_WARRIOR:
		return NewWarrior(ptHero)
	case H_HUNTER:
		return NewHunter(ptHero)
	case H_WIZARD:
		return NewWizard(ptHero)
	case H_TOTEM:
		return NewTotem(ptHero)
	default:
		return NewWarrior(ptHero)
	}
}
