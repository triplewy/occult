package main

import (
	"math/rand"
	"time"

	pb "github.com/triplewy/occult/proto"
)

// Log is applied to logStore
type Log struct {
	Command pb.Command
	Key     string
	Value   []byte
	Deps    uint64
}

// Entry is used to store writes and consists of a key, value, and causal dependency
type Entry struct {
	Key   string
	Value []byte
	Deps  uint64
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

func (fsm *InmemFSM) Apply(log *Log) error {
	// Random sleep to simulate network latency
	time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
	// Append log
	fsm.logs = append(fsm.logs, log)
	// Store entry
	entry := &Entry{
		Key:   log.Key,
		Value: log.Value,
		Deps:  log.Deps,
	}
	switch log.Command {
	case pb.Command_Insert, pb.Command_Update:
		fsm.entries[entry.Key] = entry
	case pb.Command_Delete:
		delete(fsm.entries, entry.Key)
	default:
		return ErrUnknownCommand
	}
	return nil
}

func (fsm *InmemFSM) Read(key string) (value []byte, deps uint64, err error) {
	if entry, ok := fsm.entries[key]; ok {
		return entry.Value, entry.Deps, nil
	}
	return nil, 0, ErrKeyNotFound
}
