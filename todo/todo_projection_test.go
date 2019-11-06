package todo

import (
	"encoding/json"
	"testing"

	"github.com/pborman/uuid"

	"github.com/ScottAI/go-sample-es-cqrs/common"
	"github.com/ScottAI/go-sample-es-cqrs/event"
)

func TestCreateTodo(t *testing.T) {
	bus := event.NewDefaultBus()
	projection := NewProjection(bus)
	id := uuid.New()
	data, _ := json.Marshal(&Todo{ID: id})
	raw := json.RawMessage(data)
	e := &common.EventMessage{
		Name:    eventTodoItemCreated,
		Data:    &raw,
		Version: 1,
	}
	projection.HandleEvent(e)
}

func TestCreateAndRemoveTodo(t *testing.T) {
	bus := event.NewDefaultBus()
	projection := NewProjection(bus)
	id := uuid.New()
	data, _ := json.Marshal(&Todo{ID: id})
	raw := json.RawMessage(data)
	e := &common.EventMessage{
		Name:    eventTodoItemCreated,
		Data:    &raw,
		Version: 1,
	}
	projection.HandleEvent(e)

	raw = json.RawMessage(id)
	e = &common.EventMessage{
		Name: eventTodoItemRemoved,
		Data: &raw,
	}
}
