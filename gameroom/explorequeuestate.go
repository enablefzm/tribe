package gameroom

import (
	"tribe/gameroom/exploreaction"
)

// 探索队状态
type ExploreQueueState struct {
	nowAct     exploreaction.IFAction // 当前探索的动作
	state      uint8                  // 当前状态
	stateValue uint8                  // 当前状态的具体描述
}

// 构造当前状态对象
func NewExploreQueueState() *ExploreQueueState {
	return &ExploreQueueState{}
}

// 获取当前的动作值如果当前没有动作则返回nil
//  @return
//      exploreaction.IFAction     当前的动作对象
//      bool                       是否有动作对象
func (this *ExploreQueueState) GetNowAct() (exploreaction.IFAction, bool) {
	if this.nowAct == nil {
		return nil, false
	} else {
		return this.nowAct, true
	}
}

// 设定当前动作
func (this *ExploreQueueState) SetNowAct(act exploreaction.IFAction) {
	this.nowAct = act
}

// 设定当前状态信息
func (this *ExploreQueueState) SetState(vState, vStateValue uint8) {
	this.state = vState
	this.stateValue = vStateValue
}

// 更改状态为主动退出
func (this *ExploreQueueState) SetActiveQuit() {
	this.SetState(0, 0)
}

// 更改状态为没有粮食退出
//	状态值为 0, 1
func (this *ExploreQueueState) SetNoFoodQuit() {
	this.SetState(0, 1)
}

// 未知原因退出
func (this *ExploreQueueState) SetNoOtherQuit(vStatValue uint8) {
	this.SetState(0, vStatValue)
}
