package main

import (
	"context"

	pb "github.com/triplewy/occult/proto"
)

func (node *Node) ConnectRPC(ctx context.Context, msg *pb.EmptyMsg) (*pb.CausalTsMsg, error) {
	node.shardstampLock.RLock()
	shardstamp := node.shardstamp
	node.shardstampLock.RUnlock()

	return &pb.CausalTsMsg{Shardstamp: shardstamp}, nil
}

func (node *Node) WriteRPC(ctx context.Context, msg *pb.WriteMsg) (*pb.ShardTsMsg, error) {
	if node.State != Leader {
		return nil, ErrNotLeader
	}
	entry := NewEntry(msg.Key, msg.Value, msg.Deps)
	shardstamp, err := node.apply(entry)
	if err != nil {
		return nil, err
	}
	return &pb.ShardTsMsg{Shardstamp: shardstamp}, nil
}

func (node *Node) ReadRPC(ctx context.Context, msg *pb.ReadMsg) (*pb.EntryMsg, error) {
	node.shardstampLock.RLock()
	shardstamp := node.shardstamp
	node.shardstampLock.RUnlock()

	value, deps, err := node.fsm.Read(msg.Key)
	if err != nil {
		return nil, err
	}
	return &pb.EntryMsg{
		Value:      value,
		Deps:       deps,
		Shardstamp: shardstamp,
	}, nil
}

func (node *Node) ReplicateRPC(ctx context.Context, msg *pb.ReplicateMsg) (*pb.EmptyMsg, error) {
	node.shardstampLock.Lock()
	node.shardstamp = msg.Shardstamp
	node.shardstampLock.Unlock()

	entry := NewEntry(msg.Key, msg.Value, msg.Deps)
	_, err := node.apply(entry)
	return &pb.EmptyMsg{}, err
}
