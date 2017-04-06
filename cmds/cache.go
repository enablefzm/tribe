package cmds

import (
	"fmt"
	"strings"
	"tribe/gameroom"
)

func cache(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	il := len(args)
	if il < 1 {
		return "cache"
	}
	switch args[0] {
	case "list":
		result := gameroom.OBManageCache.GetList()
		return strings.Join(result, " ")
	// 进行操作
	case "operate":
		if il < 2 {
			return "需要指定Cache名称"
		}
		if obCache, err := gameroom.OBManageCache.GetCache(args[1]); err != nil {
			return err.Error()
		} else {
			if il == 2 {
				return fmt.Sprint(obCache.Name(), "当前缓存里有", obCache.Count(), "对象")
			} else {
				switch args[2] {
				case "release":
					obCache.Release()
					return "成功清空：" + obCache.Name()
				}
			}
		}
	}
	return "do nothing."
}

func init() {
	regCMD("cache", cache)
}
