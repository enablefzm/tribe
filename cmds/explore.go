package cmds

import (
	"tribe/gameroom"
	"tribe/inte"
	"vava6/vatools"
)

func explore(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	il := len(args)
	if il < 1 {
		return "缺少参数"
	}
	switch args[0] {
	// 创建新队列
	case "create":
		obRes, _ := inte.NewResMessageInfo("ExploreCreate")
		ptQues := p.GetLead().GetExplores()
		ptQue, err := ptQues.CreateNewExplore()
		if err != nil {
			obRes.SetRes(false, err.Error())
			p.Send(obRes.GetString())
			return ""
		}
		obRes.SetInfo(ptQue.GetMapInfo())
		p.Send(obRes.GetString())
		return ""
	// 查看指定的探索队列消息
	case "info":
		var ptQue *gameroom.ExploreQueue
		var err error
		// 没有指定哪个队列ID
		// 	则默认获取玩家队列第一个信息
		if il < 2 {
			// 玩家角色对象
			obLead := p.GetLead()
			// 获得玩家队列管理组
			ptQues := obLead.GetExplores()
			// 获得默认探索队象
			ptQue, err = ptQues.GetDefaultExploreQueue()
			if err != nil {
				return err.Error()
			}
		} else {
			ptQue, err = gameroom.OBManageExplore.GetExplore(vatools.SInt(args[1]))
			if err != nil {
				return err.Error()
			}
		}
		res := inte.NewResMessage("ExploreInfo")
		res.SetInfo(ptQue.GetMapInfo())
		p.Send(res.GetString())
		return ""
		// 获取队列所有信息
	case "infos":
		var obLead *gameroom.Lead
		res := inte.NewResMessage("ExploreInfos")
		if il < 2 {
			// 获取自己的队列信息
			obLead = p.GetLead()
		} else {
			var err error
			obLead, err = gameroom.OBManageLead.GetLeadNoCreate(vatools.SInt(args[1]))
			if err != nil {
				res.SetRes(false, err.Error())
				p.Send(res.GetString())
				return ""
			}
		}
		ptQues := obLead.GetExplores()
		arrQues := ptQues.GetExploreQueues()
		// 获取探索所有队列信息
		arrExpQue := make([]map[string]interface{}, 0, len(arrQues))
		for _, ptQue := range arrQues {
			arrExpQue = append(arrExpQue, ptQue.GetMapInfo())
		}
		res.SetInfo(arrExpQue)
		p.Send(res.GetString())
		return ""
	// 探索队加入到指定的zone
	// explore join zoneid expqueueid
	case "join":
		res := inte.NewResMessage("ExploreJoin")
		if il < 3 {
			res.SetRes(false, "没有指定要加入哪个ZONE和自己哪个探索队ID")
			p.Send(res.GetString())
			return ""
		}
		// 获得指定的Expquid
		queueID := vatools.SInt(args[2])
		ptQue, err := gameroom.OBManageExplore.GetExplore(queueID)
		if err != nil {
			res.SetRes(false, "没有你需要的探索队列")
			p.Send(res.GetString())
			return ""
		}
		// 判断探索队列是不是操作者本人的
		if ptQue.GetLeadID() != p.GetID() {
			res.SetRes(false, "这个队列不是你的")
			p.Send(res.GetString())
			return ""
		}
		// 获得Zone
		obZone, err := gameroom.OBManageZone.GetCanchZone(vatools.SInt(args[1]))
		if err != nil {
			res.SetRes(false, err.Error())
			p.Send(res.GetString())
			return ""
		}
		err = obZone.JoinExploreQueue(ptQue)
		if err != nil {
			res.SetRes(false, err.Error())
			p.Send(res.GetString())
			return ""
		}
		res.SetRes(true, "加载成功")
		p.Send(res.GetString())
		return "ExploreJoin"
	}
	return "探索"
}

func init() {
	regCMD("explore", explore)
}
