package event

import (
	"log"
	"time"

	"github.com/ScottAI/go-sample-es-cqrs/internal/common"
)

//Bus interface/事件总线接口
type Bus interface {
	Notify(*common.EventMessage)
	Subscribe(string, ...string) *Subscription
	Start()
}

//DefaultBus implementation/事件总线的一个实现
type DefaultBus struct {
	subscriptions []*Subscription
	notifyChan    chan *common.EventMessage
}

// NewDefaultBus 创建一个默认事件总线
func NewDefaultBus() *DefaultBus {
	return &DefaultBus{
		subscriptions: make([]*Subscription, 0),
		notifyChan:    make(chan *common.EventMessage, 0),
	}

}

//监听器订阅事件总线
func (d *DefaultBus) Subscribe(name string, eventType ...string) *Subscription {
	eventTypeMap := make(map[string]bool)
	for _, v := range eventType {
		eventTypeMap[v] = true
	}

	subscription := &Subscription{
		Name:      name,
		EventChan: make(chan *common.EventMessage, 1),
		eventType: eventTypeMap,
	}
	d.subscriptions = append(d.subscriptions, subscription)
	return subscription
}

//向监听器发布事件信息
func (d *DefaultBus) Notify(event *common.EventMessage) {
	d.notifyChan <- event
}

//Start应该在一个goroutine中运行, 该方法启动事件总线
func (d *DefaultBus) Start() {
	for {
		select {
		case event := <-d.notifyChan:
			for i := len(d.subscriptions) - 1; i >= 0; i-- {
				li := i
				subscription := d.subscriptions[li]
				if subscription.destroyed {
					d.subscriptions = append(d.subscriptions[:li], d.subscriptions[li+1:]...)
				} else if len(subscription.eventType) == 0 || subscription.eventType[event.Name] {
					go func() {
						select {
						case subscription.EventChan <- event:
							log.Printf("Sending event %s to %s", event.Name, subscription.Name)
						case <-time.After(3 * time.Second):
							log.Printf("Sending event to %s timed out!", subscription.Name)
						}
					}()
				}
			}
		}
	}
}
