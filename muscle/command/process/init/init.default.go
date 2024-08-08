package init

import "fmt"

type InitDefault struct {
	Config map[string]string
}

func (i *InitDefault) CheckConfig() error {
	// CheckConfig
	return nil
}

// -- Arguments List
var EssentialArgList = []string{}

var OptionalArgList = []string{
	// Optional Arguments
	"git-ssh",
}

func (i *InitDefault) CheckArgValidate() error {
	// CheckArgValidate

	// Check Essential Arguments
	for _, essentialArg := range EssentialArgList {
		if _, ok := i.Config[essentialArg]; !ok {
			return fmt.Errorf("essential argument '%s' is missing", essentialArg)
		}
	}

	return nil
}

func (i *InitDefault) InputConfig() error {
	// InputConfig

	if _, ok := i.Config["input"]; ok { // if input is nil, load from 'muscle.conf' file

	}

	return nil
}

func (i *InitDefault) Run() error {

	// Run
	return nil
}
