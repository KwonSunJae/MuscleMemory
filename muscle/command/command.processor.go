package command

type Processor interface {
	Init() error
	Process() error
}

type CommandProcessor struct {
	// contains filtered or unexported fields
}

func NewCommandProcessor() *CommandProcessor {
	return &CommandProcessor{}
}

func (c *CommandProcessor) Init() error {
	// Init
	return nil
}

func (c *CommandProcessor) Process() error {
	// Process
	return nil
}
