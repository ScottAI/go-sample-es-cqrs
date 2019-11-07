package jsstore

import "os"

//JSStore ...
type JSStore interface {
	Set(string, interface{}) error
	Remove(string) error
	Get(string, interface{}) error
	AddToCollection(string, string, interface{}) error
	RemoveFromCollection(string, string) error
	GetCollection(string) (map[string]interface{}, error)
	RemoveAll() error
	GetDataDir() string
}

//数据存储的根目录
var DataDir = os.TempDir()
