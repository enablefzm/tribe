package npc

import (
	"errors"
	"fmt"
	"tribe/sqldb"
	"vava6/vatools"
)

func NewMonster(id int) (*Monster, error) {
	rss, err := sqldb.Querys("d_monster", "id, name", fmt.Sprint("id=", id))
	if err != nil {
		return nil, err
	}
	if len(rss) != 1 {
		return nil, errors.New(fmt.Sprint("不存在ID为", id, "的Monster！"))
	}
	rs := rss[0]
	return &Monster{
		id:   vatools.SInt(rs["id"]),
		name: rs["name"],
	}, nil
}

type Monster struct {
	id   int
	name string
}

func (this *Monster) GetID() int {
	return this.id
}

func (this *Monster) GetName() string {
	return this.name
}

func (this *Monster) GetFieldInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":   this.GetID(),
		"name": this.GetName(),
	}
}
