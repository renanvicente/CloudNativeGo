package core

import (
	"errors"
	"sync"
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
}

type KeyValueStore struct {
	sync.RWMutex
	m        map[string]string
	transact TransactionLogger
}

func NewKeyValueStore(tl TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: tl,
	}
}

var ErrorNoSuchKey = errors.New("no such key")

func (store *KeyValueStore) Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	store.transact.WritePut(key, value)
	return nil
}

func (store *KeyValueStore) Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func (store *KeyValueStore) Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	store.transact.WriteDelete(key)
	return nil
}
