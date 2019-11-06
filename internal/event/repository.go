package event

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/ScottAI/go-sample-es-cqrs/internal/common"
)

//负责处理事件日志，读事件日志，写事件日志
type EventHandler interface {
	Write(*common.EventMessage) error
	Read() error
}

//DefaultEventHandler 是EventHandler的一个实现
type DefaultEventHandler struct {
	r        io.Reader
	w        io.Writer
	eventBus Bus
}

//NewDefaultRepository instantiates a new DefaultRepository
func NewDefaultRepository(r io.Reader, w io.Writer, bus Bus) *DefaultEventHandler {
	return &DefaultEventHandler{
		r:        r,
		w:        w,
		eventBus: bus,
	}
}

func (d *DefaultEventHandler) Write(event *common.EventMessage) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := d.w.Write(append(jsonEvent, '\n')); err != nil {
		return err
	}

	json.Unmarshal(jsonEvent, event)
	d.eventBus.Notify(event)
	return nil
}

func (d *DefaultEventHandler) Read() error {
	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		event := &common.EventMessage{}
		if err := json.Unmarshal(scanner.Bytes(), event); err != nil {
			return err
		}
		log.Printf("Event: %s, version: %d", event.Name, event.Version)
		d.eventBus.Notify(event)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
