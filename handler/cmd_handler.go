package handler

import (
	"fmt"
	"log"

	"github.com/ScottAI/go-sample-es-cqrs/initial"
	"github.com/ScottAI/go-sample-es-cqrs/internal/common"
)

//CommandFunc 命令信息的处理函数定义
type CommandFunc func(*common.CommandMessage, chan<- *common.EventMessage) error

//CommandHandler 处理命令信息的接口
type CommandHandler interface {
	RegisterCommand(string, ...CommandFunc) error
	HandleCommandMessage(*common.CommandMessage) error
	Start()
}

//DefaultCommandHandler 默认的信息处理结构体
type DefaultCommandHandler struct {
	commands  map[string][]CommandFunc
	eventChan chan *common.EventMessage
}

//RegisterCommand 把函数和命令绑定在一起
func (d *DefaultCommandHandler) RegisterCommand(cmd string, handlers ...CommandFunc) error {
	if _, exists := d.commands[cmd]; exists {
		return fmt.Errorf("Command: %s already exists", cmd)
	}
	d.commands[cmd] = handlers
	return nil
}

//HandleCommandMessage 处理 common.CommandMessage 并且将其传递给注册处理函数
func (d *DefaultCommandHandler) HandleCommandMessage(message *common.CommandMessage) error {
	log.Printf("Received command: %s", message.Name)
	if handlers, exists := d.commands[message.Name]; exists {
		var err error
		for _, handler := range handlers {
			err = handler(message, d.eventChan)
			if err != nil {
				break
			}
		}
		return err
	}
	return fmt.Errorf("No such command %s", message.Name)
}

//Start 启动简单evenChan
func (d *DefaultCommandHandler) Start() {
	for {
		select {
		case event := <-d.eventChan:
			initial.EventHandler.Write(event)
		}
	}
}

//NewDefaultCommandHandler 创建并返回 DefaultCommandHandler
func NewDefaultCommandHandler() *DefaultCommandHandler {
	return &DefaultCommandHandler{
		commands:  make(map[string][]CommandFunc),
		eventChan: make(chan *common.EventMessage),
	}
}
