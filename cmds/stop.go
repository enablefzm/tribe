package cmds

import (
	"strconv"
	"tribe/gameroom"
)

func stop(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	if len(args) > 0 {
		lastTime, err := strconv.Atoi(args[0])
		if err != nil {
			return "时间不太对 " + args[0]
		} else {
			g.LastStop(lastTime)
			return "你执行了服务器延迟停止"
		}
	} else {
		g.BordCastInfo(p.Name() + " 执行了服务器停止")
		g.Stop()
		return "你执行了服务器停止"
	}
}

func init() {
	regCMD("stop", stop)
}
