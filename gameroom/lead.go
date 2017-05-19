// 角色
package gameroom

import (
	"errors"
	"fmt"
	"sync"
	"tribe/gameroom/item"
	"tribe/sqldb"
)

func NewLead(id int) *Lead {
	ob, err := newLeadInDB(id)
	if err != nil {
		// 数据库里没有创建一个新对象
		ob = newCreateLead(id)
		ob.name = ""
		ob.isNew = true
		// 构造玩家已开启的Zone信息
		ob.obZones = &LeadZones{
			zones: []*LeadZone{&LeadZone{id: 1001}},
		}
		fmt.Println("Lead 数据库里没有，新建玩家角色", id)
	}
	return ob
}

// 从数据库里加载
func newLeadInDB(id int) (*Lead, error) {
	// 通过数据库加载
	rss, _ := sqldb.Querys("u_lead", "*", fmt.Sprintf("id=%d", id))
	if len(rss) < 1 {
		fmt.Println("没有找到这个玩家信息")
		return nil, errors.New("NO")
	}
	ob := newCreateLead(id)
	rs := rss[0]
	ob.name = rs["name"]
	ob.isNew = false
	ob.obZones = NewLeadZonesOnStr(rs["zones"])
	fmt.Println("Lead 从数据库里检索玩家", id)
	return ob, nil
}

// 创建一个新的Lead对象
func newCreateLead(id int) *Lead {
	ob := &Lead{
		lastTime: 0,
		id:       id,
		lk:       new(sync.Mutex),
		isOnline: false,
	}
	ob._init()
	return ob
}

type Lead struct {
	id         int
	name       string
	OBResource *Resource
	lastTime   int
	isNew      bool
	obExplores *ExploreQueues
	obBag      *item.Bag
	obZones    *LeadZones
	lk         *sync.Mutex
	isOnline   bool
}

func (this *Lead) _init() {
	this.OBResource = NewResource(this.id)
}

func (this *Lead) GetID() int {
	return this.id
}

func (this *Lead) GetName() string {
	return this.name
}

func (this *Lead) SetName(newName string) {
	this.name = newName
}

func (this *Lead) IsOnline() bool {
	return this.isOnline
}

func (this *Lead) SetIsOnline() {
	this.isOnline = true
}

func (this *Lead) SetIsDown() {
	this.isOnline = false
}

func (this *Lead) GetFieldInfo() map[string]interface{} {
	info := make(map[string]interface{}, 5)
	info["id"] = this.GetID()
	info["name"] = this.GetName()
	info["isOnline"] = this.IsOnline()
	infoResource := make(map[string]interface{}, 2)
	infoResource["conch"] = this.OBResource.Conch()
	infoResource["food"] = this.OBResource.Food()
	info["resource"] = infoResource
	info["zones"] = this.obZones.GetZonesInfo()
	return info
}

// 获取玩家的探索队列管理对象
//	@return
//		*ExploreQueues
func (this *Lead) GetExplores() *ExploreQueues {
	if this.obExplores == nil {
		this.lk.Lock()
		if this.obExplores == nil {
			this.obExplores = NewExploreQueues(this.id)
		}
		this.lk.Unlock()
	}
	return this.obExplores
}

// 获取玩家的背包对象
func (this *Lead) GetBag() *item.Bag {
	if this.obBag == nil {
		this.lk.Lock()
		if this.obBag == nil {
			var err error
			this.obBag, err = item.NewBag(this.id)
			if err != nil {
				// 这里有隐患
				fmt.Println("生成角色背包出错：", err.Error())
			}
		}
		this.lk.Unlock()
	}
	return this.obBag
}

// 获取玩家的Zones信息对象
func (this *Lead) GetZones() *LeadZones {
	return this.obZones
}

// 保存玩家所有数据
//	@return
//		nil
func (this *Lead) SaveAll() {
	_ = this.Save()
	_ = this.OBResource.Save()
	if this.obBag != nil {
		_ = this.obBag.Save()
	}
}

// 保存玩家角色本身数据
func (this *Lead) Save() error {
	info := make(map[string]interface{})
	info["name"] = this.name
	info["zones"] = this.obZones.GetSaveJson()
	if this.isNew {
		info["id"] = this.id
		_, err := sqldb.Insert("u_lead", info)
		if err != nil {
			return err
		}
		this.isNew = false
	} else {
		key := make(map[string]interface{})
		key["id"] = this.id
		_, err := sqldb.Update("u_lead", info, key)
		if err != nil {
			return err
		}
	}
	return nil
}
