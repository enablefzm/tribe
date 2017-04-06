package gameroom

import (
	"strconv"
	"vava6/vaconn"
	"vava6/valog"
)

// 通迅接口
type IFRes interface {
	GetString() string
}

// 玩家对象
type Player struct {
	id           int    // 玩家ID
	uid          string // 玩家UID
	name         string // 玩家名称
	pass         string // 密码
	cmdTime      int    // 记录最后一次请求时间
	isLogin      bool   // 是否已登入
	isOtherLogin bool   // 是否被第三方登入替代
	Conn         vaconn.MConn
}

func NewPlayer(tID int, tUID, tName string, conn vaconn.MConn) *Player {
	newPlayer := &Player{
		id:           tID,
		uid:          tUID,
		name:         tName,
		Conn:         conn,
		isLogin:      false,
		isOtherLogin: false,
	}
	return newPlayer
}

// 关闭连接并保存退出
func (this *Player) CloseAndSave() {
	// 断开连接
	this.Conn.Close()
	if !this.isOtherLogin && this.IsLogin() {
		// 保存角色信息
		obLead, err := OBManageLead.GetLeadNoCreate(this.id)
		if err == nil {
			// ***************************************************************
			// 测试里使用退出就保存，如果正式则要不在这里使用保存对象，让缓存管理器去保存对象
			// obLead.Save()
			// ****************************************************************
			obLead.SetIsDown()
			valog.OBLog.LogMessage(this.Conn.GetIPInfo() + " " + this.name + " 保存退出游戏")
		} else {
			valog.OBLog.LogMessage(this.Conn.GetIPInfo() + " " + this.name + " 没有这名玩家Lead,退出游戏")
		}
	} else {
		valog.OBLog.LogMessage(this.Conn.GetIPInfo() + " " + this.name + " 退出游戏")
	}
}

func (this *Player) CONN() vaconn.MConn {
	return this.Conn
}

func (this *Player) GetID() int {
	return this.id
}

func (this *Player) Name() string {
	return this.name
}

func (this *Player) Uid() string {
	return this.uid
}

func (this *Player) IsLogin() bool {
	return this.isLogin
}

func (this *Player) IsOtherLogin() bool {
	return this.isOtherLogin
}

func (this *Player) Send(msg string) error {
	err := this.Conn.Send(msg)
	return err
}

func (this *Player) SendRes(res IFRes) error {
	err := this.Conn.Send(res.GetString())
	return err
}

// 通Map加载信息
func (this *Player) InitMap(rs map[string]string) {
	this.id, _ = strconv.Atoi(rs["id"])
	this.uid = rs["uid"]
	this.name = rs["name"]
}

// 设定被登入
func (this *Player) SetLogin() {
	this.isLogin = true
}

// 设定被第三方登入
func (this *Player) SetOtherLogin() {
	this.isOtherLogin = true
}

func (this *Player) GetLead() *Lead {
	return OBManageLead.GetLead(this.id)
}
