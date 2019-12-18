package main

import "fmt"

func (node *Node) runFollower() stateFunction {
	fmt.Println("Running follower")
	for {
		select {
		case f := <-node.applyChan:
			err := node.fsm.Apply(f.log)
			if err != nil {
				f.errChan <- err
			} else {
				f.resChan <- 0
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
