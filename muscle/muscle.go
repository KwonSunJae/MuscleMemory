package main

import (
	"muscle/command"
	"muscle/logger"
	"os"
)

func main() {
	// Get system arguments
	args := os.Args[1:]

	// Initialize logger
	log := logger.GetInstance()
	defer log.Sync()

	// Dispatch command
	commandDispatcher := command.NewCommandDispatcher()
	commandDispatcher.CommandDispatch(args)

}
