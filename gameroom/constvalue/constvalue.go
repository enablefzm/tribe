package constvalue

// 探索的动作类型
const (
	ACT_REST       = 1 // 休息
	ACT_WALK       = 2 // 行走
	ACT_HUNT       = 3 // 守猎
	ACT_TREASURE   = 4 // 寻宝
	ACT_COLLECTION = 5 // 采集
	ACT_INSIGHT    = 6 // 洞察
)

// 探索的动作字符KEY
const (
	STR_REST       = "REST"       // 休息
	STR_WALK       = "WALK"       // 行走
	STR_HUNT       = "HUNT"       // 狩猎
	STR_TREASURE   = "TREASURE"   // 寻宝
	STR_COLLECTION = "COLLECTION" // 采集
	STR_INSIGHT    = "INSIGHT"    // 洞察
)

// 通过字符KEY获取相对应的idx
func GetActionIdx(strKey string) uint8 {
	switch strKey {
	case STR_REST:
		return ACT_REST
	case STR_WALK:
		return ACT_WALK
	case STR_HUNT:
		return ACT_HUNT
	case STR_TREASURE:
		return ACT_TREASURE
	case STR_COLLECTION:
		return ACT_COLLECTION
	case STR_INSIGHT:
		return ACT_INSIGHT
	default:
		return ACT_REST
	}
}

// 通过索引Idx获得字符Key
func GetActionKey(idx uint8) string {
	switch idx {
	case ACT_REST:
		return STR_REST
	case ACT_WALK:
		return STR_WALK
	case ACT_HUNT:
		return STR_HUNT
	case ACT_TREASURE:
		return STR_TREASURE
	case ACT_COLLECTION:
		return STR_COLLECTION
	case ACT_INSIGHT:
		return STR_INSIGHT
	default:
		return STR_REST
	}
}
