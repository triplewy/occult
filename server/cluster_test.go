package main

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/triplewy/occult/proto"
)

func TestSimpleCluster(t *testing.T) {
	nodes, err := BootstrapSimpleCluster([]int{50000, 50001})
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	leader := nodes[0]
	follower := nodes[1]
	if leader.State == Follower {
		leader = nodes[1]
		follower = nodes[0]
	}

	writeReply, err := leader.InsertRPC(context.Background(), &pb.WriteMsg{
		Key:   "test",
		Value: []byte("value"),
		Deps:  0,
	})

	ts := writeReply.Shardstamp

	readReply, err := follower.ReadRPC(context.Background(), &pb.KeyMsg{Key: "test"})
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	fmt.Println(ts)
	fmt.Println(readReply)
}
