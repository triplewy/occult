package main

import (
	"context"

	pb "github.com/triplewy/occult/proto"
)

func (node *Node) InsertRPC(ctx context.Context, msg *pb.WriteMsg) (*pb.ShardTsMsg, error) {
	if node.State != Leader {
		return nil, ErrNotLeader
	}
	_, _, err := node.fsm.Read(msg.Key)
	if err == nil {
		return nil, ErrKeyExists
	}
	if err != ErrKeyNotFound {
		return nil, err
	}
	log := &Log{
		Command: pb.Command_Insert,
		Key:     msg.Key,
		Value:   msg.Value,
		Deps:    msg.Deps,
	}
	shardstamp, err := node.Apply(log)
	if err != nil {
		return nil, err
	}
	return &pb.ShardTsMsg{Shardstamp: shardstamp}, nil
}

func (node *Node) UpdateRPC(ctx context.Context, msg *pb.WriteMsg) (*pb.ShardTsMsg, error) {
	if node.State != Leader {
		return nil, ErrNotLeader
	}
	_, _, err := node.fsm.Read(msg.Key)
	if err != nil {
		return nil, err
	}
	log := &Log{
		Command: pb.Command_Update,
		Key:     msg.Key,
		Value:   msg.Value,
		Deps:    msg.Deps,
	}
	shardstamp, err := node.Apply(log)
	if err != nil {
		return nil, err
	}
	return &pb.ShardTsMsg{Shardstamp: shardstamp}, nil
}

func (node *Node) DeleteRPC(ctx context.Context, msg *pb.KeyMsg) (*pb.ShardTsMsg, error) {
	if node.State != Leader {
		return nil, ErrNotLeader
	}
	log := &Log{
		Command: pb.Command_Delete,
		Key:     msg.Key,
		Value:   nil,
		Deps:    0,
	}
	shardstamp, err := node.Apply(log)
	if err != nil {
		return nil, err
	}
	return &pb.ShardTsMsg{Shardstamp: shardstamp}, nil
}

func (node *Node) ReadRPC(ctx context.Context, msg *pb.KeyMsg) (*pb.EntryMsg, error) {
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

	log := &Log{
		Command: msg.Command,
		Key:     msg.Key,
		Value:   msg.Value,
		Deps:    msg.Deps,
	}
	_, err := node.Apply(log)
	return &pb.EmptyMsg{}, err
}
