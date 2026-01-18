package main

import (
	"log"
	"os"
	"runtime/debug"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC RECOVERED: %v\n", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	StartServer()
}
