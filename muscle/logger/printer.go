package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Printer struct
type Printer struct {
	taskName string
	done     chan struct{}
}

// Ask function to ask a question and return the answer
func (p *Printer) Ask(question string) string {
	log := GetInstance()
	log.Info("Asking question", zap.String("question", question))

	var answer string
	fmt.Print(question + ": ")
	fmt.Scanln(&answer)

	log.Info("Received answer", zap.String("answer", answer))
	return answer
}

// Start function to start a task with a loading indicator
func (p *Printer) Start(task string) {
	log := GetInstance()
	log.Info("Starting task", zap.String("task", task))

	p.taskName = task
	p.done = make(chan struct{})

	go p.loadingIndicator()

	log.Info("Task started", zap.String("task", task))
}

// loadingIndicator displays a loading animation
func (p *Printer) loadingIndicator() {
	for {
		select {
		case <-p.done:
			return
		default:
			for _, pattern := range []string{"...", ".. ", ".  ", "   "} {
				fmt.Printf("\r%s %s", p.taskName, pattern)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

// Warn function to log a warning message
func (p *Printer) Warn(message string) {
	log := GetInstance()
	log.Warn("Warning", zap.String("task_name", p.taskName), zap.String("message", message))
	fmt.Printf("Warning: %s\n", message)
}

// Error function to log an error message and stop the loading indicator
func (p *Printer) Error(action string) {
	log := GetInstance()
	log.Error("Error occurred", zap.String("task_name", p.taskName), zap.String("action", action))

	close(p.done) // Stop the loading indicator
	fmt.Printf("\rError: %s\n", action)
}

// Done function to indicate task completion
func (p *Printer) Done() {
	log := GetInstance()
	log.Info("Task completed", zap.String("task_id", p.taskName))

	close(p.done) // Stop the loading indicator
	fmt.Printf("\rTask %s completed successfully.\n", p.taskName)
}

func NewPrinter() *Printer {
	return &Printer{}
}
