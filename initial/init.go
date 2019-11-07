package initial

import (
	"github.com/ScottAI/go-sample-es-cqrs/internal/event"
	"github.com/ScottAI/go-sample-es-cqrs/internal/jsstore"
	"github.com/ScottAI/go-sample-es-cqrs/todo"
	"io"
	"log"
	"os"
	"path/filepath"
)

var EventBus event.Bus
var EventLogFile string
var EventLogWriter io.Writer
var EventLogReader io.Reader
var EventHandler event.EventHandler
var TodoProjection *todo.Projection
var StaticPath string

func init() {
	dir,err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	StaticPath = filepath.Join(dir, "static")
	log.Println("staticPath:",StaticPath)
	jsstore.DataDir = filepath.Join(StaticPath, "api")
	EventBus = event.NewDefaultBus()
	EventLogFile = filepath.Join(os.TempDir(), "eventlog")
	TodoProjection = todo.NewProjection(EventBus)
	EventLogWriter, _ = os.OpenFile(EventLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	EventLogReader, _ = os.Open(EventLogFile)
	EventHandler = event.NewDefaultRepository(EventLogReader, EventLogWriter, EventBus)

	log.SetFlags(log.Flags() | log.Lshortfile)
}