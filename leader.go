package main

import "fmt"

func (node *Node) runLeader() stateFunction {
	for {
		select {
		case event := <-node.watch:
			fmt.Println(event)
		}

	}
}
