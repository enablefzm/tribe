package gameroom

import (
	"strconv"
	"tribe/baseob"
)

func NewResource(id int) *Resource {
	ob := &Resource{}
	ob.id = id
	key := make(map[string]interface{})
	key["id"] = id
	rs, _ := ob.LoadDB("u_resource", "*", "id", key)
	if ob.IsNew() {
		ob.food = 2000
		ob.conch = 1000
	} else {
		ob.food, _ = strconv.Atoi(rs["food"])
		ob.conch, _ = strconv.Atoi(rs["conch"])
	}
	return ob
}

// 玩家资源对象
type Resource struct {
	baseob.BaseOB
	id    int
	food  int
	conch int
}

func (this *Resource) ID() int {
	return this.id
}

func (this *Resource) Food() int {
	return this.food
}

func (this *Resource) Conch() int {
	return this.conch
}

func (this *Resource) Save() error {
	saveInfo := make(map[string]interface{})
	saveInfo["id"] = this.id
	saveInfo["food"] = this.food
	saveInfo["conch"] = this.conch
	this.SetInfo(saveInfo)
	return this.BaseOB.Save()
}
