package gameroom

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"tribe/gameroom/item"
	"tribe/inte"
	"vava6/vaconn"
)

// 游戏世界
type TribeWord struct {
	gameName       string
	gameVer        string
	gID            int
	lkGID          *sync.Mutex
	obManagePlayer *ManagePlayer
	CHClose        chan bool
	isRunStop      bool
}

func NewTribeWord() *TribeWord {
	ob := &TribeWord{
		gameName:       "部落世界",
		gameVer:        "0.1.20160818",
		gID:            0,
		lkGID:          new(sync.Mutex),
		obManagePlayer: NewManagePlayer(),
		CHClose:        make(chan bool),
		isRunStop:      false,
	}
	// 游戏构造时需要进进一些初始化
	// 初始化ZONE
	// OBManageZone.InitZones()
	// 将缓存对象放到缓存对象管理器里
	OBManageCache.PutCache(item.OBManageItem)
	OBManageCache.PutCache(OBManageLead)
	OBManageCache.PutCache(OBManageExplore)
	OBManageCache.PutCache(OBManageZone)
	return ob
}

func (this *TribeWord) GetTribeWordInfo() string {
	return fmt.Sprintf("%s - %s", this.gameName, this.gameVer)
}

func (this *TribeWord) getGID() int {
	this.lkGID.Lock()
	this.gID++
	this.lkGID.Unlock()
	return this.gID
}

// 延迟关闭服务器
//	@parames
//		sec		等待关闭的秒数
func (this *TribeWord) LastStop(sec int) error {
	this.lkGID.Lock()
	defer this.lkGID.Unlock()
	var err error
	if this.isRunStop == true {
		err = errors.New("正在运运行迟关闭服务器")
		return err
	}
	go func() {
		this.BordCastAllInfo("服务器将在" + strconv.Itoa(sec) + "秒后关闭")
		r := sec - 10
		h := 0
		if r > 0 {
			h = 10
			time.Sleep(time.Second * time.Duration(sec))
		} else {
			h = sec
		}

		for i := 0; i < h; i++ {
			this.BordCastAllInfo(strconv.Itoa(h-i) + " 秒后关闭服务器")
			time.Sleep(time.Second * 1)
		}
		this.Stop()
	}()
	this.isRunStop = true
	return err
}

func (this *TribeWord) BordCastAllInfo(msg string) {
	this.BordCastInfo(msg)
	this.BordCastNologinInfo(msg)
}

func (this *TribeWord) BordCastInfo(msg string) {
	obRes := inte.NewResMessage("SYSTEMINFO")
	obRes.SetInfo(msg)
	str := obRes.GetString()
	this.obManagePlayer.ChSendInfo <- str
}

func (this *TribeWord) BordCastNologinInfo(msg string) {
	obRes := inte.NewResMessage("SYSTEMINFO")
	obRes.SetInfo(msg)
	str := obRes.GetString()
	this.obManagePlayer.ChSendNoLogin <- str
}

func (this *TribeWord) BordCastChat(msg string) {
	this.obManagePlayer.ChSendInfo <- msg
}

// 向游戏世界发送停止信号
func (this *TribeWord) Stop() {
	this.CHClose <- true
}

// 创建新连接玩家
//	@parames
//		conn 	连接对象
//	@return
//		Player
func (this *TribeWord) CreateNewLinkPlayer(conn vaconn.MConn) *Player {
	gid := this.getGID()
	tUID := fmt.Sprint("p_", gid)
	newLinkPlayer := NewPlayer(gid, tUID, "USER", conn)
	// 添加到连接但未登入戏
	this.obManagePlayer.chAddNotLogin <- newLinkPlayer
	return newLinkPlayer
}

// 移除玩家，将玩家从游戏世界里移除
//	@parames
//		p	*Player		玩家指针对象
func (this *TribeWord) MovePlayer(p *Player) {
	if p.isLogin {
		this.obManagePlayer.chMoveLogin <- p
	} else {
		this.obManagePlayer.chMoveNotLogin <- p
	}
	p.CloseAndSave()
}

// 向已登入游戏玩家广播信息
//	@parames
//		msg		string	消息
func (this *TribeWord) CastBoardLogin(msg string) {
	obRes := inte.NewResMessage("SYSTEMINFO")
	obRes.SetInfo(msg)
	this.obManagePlayer.ChSendInfo <- obRes.GetString()
}

// 玩家Login游戏
// 	@parames
//		p		Player	玩家对象
//		id		string	玩家唯一ID
func (this *TribeWord) LoginPlayer(p *Player) {
	if p.isLogin {
		return
	}
	// 广播登入的玩家信息
	this.CastBoardLogin(fmt.Sprint(p.Name(), "-登入", this.gameName))
	// 判断当前玩家列表是否有这名玩家存在
	oldp, err := this.obManagePlayer.GetLoginPlayerInUID(p.uid)
	if err == nil {
		oldp.SetOtherLogin()
		oldp.Send(inte.ResMessageCMD("您从其它地方登入游戏世界"))
		oldp.CONN().Close()
	}
	// 标记已被登入
	p.SetLogin()
	// 获得角色信息
	iptLead := OBManageLead.GetLead(p.id)
	iptLead.SetIsOnline()
	this.obManagePlayer.chAddLogin <- p
	this.obManagePlayer.chMoveNotLogin <- p
	// 临时增加姓名
	if iptLead.isNew {
		iptLead.name = p.name
	}
}

// 查看游戏世界里的玩家信息列表
//	@parames
//		nil
//	@return
//		string		玩家在线信息字符
func (this *TribeWord) WordPlayers() string {
	str := ""
	for _, p := range this.obManagePlayer.loginPlayer {
		str += p.Name() + ","
	}
	strNotLogin := ""
	i := 0
	for _, p := range this.obManagePlayer.notLoginPlayer {
		strNotLogin += p.Name() + ","
		i++
		if i > 12 {
			strNotLogin += ",更多..."
			break
		}
	}
	str += "\n======================================================\nCount :" + strconv.Itoa(this.obManagePlayer.CountLogin())
	strNotLogin += "\n======================================================\nNotLogin :" + strconv.Itoa(this.obManagePlayer.CountNotlogin())
	return str + "\n\n" + strNotLogin
}

// 获取游戏名称
//	@return
//		string 		游戏名称
func (this *TribeWord) GetGameName() string {
	return this.gameName
}
