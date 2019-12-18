package main

import (
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func (c *Client) retrieveNodes() error {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		return err
	}
	znodes, _, err := conn.Children("/occult")
	if err != nil {
		return err
	}
	nodes := []*RemoteNode{}
	smallestSeq := math.MaxInt32
	for _, znode := range znodes {
		data, _, err := conn.Get(filepath.Join("/occult", znode))
		if err != nil {
			return err
		}
		seq, err := strconv.Atoi(znode[len(znode)-10:])
		if err != nil {
			return err
		}
		remote := &RemoteNode{Addr: string(data)}
		nodes = append(nodes, remote)
		if seq < smallestSeq {
			c.leader = remote
		}
	}
	c.nodes = nodes
	return nil
}
