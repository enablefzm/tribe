package cmds

import (
	"fmt"
	"tribe/gameroom"
	"tribe/sqldb"
)

func test(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	c := sqldb.NewSqlDbCfg()
	fmt.Println(c.User)
	fmt.Println(c.Pass)
	fmt.Println(c.DBName)
	ptFieldDB, err := gameroom.NewFieldDb(1001)
	if err != nil {
		return "错误的ID"
	}
	acts := ptFieldDB.GetActs()
	for _, v := range acts {
		et, ok := ptFieldDB.GetEvent(v)
		if !ok {
			return "不存在这个事件"
		}
		fmt.Println(et)
	}
	return ""
}

func init() {
	regCMD("test", test)
	regCMD("t", test)
}
