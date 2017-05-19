package main

import (
	"fmt"
	"io"
	"strings"
	"tribe/cmds"
	"tribe/gameroom"
	"tribe/inte"
	"vava6/vaconn"
	"vava6/valog"
)

var chClose = make(chan bool, 1)      // 关闭通道
var obTribe = gameroom.NewTribeWord() // 游戏主房间

func main() {
	// 运行游戏世界
	Run()
}

// 收到来自客户端玩家发送的命令
//	游戏世界收到玩家发起的命令入口
//	@parames
//		player		naviinterface.IFPlayer		连接进入的玩家对象
//		strCmd		string						命令字符串
//	@return
//		nil
func doCmd(player *gameroom.Player, strCmd string) {
	res := cmds.Do(strCmd, player, obTribe)
	if len(res) < 1 {
		return
	}
	// player.Send(naviinterface.MessageMsg(res))
	player.Send(inte.ResMessageCMD(res))
}

// 新玩家连线进入处理
//	有侦听端口监听到的新连接转入都有这个函数进行处理
//	@parames
//		Conn vaconn.MConn	类型的连接
func handleNewPlayerConn(Conn vaconn.MConn) {
	defer Conn.Close()
	Conn.Send(fmt.Sprintln("欢迎来到【", obTribe.GetTribeWordInfo(), "】"))
	player := obTribe.CreateNewLinkPlayer(Conn)
	for {
		strCmd, err := Conn.Read()
		if err != nil {
			if err == io.EOF {
				valog.OBLog.LogMessage("客户端主动断开连接")
			} else {
				valog.OBLog.LogMessage("错误：" + err.Error())
			}
			obTribe.MovePlayer(player)
			return
		}
		strCmd = strings.Trim(strCmd, "\n\r")
		if len(strCmd) == 0 {
			continue
		}
		// 转交给命令对象处理
		doCmd(player, strCmd)
	}
}

// 运行游戏 开始侦听网络连接端口
func Run() {
	valog.OBLog.LogMessage(obTribe.GetTribeWordInfo())
	obWebSocket := vaconn.NewNaviConnect(vaconn.TypeWEBSocket, "2017", "", handleNewPlayerConn)
	obWebSocket.SetOnCountFunc(obTribe.WordPlayers)
	obSocket := vaconn.NewNaviConnect(vaconn.TypeSocket, "2016", "", handleNewPlayerConn)
	go obWebSocket.Listen()
	go obSocket.Listen()
	for {
		select {
		case c := <-obTribe.CHClose:
			if c == true {
				close(obTribe.CHClose)
				return
			}
		}
	}
}
