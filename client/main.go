package main

import (
	"errors"
	"log"

	"github.com/abiosoft/ishell"
)

func main() {
	client, err := CreateClient()
	if err != nil {
		log.Fatal(err)
	}

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "insert",
		Help: "insert <key> <value>",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Err(errors.New("Args should be length 2"))
				return
			}
			key := c.Args[0]
			value := []byte(c.Args[1])
			err := client.Write(key, value)
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("Success")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: "update <key> <value>",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Err(errors.New("Args should be length 2"))
				return
			}
			key := c.Args[0]
			value := []byte(c.Args[1])
			err := client.Update(key, value)
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("Success")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete <key>",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Err(errors.New("Args should be length 1"))
				return
			}
			key := c.Args[0]
			err := client.Delete(key)
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("Success")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "read",
		Help: "read <key>",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Err(errors.New("Args should be length 1"))
				return
			}
			key := c.Args[0]
			value, err := client.Read(key)
			if err != nil {
				c.Err(err)
				return
			}
			c.Printf("Key: %v, Value: %v\n", key, string(value))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "timestamp",
		Help: "client causal timestamp",
		Func: func(c *ishell.Context) {
			c.Println(client.Timestamp)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "exit shell",
		Func: func(c *ishell.Context) {
			shell.Stop()
		},
	})

	shell.Run()
}
