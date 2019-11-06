package main

import (
	"encoding/json"
	"github.com/ScottAI/go-sample-es-cqrs/handler"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"path/filepath"

	"golang.org/x/net/websocket"

	"github.com/ScottAI/go-sample-es-cqrs/ws"

	"github.com/ScottAI/go-sample-es-cqrs/initial"
	"github.com/ScottAI/go-sample-es-cqrs/internal/common"
	"github.com/ScottAI/go-sample-es-cqrs/todo"
)



func main() {
	cmdHandler := handler.NewDefaultCommandHandler()
	cmdHandler.RegisterCommand("createTodoItem", todo.CreateTodoItem)
	cmdHandler.RegisterCommand("removeTodoItem", todo.RemoveTodoItem)
	cmdHandler.RegisterCommand("updateTodoItem", todo.UpdateTodoItem)
	go cmdHandler.Start()
	go initial.EventBus.Start()

	//读取事件日志 event log
	go func() {
		log.Println("Reading the event log...")
		err := initial.EventHandler.Read()
		if err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/cmd/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			cmd := new(common.CommandMessage)
			err = json.Unmarshal(data, cmd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if err := cmdHandler.HandleCommandMessage(cmd); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	http.Handle("/api/", http.FileServer(http.Dir(initial.StaticPath)))
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(initial.StaticPath, "app"))))
	http.Handle("/ws/", websocket.Handler(wsHandler))

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		log.Printf("Listening on interface: %s", addr.String())
	}

	log.Println("Listening on port 8787")
	err := http.ListenAndServe(":8787", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func wsHandler(conn *websocket.Conn) {
	log.Println("New WS client")
	ws := ws.NewClient(conn, initial.EventBus)
	ws.Listen()
}
