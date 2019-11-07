package todo

import (
	"encoding/json"

	"github.com/pborman/uuid"

	"github.com/ScottAI/go-sample-es-cqrs/internal/common"
)

const (
	eventTodoItemCreated = "todoItemCreated"
	eventTodoItemRemoved = "todoItemRemoved"
	eventTodoItemUpdated = "todoItemUpdated"
)

//CreateTodoItem 根据命令信息创建待办事项
func CreateTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	var todo Todo
	if err := json.Unmarshal(*cmd.Data, &todo); err != nil {
		return err
	}
	todo.ID = uuid.New()

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	raw := json.RawMessage(data)

	event := &common.EventMessage{
		Name: eventTodoItemCreated,
		Data: &raw,
	}
	eventChan <- event
	return nil
}

//RemoveTodoItem 根据命令信息删除待办事项
func RemoveTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	event := &common.EventMessage{
		Name: eventTodoItemRemoved,
		Data: cmd.Data,
	}
	eventChan <- event
	return nil
}

//UpdateTodoItem 根据命令信息更新代办事项
func UpdateTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	event := &common.EventMessage{
		Name: eventTodoItemUpdated,
		Data: cmd.Data,
	}
	eventChan <- event
	return nil
}
