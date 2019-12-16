package main

import (
	"log"
	"os"
)

// Error is logger to Stderr
var Error *log.Logger

// Initialize the loggers
func init() {
	Error = log.New(os.Stderr, "ERROR: ", log.Ltime|log.Lshortfile)
}
