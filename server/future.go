package main

type ApplyFuture struct {
	log *Log

	resChan chan uint64
	errChan chan error
}

func (f *ApplyFuture) Response() (uint64, error) {
	select {
	case shardstamp := <-f.resChan:
		return shardstamp, nil
	case err := <-f.errChan:
		return 0, err
	}
}
