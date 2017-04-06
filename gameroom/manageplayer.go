package gameroom

import (
	"errors"
	"sync"
	"vava6/vaconn"
)

// 玩家管理
type ManagePlayer struct {
	notLoginPlayer map[vaconn.MConn]*Player
	loginPlayer    map[vaconn.MConn]*Player
	chAddNotLogin  chan *Player
	chMoveNotLogin chan *Player
	chAddLogin     chan *Player
	chMoveLogin    chan *Player
	ChSendInfo     chan string
	ChSendNoLogin  chan string
	LkLogin        *sync.RWMutex
	LKNotLogin     *sync.RWMutex
}

func NewManagePlayer() *ManagePlayer {
	ob := &ManagePlayer{
		notLoginPlayer: make(map[vaconn.MConn]*Player),
		loginPlayer:    make(map[vaconn.MConn]*Player),
		chAddNotLogin:  make(chan *Player, 100),
		chAddLogin:     make(chan *Player, 100),
		chMoveLogin:    make(chan *Player, 100),
		chMoveNotLogin: make(chan *Player, 100),
		ChSendInfo:     make(chan string, 100),
		ChSendNoLogin:  make(chan string, 100),
		LkLogin:        new(sync.RWMutex),
		LKNotLogin:     new(sync.RWMutex),
	}
	go ob.run()
	return ob
}

func (this *ManagePlayer) CountLogin() int {
	return len(this.loginPlayer)
}

func (this *ManagePlayer) CountNotlogin() int {
	return len(this.notLoginPlayer)
}

func (this *ManagePlayer) run() {
	for {
		select {
		// 加入连接未登入验证玩家
		case chp := <-this.chAddNotLogin:
			if _, ok := this.notLoginPlayer[chp.CONN()]; !ok {
				this.LKNotLogin.Lock()
				this.notLoginPlayer[chp.CONN()] = chp
				this.LKNotLogin.Unlock()
			}
		// 移除未登入玩家
		case chp := <-this.chMoveNotLogin:
			this.LKNotLogin.Lock()
			delete(this.notLoginPlayer, chp.CONN())
			this.LKNotLogin.Unlock()
		// 添加已登入玩家
		case chp := <-this.chAddLogin:
			if _, ok := this.loginPlayer[chp.CONN()]; !ok {
				this.LkLogin.Lock()
				this.loginPlayer[chp.CONN()] = chp
				this.LkLogin.Unlock()
			}
		// 移除已登入玩家
		case chp := <-this.chMoveLogin:
			this.LkLogin.Lock()
			delete(this.loginPlayer, chp.CONN())
			// fmt.Println("移除已登入玩家", chp.Name())
			this.LkLogin.Unlock()
		// 向所有已登入玩家发送消息通道
		case s := <-this.ChSendInfo:
			this.LkLogin.RLock()
			for _, p := range this.loginPlayer {
				_ = p.Send(s)
			}
			this.LkLogin.RUnlock()
		case s := <-this.ChSendNoLogin:
			this.LKNotLogin.RLock()
			for _, p := range this.notLoginPlayer {
				_ = p.Send(s)
			}
			this.LKNotLogin.RUnlock()
		}
	}
}

//func (this *ManagePlayer) GetLoginPlayer() map[vaconn.MConn]*Player {
//	return this.loginPlayer
//}

func (this *ManagePlayer) GetLoginPlayerInUID(uid string) (*Player, error) {
	this.LkLogin.RLock()
	for _, p := range this.loginPlayer {
		if p.uid == uid {
			this.LkLogin.RUnlock()
			return p, nil
		}
	}
	this.LkLogin.RUnlock()
	return nil, errors.New("no")
}

func (this *ManagePlayer) GetLoginPlayerInID(id int) (*Player, error) {
	this.LkLogin.RLock()
	for _, p := range this.loginPlayer {
		if p.id == id {
			this.LkLogin.RUnlock()
			return p, nil
		}
	}
	this.LkLogin.RUnlock()
	return nil, errors.New("no")
}
