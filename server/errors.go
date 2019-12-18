package main

import "errors"

// Public errors
var (
	ErrKeyNotFound    = errors.New("occult: Key not found")
	ErrKeyExists      = errors.New("occult: Key already exists")
	ErrNotLeader      = errors.New("occult: Not leader")
	ErrUnknownCommand = errors.New("occult: Unknown log command")
)
