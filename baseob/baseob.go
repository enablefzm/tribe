package baseob

import (
	"errors"
	"fmt"
	"tribe/sqldb"
)

const (
	ERR_NULL    string = "ERR_NULL"
	ERR_PARAMES string = "ERR_PARAMES"
	ERR_SQLERR  string = "ERR_SQLERR"
)

func NewBaseOB(table, field, idx string, key map[string]interface{}) *BaseOB {
	return &BaseOB{
		table: table,
		field: field,
		idx:   idx,
	}
}

type BaseOB struct {
	table      string
	field      string
	key        map[string]interface{}
	idx        string
	isNew      bool
	info       map[string]interface{}
	lastAutoID int64
}

func (this *BaseOB) LoadDB(table, field, idx string, key map[string]interface{}) (map[string]string, error) {
	this.SetDBInfo(table, field, idx, key)
	return this.readDB(key)
}

// 通过主键值获取数据
func (this *BaseOB) LoadDbOnIdx(k interface{}) (map[string]string, error) {
	if this.idx == "" {
		return nil, errors.New(ERR_PARAMES)
	}
	return this.readDB(map[string]interface{}{this.idx: k})
}

func (this *BaseOB) readDB(key map[string]interface{}) (map[string]string, error) {
	this.key = key
	if this.table == "" || this.field == "" {
		return nil, errors.New(ERR_PARAMES)
	}
	var strKey string
	for k, v := range this.key {
		strKey = fmt.Sprint(k, "=", v)
	}
	rss, err := sqldb.Querys(this.table, this.field, strKey)
	if err != nil {
		this.isNew = true
		return nil, errors.New(ERR_SQLERR)
	}
	if len(rss) < 1 {
		this.isNew = true
		return nil, errors.New(ERR_NULL)
	} else {
		this.isNew = false
		return rss[0], nil
	}
}

func (this *BaseOB) Save() error {
	if len(this.table) < 1 || this.field == "" || this.idx == "" {
		return nil
	}
	if this.info == nil {
		return nil
	}
	saveMap := make(map[string]interface{})
	for k, v := range this.info {
		saveMap[k] = v
	}
	var err error
	if this.isNew {
		res, err := sqldb.Insert(this.table, saveMap)
		if err != nil {
			return err
		}
		this.lastAutoID, _ = res.LastInsertId()
		this.isNew = false
		fmt.Println("DB新增保存", this.info)
	} else {
		delete(saveMap, this.idx)
		_, err = sqldb.Update(this.table, saveMap, this.key)
		fmt.Println("DB更新保存", this.info)
	}
	return err
}

func (this *BaseOB) GetLastAutoID() int64 {
	return this.lastAutoID
}

func (this *BaseOB) SetNew(blnNew bool) {
	this.isNew = blnNew
}

func (this *BaseOB) IsNew() bool {
	return this.isNew
}

func (this *BaseOB) SetInfo(saveInfo map[string]interface{}) {
	this.info = saveInfo
}

// 设定DB关联数据
func (this *BaseOB) SetDBInfo(table, field, idx string, key map[string]interface{}) {
	this.table = table
	this.field = field
	this.idx = idx
	this.SetKey(key)
}

// 设定KEY
func (this *BaseOB) SetKey(key map[string]interface{}) {
	this.key = key
}
