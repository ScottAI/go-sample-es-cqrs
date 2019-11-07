package event

import "github.com/ScottAI/go-sample-es-cqrs/internal/common"

// 订阅处理来自事件总线的多个事件中的一个事件
type Subscription struct {
	Name      string
	EventChan chan *common.EventMessage
	eventType map[string]bool
	destroyed bool
}

// 将当前订阅的事件更改为一个新事件
func (s *Subscription) ChangeSubscription(eventTypes ...string) {
	newMap := make(map[string]bool)
	for _, eventType := range eventTypes {
		newMap[eventType] = true
	}
	s.eventType = newMap
}

// 销毁该订阅以将其删除并且关闭事件频道
func (s *Subscription) Destroy() {
	if !s.destroyed {
		s.destroyed = true
		close(s.EventChan)
	}
}
