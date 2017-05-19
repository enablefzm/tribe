package gameroom

import (
	"tribe/jsondb"
)

// 通过Json格式字符构造一个Lead已开启的zone
//	@parames
//		strJson string
//	@return
//		*LeadZones
func NewLeadZonesOnStr(strJson string) *LeadZones {
	// 加载数据
	var jsonDB jsondb.LeadZones
	jsondb.JsdbUnJson(strJson, &jsonDB)
	ptLeadZones := &LeadZones{
		zones: make([]*LeadZone, len(jsonDB.Zones)),
	}
	for k, v := range jsonDB.Zones {
		ptLeadZones.zones[k] = &LeadZone{id: v.Id}
	}
	if len(ptLeadZones.zones) < 1 {
		ptLeadZones.AddZone(&LeadZone{id: 1001})
	}
	return ptLeadZones
}

// 玩家已打开的Zone信息
type LeadZones struct {
	zones []*LeadZone
}

func (this *LeadZones) GetZones() []*LeadZone {
	return this.zones
}

func (this *LeadZones) GetZonesInfo() []map[string]interface{} {
	res := make([]map[string]interface{}, len(this.zones))
	i := 0
	for _, v := range this.zones {
		res[i] = v.GetInfo()
		i++
	}
	return res
}

func (this *LeadZones) AddZone(ptZone *LeadZone) {
	this.zones = append(this.zones, ptZone)
}

func (this *LeadZones) GetSaveJson() string {
	db := &jsondb.LeadZones{
		Zones: make([]*jsondb.LeadZone, len(this.zones)),
	}
	for k, v := range this.zones {
		db.Zones[k] = &jsondb.LeadZone{Id: v.id}
	}
	str, err := jsondb.JsdbJson(db)
	if err != nil {
		return ""
	}
	return str
}

type LeadZone struct {
	id int
}

func (this *LeadZone) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"id": this.id,
	}
}
