package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/samuel/go-zookeeper/zk"
	pb "github.com/triplewy/occult/proto"
	"google.golang.org/grpc"
)

// NodeState represents one of four possible states a Raft node can be in.
type NodeState uint8

// Possible states
const (
	Leader NodeState = iota
	Follower
)

type stateFunction func() stateFunction

// RemoteNode is used for RPC purposes
type RemoteNode struct {
	Addr string
}

// Node contains connection to zookeeper
type Node struct {
	Config *Config
	State  NodeState

	server *grpc.Server
	fsm    FSM

	addr  string
	nodes map[string]*RemoteNode

	shardstamp     uint64
	shardstampLock sync.RWMutex

	zkSeq int
	zk    *zk.Conn
	watch <-chan zk.Event

	applyChan chan *Entry
}

// CreateNode creates a new node
func CreateNode(config *Config) (*Node, error) {
	node := &Node{
		Config:    config,
		fsm:       NewInmemFSM(),
		addr:      "localhost:" + strconv.Itoa(config.port),
		nodes:     make(map[string]*RemoteNode),
		applyChan: make(chan *Entry),
	}

	err := node.setupRPC()
	if err != nil {
		return nil, err
	}

	conn, err := connectZk()
	if err != nil {
		return nil, err
	}
	node.zk = conn
	path, err := node.zk.CreateProtectedEphemeralSequential("/occult/node", []byte(node.addr), zk.WorldACL(zk.PermAll))
	if err != nil {
		return nil, err
	}
	node.zkSeq, err = strconv.Atoi(path[len(path)-10:])
	if err != nil {
		return nil, err
	}

	err = node.checkZk()
	if err != nil {
		return nil, err
	}

	go node.run()

	return node, nil
}

func (node *Node) setupRPC() error {
	addr := fmt.Sprintf(":%d", node.Config.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Start RPC server
	serverOptions := []grpc.ServerOption{}
	node.server = grpc.NewServer(serverOptions...)
	pb.RegisterOccultServer(node.server, node)

	go func() {
		if err := node.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return nil
}

func (node *Node) run() {
	curr := node.runFollower
	if node.State == Leader {
		curr = node.runLeader
	}
	for curr != nil {
		curr = curr()
	}
}

func (node *Node) apply(entry *Entry) (shardstamp uint64, err error) {
	node.applyChan <- entry
	return entry.Response()
}
