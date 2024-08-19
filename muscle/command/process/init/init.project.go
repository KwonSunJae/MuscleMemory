package init

import (
	"fmt"
	process_error "muscle/command/error"
	git "muscle/command/git"
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

			return process_error.NewError(fmt.Sprintf("essential argument '%s' is missing", essentialArg), nil)
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
				return process_error.NewError("Please enter your project name with 'n' or 'name' flag", nil)
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
	i.Config["repository-git-url"] = conf["repository-git-url"]
	pr.Done()

	return nil
}

func (i *InitProject) Run() error {
	// Run
	pr := logger.NewPrinter()
	cmd := systemCMD.NewCommandSystemExecutor()
	git := git.NewGit(i.Config["dir"])

	// 1.Read muscle.init file
	dir := i.Config["dir"]

	// 2. check project is already exist
	pr.Start("Check Main Branch") // IF CHECK OUT, IT may CAUSE Lock Error
	if _, err := os.Stat(dir); err == nil {
		if errs := git.Fetch(); errs != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", errs)
		}
		git.Pull()
	} else {
		if err := git.Clone(i.Config["repository-git-url"]); err != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", err)
		}
	}
	pr.Done()

	//3. create branch new project at blank branch
	pr.Start("Create branch project from blank branch")
	if err := git.Checkout("blank"); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New branch failed.", err)
	}
	if err := git.NewBlankBranch(i.Config["name"]); err != nil {
		pr.Error("You should check your branch name. Maybe name is duplicated")
		return process_error.NewError("Create New branch failed.", err)
	}
	pr.Done()

	//4. create a dir and commit then push
	pr.Start("Create a dir and create project.conf file and commit push")

	if err := os.Mkdir(dir+"/"+i.Config["name"], 0755); err != nil {
		pr.Error("You should check your project name. It is already exist")
		return process_error.NewError("Create New project failed.", err)
	}
	if err := cmd.Execute("touch", dir+"/"+i.Config["name"]+"/project.conf"); err != nil {
		pr.Error("project.conf file is not created")
		return process_error.NewError("Create New project failed.", err)
	}
	if err := git.AddAll(); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}
	if err := git.Commit("Create project.conf file"); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}
	if err := git.PushBranch(i.Config["name"]); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Create New project failed.", err)
	}
	pr.Done()

	// 5. checkout main branch
	pr.Start("Checkout Main Branch")
	if err := git.Checkout("main"); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Checkout main branch failed.", err)
	}
	pr.Done()

	return nil
}
