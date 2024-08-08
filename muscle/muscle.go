package main

import (
	"muscle/command"
	"os"
)

func main() {
	// Get system arguments
	args := os.Args[1:]

	commandDispatcher := command.NewCommandDispatcher()
	commandDispatcher.CommandDispatch(args)

}
