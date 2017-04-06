package cmds

import (
	"tribe/gameroom"
	CHero "tribe/gameroom/hero"
	"tribe/inte"
	//	"vava6/vatools"
)

func hero(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	arrLen := len(args)
	if arrLen < 1 {
		return "缺少参数"
	}
	switch args[0] {
	// 酒馆创建
	case "createPub":
		res := inte.NewResMessage("HeroCreatePub")
		arr := make([]map[string]interface{}, 0, 3)
		for i := 0; i < 3; i++ {
			obHero := CHero.CreatePub()
			arr = append(arr, obHero.GetMapInfo())
		}
		res.SetInfo(arr)
		p.SendRes(res)
	default:
		return "没有相应的参数"
	}
	return ""
}

func init() {
	regCMD("hero", hero)
}
