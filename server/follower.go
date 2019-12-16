package main

import "fmt"

func (node *Node) runFollower() stateFunction {
	fmt.Println("Running follower")
	for {
		select {
		case entry := <-node.applyChan:
			err := node.fsm.Apply(entry)
			if err != nil {
				entry.errChan <- err
			} else {
				entry.resChan <- 0
			}
		case <-node.watch:
			err := node.checkZk()
			if err != nil {
				fmt.Println(err)
			}
			if node.State == Leader {
				return node.runLeader
			}
		}
	}
}
