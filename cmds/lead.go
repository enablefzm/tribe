package cmds

import (
	"tribe/gameroom"
	"tribe/inte"
	"vava6/vatools"
)

func lead(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	if len(args) < 1 {
		return "缺少参数"
	}
	switch args[0] {
	// 显示玩家信息
	case "info":
		var id int
		if len(args) > 1 {
			id = vatools.SInt(args[1])
		} else {
			id = p.GetID()
		}
		obRes := inte.NewResMessage("LeadInfo")
		// 获取玩家角色信息
		obLead, err := gameroom.OBManageLead.GetLeadNoCreate(id)
		if err != nil {
			obRes.SetRes(false, "不存在指定的对象")
		} else {
			if obLead == nil {
				obRes.SetRes(false, "指针指向nil")
			} else {
				info := obLead.GetFieldInfo()
				// 获取玩家信息
				info["userInfo"] = map[string]interface{}{
					"userName": p.Name(),
					"uid":      p.Uid(),
				}
				obRes.SetInfo(info)
			}
		}
		p.Send(obRes.GetString())
	// 获取玩家已开启的Zone信息
	case "zones":
		obRes := inte.NewResMessage("LeadZones")
		obLead, err := gameroom.OBManageLead.GetLeadNoCreate(p.GetID())
		if err != nil {
			obRes.SetRes(false, "不存在这名玩家信息")
		} else {
			if obLead == nil {
				obRes.SetRes(false, "玩家指针为null")
			} else {
				// 获取玩家的Zone信息
				ptZones := obLead.GetZones()
				obRes.SetInfo(ptZones.GetZonesInfo())
			}
		}
		p.SendRes(obRes)
	default:
		return "没有找到相关的参数"
	}
	return ""
}

func init() {
	regCMD("lead", lead)
}
