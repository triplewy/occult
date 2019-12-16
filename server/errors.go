package main

import "errors"

// Public errors
var (
	ErrKeyNotFound      = errors.New("occult: Key not found")
	ErrNotLeader        = errors.New("occult: Not leader")
	ErrReconfigDisabled = errors.New("attempts to perform a reconfiguration operation when reconfiguration feature is disabled")
	ErrBadArguments     = errors.New("invalid arguments")
)
