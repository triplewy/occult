package main

// Future is used to represent an action that may occur in the future.
type Future interface {
	// Error blocks until the future arrives and then
	// returns the error status of the future.
	// This may be called any number of times - all
	// calls will return the same value.
	// Note that it is not OK to call this method
	// twice concurrently on the same Future instance.
	Error() error
}

// ShardstampFuture is used for future actions that can result in a log entry
// being created.
type ShardstampFuture interface {
	Future

	// Shardstamp holds the shardstamp of the newly applied log entry.
	// This must not be called until after the Error method has returned.
	Shardstamp() uint64
}

// ApplyFuture is used for Apply and can return the FSM response.
type ApplyFuture interface {
	ShardstampFuture

	// Response returns the FSM response as returned
	// by the FSM.Apply method. This must not be called
	// until after the Error method has returned.
	Response() interface{}
}

// deferError can be embedded to allow a future
// to provide an error in the future.
type deferError struct {
	err       error
	errCh     chan error
	responded bool
}

// logFuture is used to apply a log entry and waits until
// the log is considered committed.
type logFuture struct {
	deferError
	log      Log
	response interface{}
}
