package world

import (
	"tribe/gameroom/event"
)

type DbField struct {
	id         int                      // 格子的ID
	name       string                   // 格子名称
	fType      uint8                    // 格子类型
	image      string                   // 图片名称
	defense    uint                     // 被对抗的值-接受探索队的攻击
	decoration string                   // 上面物件名称
	canAction  map[string]bool          // 可以接受的动作
	mpEvent    map[string]event.IFEvent // 事件列表
}
