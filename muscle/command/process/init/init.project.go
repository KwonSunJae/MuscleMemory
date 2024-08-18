package init

import (
	"fmt"
	process_error "muscle/command/error"
	systemCMD "muscle/command/system"
	"muscle/logger"
	"os"
)

type InitProject struct {
	Config map[string]string
}

func (i *InitProject) CheckArgValidate() error {
	// CheckArgValidate
	// -- Arguments List
	var essentialArgList = []string{
		// Essential Arguments
		"name",
		"dir",
	}

	// Check Essential Arguments
	for _, essentialArg := range essentialArgList {
		if _, ok := i.Config[essentialArg]; !ok {

			return fmt.Errorf("essential argument '%s' is missing", essentialArg)
		}
	}

	return nil
}

func (i *InitProject) InputConfig() error {
	// InputConfig
	pr := logger.NewPrinter()

	if _, ok := i.Config["input"]; ok { // if input is nil, load from 'muscle.ini' file
		name := pr.Ask("Please enter your project name")
		i.Config["name"] = name
	} else {

		if _, ok := i.Config["name"]; !ok {
			if _, ok := i.Config["n"]; ok {
				i.Config["name"] = i.Config["n"]
			} else {
				return fmt.Errorf("essential argument 'name' or 'n' flag is missing")
			}

		}
	}

	pr.Start("Read muscle.init file")
	conf, err := LoadConfig("muscle.init")
	if err != nil {
		pr.Error("muscle.init file is not exist. Please init")
		return process_error.NewError("muscle.init file is not exist", err)
	}
	i.Config["dir"] = conf["dir"]
	pr.Done()

	return nil
}

func (i *InitProject) Run() error {
	// Run
	pr := logger.NewPrinter()
	cmd := systemCMD.NewCommandSystemExecutor()
	// 1.Read muscle.init file
	dir := i.Config["dir"]

	//3. create branch new project at blank branch
	pr.Start("Create branch project from blank branch")

	err := cmd.Execute("git", "-C", dir, "checkout", "--orphan", i.Config["name"])
	if err != nil {
		pr.Error("You should check your branch name. It is already exist")
		return process_error.NewError("Create New branch failed.", err)
	}

	pr.Done()

	//4. create a dir and commit then push
	pr.Start("Create a dir and create project.conf file and commit push")
	err = os.Mkdir(dir+"/"+i.Config["name"], 0755)
	if err != nil {
		pr.Error("You should check your project name. It is already exist")
		return process_error.NewError("Create New project failed.", err)
	}

	err = cmd.Execute("touch", dir+"/"+i.Config["name"]+"/project.conf")
	if err != nil {
		pr.Error("project.conf file is not created")
		return process_error.NewError("Create New project failed.", err)
	}

	err = cmd.Execute("git", "-C", dir, "add", i.Config["name"]+"/project.conf")
	if err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}

	err = cmd.Execute("git", "-C", dir, "commit", "-m", "Create project.conf file")
	if err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}

	err = cmd.Execute("git", "-C", dir, "push", "origin", i.Config["name"])
	if err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}
	pr.Done()

	return nil
}
