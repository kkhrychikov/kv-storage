package kvstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Storage struct {
	storage map[string]string
	mt      *sync.RWMutex
}

// NewStorage creates simple [string -> string] storage
func NewStorage() *Storage {
	return &Storage{make(map[string]string), &sync.RWMutex{}}
}

// Insert adding key and value to storage.
// If key already in storage error will be returned
func (storage *Storage) Insert(key, val string) error {
	if key == "" {
		return errors.New("Empty key")
	}
	storage.mt.Lock()
	if _, ok := storage.storage[key]; ok {
		storage.mt.Unlock()
		return errors.New("Key already in storage")
	}
	storage.storage[key] = val
	storage.mt.Unlock()
	return nil
}

// Update is replacing value of the key in storage.
// If key dont exists in storage error will be returned
func (storage *Storage) Update(key, val string) error {
	if key == "" {
		return errors.New("Empty key")
	}
	storage.mt.Lock()
	if _, ok := storage.storage[key]; !ok {
		storage.mt.Unlock()
		return errors.New("Key not in storage")
	}
	storage.storage[key] = val
	storage.mt.Unlock()
	return nil
}

// Select returns value of the key in storage.
// If key dont exists in storage error will be returned
func (storage *Storage) Select(key string) (string, error) {
	if key == "" {
		return "", errors.New("Empty key")
	}
	storage.mt.RLock()
	defer storage.mt.RUnlock()
	if v, ok := storage.storage[key]; ok {
		return v, nil
	}
	return "", errors.New("Key not in storage")
}

// Delete clears key/value pair from storage.
// If key dont exists in storage error will be returned
func (storage *Storage) Delete(key string) error {
	if key == "" {
		return errors.New("Empty key")
	}
	storage.mt.Lock()
	if _, ok := storage.storage[key]; !ok {
		storage.mt.Unlock()
		return errors.New("Key not in storage")
	}
	delete(storage.storage, key)
	storage.mt.Unlock()
	return nil
}

// Dump return storage elements in JSON format
func (storage *Storage) Dump() ([]byte, error) {
	return json.Marshal(storage.storage)
}

// Reset clears all data
func (storage *Storage) Reset() {
	storage.storage = make(map[string]string)
	return
}

// Load wipes storage and adds new json data.
// Data must be dumped earlier from storage or be {"key":"value"} format like
// {"Name":"Adam","Age":"36","Job":"CEO"}
func (storage *Storage) Load(data []byte) (stat string, err error) {
	rawMap := make(map[string]interface{})
	err = json.Unmarshal(data, &rawMap)
	if err != nil {
		return "", err
	}
	storage.Reset()
	insertedCount := 0
	for k, v := range rawMap {
		if val, ok := v.(string); ok {
			err = storage.Insert(k, val)
			if err == nil {
				insertedCount++
			}
		}
	}
	return fmt.Sprintf("Added %d/%d pairs", insertedCount, len(rawMap)), nil
}
