package hero

type iHero interface {
	AddExp(val int)
	GetID() int
	GetInfo() map[string]interface{}
	GetMapInfo() map[string]interface{}
	GetName() string
	GetNextExp() int
	NeedFood() int
	Save() error
	SetLeadID(leadId int)
	UpLevel()
}
