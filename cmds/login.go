package cmds

import (
	"tribe/action"
	"tribe/gameroom"
	"tribe/inte"
)

func login(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	if len(args) < 2 {
		return "缺少参数"
	}
	if p.IsLogin() {
		return "你已经登入游戏了"
	}
	uid := args[0]
	pwd := args[1]
	iRes, rs := action.CheckLogin(uid, pwd)
	result := inte.NewResMessage("Login")
	if iRes > 0 {
		switch iRes {
		case 1:
			result.SetRes(false, "玩家UID不存在")
		case 2:
			result.SetRes(false, "密码不正确")
		case 3:
			result.SetRes(false, "查询错误")
		default:
			result.SetRes(false, "未知错误")
		}
	} else {
		result.SetRes(true, "登入成功")
		p.InitMap(rs)
		g.LoginPlayer(p)
	}
	p.Send(result.GetString())
	return ""
}

func init() {
	regCMD("login", login)
}
