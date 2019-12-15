package main

import (
	"fmt"
	"strconv"

	"github.com/samuel/go-zookeeper/zk"
)

// Possible states
const (
	Leader uint8 = iota
	Follower
	Joining
)

type stateFunction func() stateFunction

// Node contains connection to zookeeper
type Node struct {
	Config *Config
	State  uint8

	addr     string
	znodeSeq int
	zk       *zk.Conn
	watch    <-chan zk.Event
	nodes    []string

	errChan chan error
}

// NewNode creates a new node
func NewNode(config *Config) (*Node, error) {
	node := new(Node)
	node.Config = config
	node.State = Joining
	node.addr = "localhost:" + strconv.Itoa(config.port)
	conn, err := connectZk()
	if err != nil {
		return nil, err
	}
	node.zk = conn
	path, err := node.zk.CreateProtectedEphemeralSequential("/occult/node", []byte(node.addr), zk.WorldACL(zk.PermAll))
	if err != nil {
		return nil, err
	}
	node.znodeSeq, err = strconv.Atoi(path[len(path)-8:])
	if err != nil {
		return nil, err
	}

	go node.run()

	return node, nil
}

func (node *Node) checkZk() ([]string, error) {
	znodes, _, watch, err := node.zk.ChildrenW("/occult")
	if err != nil {
		return nil, err
	}
	node.watch = watch

	for _, znode := range znodes {
		seq, err := strconv.Atoi(znode[len(znode)-8:])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(seq)
	}

	return znodes, nil
}

func (node *Node) run() {
	var curr stateFunction = node.runFollower
	for curr != nil {
		curr = curr()
	}
}
