package kvstorage

import (
	"encoding/json"
	"errors"
	"sync"
)

type Storage interface {
	Insert(KeyVal) (interface{}, error)
	Update(KeyVal) (interface{}, error)
	Select(KeyVal) (interface{}, error)
	Delete(KeyVal) (interface{}, error)
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte) (KeyVal, error)
}

type StringsStorage struct {
	storage map[string]string
	mt      *sync.RWMutex
}

func CreateStringsStorage() StringsStorage {
	return StringsStorage{make(map[string]string), &sync.RWMutex{}}
}

func (storage StringsStorage) Insert(kv KeyVal) (interface{}, error) {
	storage.mt.RLock()
	if _, ok := storage.storage[kv.GetKey()]; ok {
		storage.mt.RUnlock()
		return nil, errors.New("Key already in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	storage.storage[kv.GetKey()] = kv.GetVal().(string)
	storage.mt.Unlock()
	return nil, nil
}
func (storage StringsStorage) Update(kv KeyVal) (interface{}, error) {
	storage.mt.RLock()
	if _, ok := storage.storage[kv.GetKey()]; !ok {
		storage.mt.RUnlock()
		return nil, errors.New("Key not in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	storage.storage[kv.GetKey()] = kv.GetVal().(string)
	storage.mt.Unlock()
	return nil, nil
}
func (storage StringsStorage) Select(kv KeyVal) (interface{}, error) {
	storage.mt.RLock()
	defer storage.mt.RUnlock()
	if v, ok := storage.storage[kv.GetKey()]; ok {
		return v, nil
	}
	return nil, errors.New("Key not in storage")
}
func (storage StringsStorage) Delete(kv KeyVal) (interface{}, error) {
	storage.mt.RLock()
	if _, ok := storage.storage[kv.GetKey()]; !ok {
		storage.mt.RUnlock()
		return nil, errors.New("Key not in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	delete(storage.storage, kv.GetKey())
	storage.mt.Unlock()
	return nil, nil
}

func (storage StringsStorage) Marshal(data interface{}) ([]byte, error) {
	if v, ok := data.(string); ok {
		return []byte(v), nil
	}
	return nil, errors.New("Can't marshal: value is not string")
}

func (storage StringsStorage) Unmarshal(b []byte) (KeyVal, error) {
	kv := keyValString{}
	err := json.Unmarshal(b, &kv)
	if err != nil {
		return nil, err
	}
	return kv, nil
}

type KeyVal interface {
	GetKey() string
	GetVal() interface{}
}

type keyValString struct {
	Key string
	Val string
}

func (kv keyValString) GetKey() string {
	return kv.Key
}

func (kv keyValString) GetVal() interface{} {
	return kv.Val
}

type keyVal struct {
	Key string
	Val interface{}
}

func (kv keyVal) GetKey() string {
	return kv.Key
}

func (kv keyVal) GetVal() interface{} {
	return kv.Val
}
