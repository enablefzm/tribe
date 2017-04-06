package cmds

import (
	"tribe/gameroom"
	"tribe/inte"
	"vava6/vatools"
)

func zone(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	il := len(args)
	if il < 2 {
		return "缺少参数"
	}
	switch args[0] {
	case "info":
		zoneID := vatools.SInt(args[1])
		obRes, _ := inte.NewResMessageInfo("ZoneInfo")
		obZone, err := gameroom.OBManageZone.GetCanchZone(zoneID)
		if err == nil {
			obRes.SetInfo(obZone.GetInfo())
		} else {
			obRes.SetRes(false, "没有找到相应的Zone")
		}
		p.Send(obRes.GetString())
	}
	return ""
}

func init() {
	regCMD("zone", zone)
}
