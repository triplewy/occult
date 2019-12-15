package main

import "fmt"

func (node *Node) runFollower() stateFunction {
	for {
		select {
		case event := <-node.watch:
			fmt.Println(event)
		}

	}
}
