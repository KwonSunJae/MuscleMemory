package init

import (
	process_error "muscle/command/error"
	git "muscle/command/git"
	"muscle/logger"
	"os"
)

type InitAnsible struct {
	Config map[string]string
}

func (i *InitAnsible) CheckArgValidate() error {
	// CheckArgValidate
	if err := CheckArgValidate(i.Config, []string{
		"repository-git-url",
	}); err != nil {
		return process_error.NewError("essential argument is missing", err)
	}

	return nil
}

func (i *InitAnsible) InputConfig() error {
	// InputConfig
	pr := logger.NewPrinter()
	if _, ok := i.Config["input"]; ok {
		projectName := pr.Ask("Please enter your project name")
		i.Config["project-name"] = projectName
	} else {
		if _, ok := i.Config["project-name"]; !ok {
			if _, ok := i.Config["project"]; ok {
				i.Config["project-name"] = i.Config["project"]
			}
			if _, ok := i.Config["p"]; ok {
				i.Config["project-name"] = i.Config["p"]
			}
			if _, ok := i.Config["n"]; ok {
				i.Config["project-name"] = i.Config["n"]
			} else {
				return process_error.NewError("Please enter your project name with 'n' or 'p' or 'project-name' flag", nil)
			}
		}
	}

	return nil
}

func (i *InitAnsible) Run() error {
	// Run

	pr := logger.NewPrinter()
	git := git.NewGit(i.Config["dir"])

	// 1. checkout project
	pr.Start("Clone Branch project") // IF CHECK OUT, IT may CAUSE Lock Error

	//check project is already exist
	if _, err := os.Stat(i.Config["project-name"]); err == nil {
		if errs := git.Fetch(); errs != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", err)
		}
		if errs := git.Pull(); errs != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", err)
		}
	} else {
		if err := git.CloneBranch(i.Config["repository-git-url"], i.Config["project-name"]); err != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", err)
		}
	}

	pr.Done()

	// 2. check project.conf file
	pr.Start("Check project.conf file")
	conf, err := LoadConfig(i.Config["project-name"] + "/" + i.Config["project-name"] + "/project.conf")
	if err != nil {
		pr.Error("Please init the project")
		return process_error.NewError("Please init the project", err)
	}
	if conf["project-type"] != "ansible" {
		pr.Error("Please init the project")
		return process_error.NewError("Please init the project", nil)
	}
	pr.Done()

	// 3. wirte project.conf

	pr.Start("Write project.conf file")
	conf["project-type"] = "ansible"
	conf["project-name"] = i.Config["project-name"]
	if err := WriteConfig(i.Config["dir"]+"/"+i.Config["name"]+"/project.conf", conf); err != nil {
		pr.Error("Write project.conf file error")
		return process_error.NewError("Write project.conf file error", err)
	}
	pr.Done()

	// 4. Commit and Push
	pr.Start("Commit and Push")
	if err := git.AddAll(); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Commit and Push failed.", err)
	}

	if err := git.Commit("Config project.conf file"); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Commit and Push failed.", err)
	}

	if err := git.PushBranch(i.Config["project-name"]); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Commit and Push failed.", err)
	}

	pr.Done()

	return nil
}
