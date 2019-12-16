package main

import "fmt"

func (node *Node) runLeader() stateFunction {
	fmt.Println("Running leader")
	for {
		select {
		case entry := <-node.applyChan:
			node.shardstamp++
			entry.Deps = node.shardstamp
			// Update local fsm
			err := node.fsm.Apply(entry)
			if err != nil {
				entry.errChan <- err
				break
			}
			// Replicate to followers
			for _, remote := range node.nodes {
				go func(remote *RemoteNode) {
					err := remote.ReplicateRPC(entry.Key, entry.Value, entry.Deps, node.shardstamp)
					if err != nil {
						Error.Println(err)
					}
				}(remote)
			}
			// Return local shardstamp to client
			entry.resChan <- node.shardstamp
		case <-node.watch:
			err := node.checkZk()
			if err != nil {
				fmt.Println(err)
			}
			if node.State == Follower {
				return node.runFollower
			}
		}
	}
}
