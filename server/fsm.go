package main

import (
	"math/rand"
	"time"
)

// FSM is interface for finite state machine
type FSM interface {
	// Apply first appends log into logStore, then updates stableStore, then returns error
	Apply(*Entry) error

	// Read returns the value and deps for a given key
	Read(key string) ([]byte, uint64, error)
}

// Log is applied to logStore
type Log struct {
	Shardstamp uint64
	Data       []byte
}

// Entry is used to store writes and consists of a key, value, and causal dependency
type Entry struct {
	Key   string
	Value []byte
	Deps  uint64

	resChan chan uint64
	errChan chan error
}

// NewEntry creates an entry an initializes response and error channels
func NewEntry(key string, value []byte, deps uint64) *Entry {
	entry := &Entry{
		Key:   key,
		Value: value,
		Deps:  deps,

		resChan: make(chan uint64, 1),
		errChan: make(chan error, 1),
	}
	return entry
}

// Response is a blocking function that will return result of apply
func (entry *Entry) Response() (uint64, error) {
	select {
	case shardstamp := <-entry.resChan:
		return shardstamp, nil
	case err := <-entry.errChan:
		return 0, err
	}
}

// InmemFSM is in-memory implementation of FSM interface
type InmemFSM struct {
	entries map[string]*Entry
	logs    []*Log
}

func NewInmemFSM() *InmemFSM {
	return &InmemFSM{
		entries: make(map[string]*Entry),
		logs:    []*Log{},
	}
}

func (fsm *InmemFSM) Apply(entry *Entry) error {
	// Random sleep to simulate network latency
	time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
	// Store log

	// Store entry
	fsm.entries[entry.Key] = entry
	return nil
}

func (fsm *InmemFSM) Read(key string) (value []byte, deps uint64, err error) {
	if entry, ok := fsm.entries[key]; ok {
		return entry.Value, entry.Deps, nil
	}
	return nil, 0, ErrKeyNotFound
}
