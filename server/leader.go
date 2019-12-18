package main

import "fmt"

func (node *Node) runLeader() stateFunction {
	fmt.Println("Running leader")
	for {
		select {
		case f := <-node.applyChan:
			node.shardstamp++
			f.log.Deps = node.shardstamp
			// Update local fsm
			err := node.fsm.Apply(f.log)
			if err != nil {
				f.errChan <- err
				break
			}
			// Replicate to followers
			for _, remote := range node.nodes {
				if remote.Addr != node.addr {
					go func(remote *RemoteNode) {
						err := remote.ReplicateRPC(f.log.Key, f.log.Value, f.log.Deps, node.shardstamp)
						if err != nil {
							Error.Println(err)
						}
					}(remote)
				}
			}
			// Return local shardstamp to client
			f.resChan <- node.shardstamp
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
