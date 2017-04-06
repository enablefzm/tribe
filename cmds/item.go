package cmds

import (
	"tribe/gameroom"
	CItem "tribe/gameroom/item"
	"tribe/inte"
)

func item(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	// 生成一个新道具
	obItem, _ := CItem.NewItems(103, 1)
	res := inte.NewResMessage("itemInfo")
	res.SetInfo(obItem.GetFieldInfo())
	// p.Send(res.GetString())
	return res.GetString()
}

func init() {
	regCMD("item", item)
}
