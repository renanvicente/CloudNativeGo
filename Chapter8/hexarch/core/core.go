package core

import (
	"errors"
	"log"
	"sync"
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type Event struct {
	Sequence  uint64    // A unique record ID
	EventType EventType // The action taken
	Key       string    // The key affected by this transaction
	Value     string    // The value of a PUT the transaction
}

type EventType byte

const (
	_           = iota // iota == 0; ignore the zero value
	EventDelete = iota // iota == 1
	EventPut           // iota == 2; implicitly repeat
)

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

func (store *KeyValueStore) Restore() error {
	var err error


	events, errors := store.transact.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors: // Retrieve any errors
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = store.Delete(e.Key)
			case EventPut: // Got a PUT event!
				err = store.Put(e.Key, e.Value)

			}
		}
	}
	store.transact.Run()
	go func() {
		for err := range store.transact.Err() {
			log.Println(err)
		}
	}()
	return err

}
