package event

type ExploreEvent struct {
	*Event       // 继续Event
	eventRnd int // 触发机率千分之几
}
