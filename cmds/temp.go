package cmds

import (
	"fmt"
	"tribe/gameroom"
	"vava6/vatools"
)

func temp(g *gameroom.TribeWord, p *gameroom.Player, cmd string, args []string) string {
	dbField, err := gameroom.NewFieldDb(vatools.SInt(args[0]))
	// dbField, err := gameroom.NewFieldDb(1001)
	if err != nil {
		return err.Error()
	} else {
		fmt.Println(dbField.GetFieldInfo())
		return "OK"
	}
}

func init() {
	regCMD("temp", temp)
}
