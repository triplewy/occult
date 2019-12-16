package main

import (
	"errors"
)

func BootstrapSimpleCluster(ports []int) ([]*Node, error) {
	if len(ports) != 2 {
		return nil, errors.New("Expected 2 ports")
	}
	node1, err := CreateNode(&Config{port: ports[0]})
	if err != nil {
		return nil, err
	}
	node2, err := CreateNode(&Config{port: ports[1]})
	if err != nil {
		return nil, err
	}

	return []*Node{node1, node2}, nil
}
