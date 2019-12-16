package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/abiosoft/ishell"
)

var port int

func init() {
	flag.IntVar(&port, "c", 50000, "connect to node")
}

func printHelp(shell *ishell.Shell) {
	shell.Println("Commands:")
	shell.Println(" - help                    Prints this help message")
	shell.Println(" - read <key>		      Read")
	shell.Println(" - write <key> <value>     Write")
	shell.Println(" - delete <key>            Delete")
	shell.Println(" - exit                    Exit CLI")
}

func main() {
	flag.Parse()

	client, err := CreateClient(&Config{Addr: fmt.Sprintf("localhost:%d", port)})
	if err != nil {
		log.Fatal(err)
	}

	shell := ishell.New()
	printHelp(shell)

	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "print help",
		Func: func(c *ishell.Context) {
			printHelp(shell)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "read",
		Help: "read key",
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
			c.Printf("Key: %v, Value: %v", key, string(value))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "exit shell",
		Func: func(c *ishell.Context) {
			shell.Stop()
		},
	})

	shell.Start()
}
