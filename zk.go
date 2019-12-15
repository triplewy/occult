package main

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func connectZk() (*zk.Conn, error) {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	_, err = c.Create("/occult", []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return nil, err
	}
	return c, nil
}
