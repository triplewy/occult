package main

import (
	"path/filepath"
	"strconv"
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

func (node *Node) checkZk() error {
	znodes, _, watch, err := node.zk.ChildrenW("/occult")
	if err != nil {
		return err
	}
	// Set watch to children in case if leader
	node.watch = watch

	// Find greatest znode seq that is less than own seq
	next := -1
	nextNode := ""
	for _, znode := range znodes {
		seq, err := strconv.Atoi(znode[len(znode)-10:])
		if err != nil {
			return err
		}
		if seq < node.zkSeq && seq > next {
			next = seq
			nextNode = znode
		}
	}

	// If no znode seq less than own seq, then self is leader
	if next == -1 {
		node.State = Leader
		// Update local nodes map
		for _, znode := range znodes {
			if _, ok := node.nodes[znode]; !ok {
				data, _, err := node.zk.Get(filepath.Join("/occult", znode))
				if err != nil {
					return err
				}
				node.nodes[znode] = &RemoteNode{Addr: string(data)}
			}
		}
	} else {
		// Watch next node for changes. Like linkedlist
		_, _, watch, err := node.zk.GetW(filepath.Join("/occult", nextNode))
		if err != nil {
			return err
		}
		node.watch = watch
		node.State = Follower
	}
	return nil
}
