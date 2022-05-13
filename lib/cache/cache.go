package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, t time.Duration) error
	SetObj(key string, value interface{}, t time.Duration) error
	Get(key string) (string, error)
	GetStruct(key string, obj interface{}) error
	Exists(key string) (bool, error)
	Delete(key string) error
	Keys(pattern string) ([]string, error)
}

type cacheItem struct {
	T time.Time
	V []byte
}

type MapCache struct {
	cache  map[string]cacheItem
	lock   sync.Mutex
	isSync bool

	filepath string
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewMapCache() *MapCache {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	m := &MapCache{
		cache:    make(map[string]cacheItem),
		lock:     sync.Mutex{},
		isSync:   true,
		filepath: path + "/cache.json",
	}
	m.rsync()
	go m.delete()
	return m
}

func (m *MapCache) delete() {
	for {
		time.Sleep(time.Hour)
		for k, v := range m.cache {
			if time.Now().After(v.T) {
				delete(m.cache, k)
			}
		}
	}
}

func (m *MapCache) rsync() {
	if !m.isSync {
		return
	}

	ok, err := pathExists(m.filepath)
	if err != nil {
		panic(err)
	}

	if !ok {
		return
	}

	b, err := os.ReadFile(m.filepath)
	if err != nil {
		panic(err)
	}

	// err = gob.NewDecoder(bytes.NewReader(b)).Decode(&m.cache)
	// err = binary.Unmarshal(b, &m.cache)

	err = json.Unmarshal(b, &m.cache)
	if err != nil {
		panic(err)
	}
	// buf := bytes.Buffer{}

	// err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &m.cache)
	// if err != nil {
	// 	panic(err)
	// }

}

func (m *MapCache) wsync() error {
	if !m.isSync {
		return nil
	}
	b, err := json.Marshal(m.cache)
	if err != nil {
		return err
	}
	// buf := bytes.Buffer{}
	// err := gob.NewEncoder(&buf).Encode(m.cache)
	// b ,err := binary.Marshal(m.cache)
	// if err != nil {
	// 	return err
	// }

	return os.WriteFile(m.filepath, b, 0644)
}

func (m *MapCache) Set(key string, value interface{}, t time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.cache[key] = cacheItem{T: time.Now().Add(t), V: v}
	if err := m.wsync(); err != nil {
		delete(m.cache, key)
		return err
	}
	return nil
}

func (m *MapCache) SetObj(key string, value interface{}, t time.Duration) error {
	return m.Set(key, value, t)
}

func (m *MapCache) Get(key string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	item, ok := m.cache[key]
	if !ok {
		return "", errors.New("cannot found " + key)
	}
	if time.Now().After(item.T) {
		delete(m.cache, key)
		return "", errors.New("cannot found " + key)
	}

	return string(item.V), nil
}

func (m *MapCache) GetStruct(key string, obj interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
		}
	}()
	item, ok := m.cache[key]
	if !ok {
		return errors.New("cannot found " + key)
	}
	if time.Now().After(item.T) {
		delete(m.cache, key)
		return errors.New("cannot found " + key)
	}

	// reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(item.V))

	return json.Unmarshal(item.V, obj)
}

func (m *MapCache) Exists(key string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	item, ok := m.cache[key]
	if !ok {
		return false, nil
	}
	if time.Now().After(item.T) {
		delete(m.cache, key)
		return false, nil
	}
	return true, nil
}

func (m *MapCache) Delete(key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.cache, key)
	return nil
}

func (m *MapCache) Keys(pattern string) ([]string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	ks := make([]string, 0)
	r := regexp.MustCompile(pattern)
	for k := range m.cache {
		if r.MatchString(k) {
			ks = append(ks, k)
		}
	}
	return ks, nil
}
