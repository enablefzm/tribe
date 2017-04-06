package cmds

import (
	"tribe/gameroom"
	"tribe/inte"
)

func chat(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	s := ""
	for i := 0; i < len(args); i++ {
		if i > 0 {
			s += " " + args[i]
		} else {
			s += args[i]
		}
	}
	il := len(s)
	if il > 0 {
		if il > 60 {
			return "你话说的太多了"
		} else {
			obRes, info := inte.NewResMessageInfo("CHAT")
			obLead := p.GetLead()
			info["name"] = obLead.GetName()
			info["chat"] = s
			info["id"] = p.GetID()
			g.BordCastChat(obRes.GetString())

		}
	} else {
		return "你想说什么？"
	}
	return ""
}

func init() {
	regCMD("chat", chat)
}
