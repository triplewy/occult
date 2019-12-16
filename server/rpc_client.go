package main

import (
	"context"
	"sync"
	"time"

	pb "github.com/triplewy/occult/proto"
	"google.golang.org/grpc"
)

var dialOptions []grpc.DialOption

func init() {
	dialOptions = []grpc.DialOption{grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithTimeout(2 * time.Second)}
}

// RPC Invocation functions

var connMap = make(map[string]*grpc.ClientConn)
var connMapLock = &sync.RWMutex{}

func closeAllConnections() {
	connMapLock.Lock()
	defer connMapLock.Unlock()
	for k, conn := range connMap {
		conn.Close()
		delete(connMap, k)
	}
}

// ClientConn creates or returns a cached RPC client for the given remote node
func (remote *RemoteNode) ClientConn() (pb.OccultClient, error) {
	connMapLock.RLock()
	if cc, ok := connMap[remote.Addr]; ok {
		connMapLock.RUnlock()
		return pb.NewOccultClient(cc), nil
	}
	connMapLock.RUnlock()

	cc, err := grpc.Dial(remote.Addr, dialOptions...)
	if err != nil {
		return nil, err
	}
	connMapLock.Lock()
	connMap[remote.Addr] = cc
	connMapLock.Unlock()

	return pb.NewOccultClient(cc), err
}

// RemoveClientConn removes the client connection to the given node, if present
func (remote *RemoteNode) RemoveClientConn() {
	connMapLock.Lock()
	defer connMapLock.Unlock()
	if cc, ok := connMap[remote.Addr]; ok {
		cc.Close()
		delete(connMap, remote.Addr)
	}
}

// connCheck checks the given error and removes the client connection if it's not nil
func (remote *RemoteNode) connCheck(err error) error {
	if err != nil {
		remote.RemoveClientConn()
	}
	return err
}

// ReplicateRPC is used by leader to replicate entries to followers
func (remote *RemoteNode) ReplicateRPC(key string, value []byte, deps uint64, shardstamp uint64) error {
	cc, err := remote.ClientConn()
	if err != nil {
		return err
	}

	_, err = cc.ReplicateRPC(context.Background(), &pb.ReplicateMsg{
		Key:        key,
		Value:      value,
		Deps:       deps,
		Shardstamp: shardstamp,
	})
	if err != nil {
		return err
	}

	return remote.connCheck(err)
}
