package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	raft "github.com/triplewy/raft/src"
)

var port int

func init() {
	flag.IntVar(&port, "p", 30000, "rpc port for node")
}

func main() {
	flag.Parse()

	raft.GetOutboundIP()

	config := &Config{port: port}
	_, err := CreateNode(config)
	if err != nil {
		panic(err)
	}

	log.Println("Occult started successfully")
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("Occult exiting")
}
