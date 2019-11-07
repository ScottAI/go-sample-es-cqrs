package jsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
)

//JSONStore 存储json文件和gzip文件
type JSONStore struct {
	dataDir string
	files   map[string]*jsonFile
}

//NewJSONStore 创建一个新的 JSONStore，包含有一个绝对路径
func NewJSONFSStore(dataDir string) (*JSONStore, error) {
	if dataDir[0] != '/' {
		dataDir = filepath.Join(DataDir, dataDir)
	}

	store := &JSONStore{
		dataDir: dataDir,
		files:   make(map[string]*jsonFile),
	}

	store.RemoveAll()
	return store, nil
}

//Set 设置一个含有id的json
func (j *JSONStore) Set(id string, data interface{}) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err

	}
	return file.set(data)
}

//Remove 根据指定id删除数据
func (j *JSONStore) Remove(id string) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err
	}

	file.remove()
	return nil
}

//Get 判指定id的json是否可以解析
func (j *JSONStore) Get(id string, v interface{}) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err
	}
	return file.get(v)
}

//RemoveAll 删除所有数据
func (j *JSONStore) RemoveAll() error {
	if err := os.RemoveAll(j.dataDir); err != nil {
		return err
	}

	if err := os.MkdirAll(j.dataDir, 0755); err != nil {
		return err
	}
	return nil
}

//GetDataDir 返回json数据存储的路径
func (j *JSONStore) GetDataDir() string {
	return j.dataDir
}

func (j *JSONStore) AddToCollection(cPath string, id string, v interface{}) error {
	file, err := j.getJsonFile(cPath)
	if err != nil {
		return err
	}

	c, err := j.getCollection(file)
	if err != nil {
		return err
	}
	c[id] = v

	return file.set(c)
}

func (j *JSONStore) RemoveFromCollection(cPath string, id string) error {
	file, err := j.getJsonFile(cPath)
	if err != nil {
		return err
	}

	c, err := j.getCollection(file)
	if err != nil {
		return err
	}
	delete(c, id)

	return file.set(c)
}

func (j *JSONStore) GetCollection(cPath string) (map[string]interface{}, error) {
	file, err := j.getJsonFile(cPath)

	if err != nil {
		return nil, err
	}
	return j.getCollection(file)
}

func (j *JSONStore) getCollection(file *jsonFile) (map[string]interface{}, error) {
	c := make(map[string]interface{})

	if err := file.get(&c); err != nil {
		return nil, err
	}
	return c, nil
}

func (j *JSONStore) getJsonFile(path string) (*jsonFile, error) {
	if v, ok := j.files[path]; !ok {
		file, err := os.OpenFile(
			filepath.Join(j.dataDir, fmt.Sprintf("%s.json", path)),
			os.O_RDWR|os.O_CREATE,
			0644,
		)

		if err != nil {
			file.Close()
			debug.PrintStack()
			return nil, err
		}

		f := newJsonFile(file)

		j.files[path] = f

		return f, nil
	} else {
		return v, nil
	}

}

func (j *JSONStore) Flush() {
	for _, f := range j.files {
		f.flush()
	}
}

func (j *JSONStore) Stop() {
	for _, f := range j.files {
		f.stop()
	}
}
