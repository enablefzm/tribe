package cmds

import (
	"strings"
	"sync"
	"tribe/gameroom"
)

var mapCmd = make(map[string]func(*gameroom.TribeWord, *gameroom.Player, string, []string) string)
var lkReg = new(sync.RWMutex)

func regCMD(cmdName string, cmdFunc func(*gameroom.TribeWord, *gameroom.Player, string, []string) string) {
	lkReg.Lock()
	mapCmd[cmdName] = cmdFunc
	lkReg.Unlock()
}

func Do(cmd string, p *gameroom.Player, g *gameroom.TribeWord) string {
	arrc := strings.Split(cmd, " ")
	c := arrc[0]
	if p.IsLogin() != true {
		if c != "login" && c != "create" {
			return "你还未登入，请先login"
		}
	}
	cmdFunc, ok := mapCmd[c]
	if ok != true {
		return "命令不存在"
	}
	return cmdFunc(g, p, cmd, arrc[1:])
}
