package common

import (
	"encoding/json"
)

//CommandMessage websocket的命令信息
type CommandMessage struct {
	Name string           `json:"name"`
	Data *json.RawMessage `json:"data"`
}

//ErrorMessage 是在websocket通信的错误信息
type ErrorMessage struct {
	Reason string `json:"reason"`
}

//EventMessage 是websocket的事件信息
type EventMessage struct {
	Name    string           `json:"name"`
	Data    *json.RawMessage `json:"data"`
	Version int              `json:"version"`
}
