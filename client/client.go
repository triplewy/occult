package main

import "math/rand"

type RemoteNode struct {
	Addr string
}

type Client struct {
	nodes     []*RemoteNode
	leader    *RemoteNode
	Timestamp uint64
}

func CreateClient() (*Client, error) {
	client := &Client{
		nodes:     []*RemoteNode{},
		leader:    nil,
		Timestamp: 0,
	}

	err := client.retrieveNodes()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Write(key string, value []byte) error {
	reply, err := c.leader.InsertRPC(key, value, c.Timestamp)
	if err != nil {
		return err
	}
	c.Timestamp = max(c.Timestamp, reply.Shardstamp)
	return nil
}

func (c *Client) Update(key string, value []byte) error {
	reply, err := c.leader.UpdateRPC(key, value, c.Timestamp)
	if err != nil {
		return err
	}
	c.Timestamp = max(c.Timestamp, reply.Shardstamp)
	return nil
}

func (c *Client) Delete(key string) error {
	reply, err := c.leader.DeleteRPC(key)
	if err != nil {
		return err
	}
	c.Timestamp = max(c.Timestamp, reply.Shardstamp)
	return nil
}

func (c *Client) Read(key string) ([]byte, error) {
	remote := c.getServer()

	entry, err := remote.ReadRPC(key)
	if err != nil {
		return nil, err
	}
	if entry.Shardstamp < c.Timestamp {
		return c.finishStaleRead(key)
	}
	c.Timestamp = max(c.Timestamp, entry.Deps)
	return entry.Value, nil
}

func (c *Client) getServer() *RemoteNode {
	// Chooses a random server
	i := rand.Intn(len(c.nodes))
	return c.nodes[i]
}

func (c *Client) finishStaleRead(key string) ([]byte, error) {
	return nil, nil
}
