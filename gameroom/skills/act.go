package skills

// 探索队基础动作
const (
	E_REST       = iota + 1 // 休息
	E_WALK                  // 行走
	E_HUNT                  // 守猎
	E_TREASURE              // 寻宝
	E_COLLECTION            // 采集
	E_INSIGHT               // 洞察
)

// 可以执行动作技能
const (
	EA_EXPLORE = iota // 探索
)
