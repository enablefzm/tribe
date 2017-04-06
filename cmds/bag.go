package cmds

import (
	"tribe/gameroom"
	// "tribe/gameroom/item"
	"tribe/inte"
	// "vava6/vatools"
)

func bag(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	il := len(args)
	if il < 1 {
		return "bag需要参数"
	}
	switch args[0] {
	// 获取背包信息
	case "info":
		res := inte.NewResMessage("bagInfo")
		info := make(map[string]interface{}, 2)
		obBag := p.GetLead().GetBag()
		info["id"] = obBag.GetID()
		info["value"] = obBag.GetValue()
		info["lenItems"] = obBag.LenItems()
		res.SetInfo(info)
		p.SendRes(res)
	case "items":
		res, info := inte.NewResMessageInfo("bagItems")
		pItems := p.GetLead().GetBag().GetItems()
		infoItems := make([]map[string]interface{}, 0, len(pItems))
		for _, pItem := range pItems {
			if pItem != nil {
				infoItems = append(infoItems, pItem.GetFieldInfo())
			}
		}
		info["playerItems"] = infoItems
		p.SendRes(res)
	default:
		return "你想对背包做什么？"
	}
	return ""
}

func init() {
	regCMD("bag", bag)
}
