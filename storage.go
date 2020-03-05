package kvstorage

import (
	"errors"
	"sync"
)

type Storage struct {
	storage map[string]string
	mt      *sync.RWMutex
}

// NewStorage creates simple [string -> string] storage
func NewStorage() Storage {
	return Storage{make(map[string]string), &sync.RWMutex{}}
}

// Insert adding key and value to storage.
// If key already in storage error will be returned
func (storage Storage) Insert(key, val string) error {
	storage.mt.RLock()
	if _, ok := storage.storage[key]; ok {
		storage.mt.RUnlock()
		return errors.New("Key already in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	storage.storage[key] = val
	storage.mt.Unlock()
	return nil
}

// Update is replacing value of the key in storage.
// If key dont exists in storage error will be returned
func (storage Storage) Update(key, val string) error {
	storage.mt.RLock()
	if _, ok := storage.storage[key]; !ok {
		storage.mt.RUnlock()
		return errors.New("Key not in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	storage.storage[key] = val
	storage.mt.Unlock()
	return nil
}

// Select returns value of the key in storage.
// If key dont exists in storage error will be returned
func (storage Storage) Select(key string) (string, error) {
	storage.mt.RLock()
	defer storage.mt.RUnlock()
	if v, ok := storage.storage[key]; ok {
		return v, nil
	}
	return "", errors.New("Key not in storage")
}

// Delete clears key/value pair from storage.
// If key dont exists in storage error will be returned
func (storage Storage) Delete(key string) error {
	storage.mt.RLock()
	if _, ok := storage.storage[key]; !ok {
		storage.mt.RUnlock()
		return errors.New("Key not in storage")
	}
	storage.mt.RUnlock()
	storage.mt.Lock()
	delete(storage.storage, key)
	storage.mt.Unlock()
	return nil
}
