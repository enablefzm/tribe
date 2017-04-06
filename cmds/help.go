package cmds

import (
	"tribe/gameroom"
)

func help(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	for k, _ := range mapCmd {
		p.Send(k)
	}
	return ""
}

func init() {
	regCMD("help", help)
}
