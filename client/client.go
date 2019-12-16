package main

type RemoteNode struct {
	Addr string
}

type Client struct {
	Config    *Config
	nodes     []*RemoteNode
	Timestamp uint64
}

func CreateClient(config *Config) (*Client, error) {
	client := &Client{
		Config:    config,
		nodes:     []*RemoteNode{},
		Timestamp: 0,
	}

	return client, nil
}

func (c *Client) getServer() (*RemoteNode, error) {
	return nil, nil
}

func (c *Client) Read(key string) ([]byte, error) {
	remote, err := c.getServer()
	if err != nil {
		return nil, err
	}
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

func (c *Client) finishStaleRead(key string) ([]byte, error) {
	return nil, nil
}
