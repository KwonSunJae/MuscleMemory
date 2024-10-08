package init

import (
	process_error "muscle/command/error"
	"muscle/command/git"
	"muscle/logger"
	"muscle/util/checker"
	"muscle/util/loader"
	"muscle/util/writer"
	"os"
)

type InitTerraform struct {
	Config map[string]string
}

func (i *InitTerraform) CheckArgValidate() error {
	// CheckArgValidate

	if err := checker.CheckArgValidate(i.Config, []string{
		"project-name",
		"dir",
		"repository-git-url",
	}); err != nil {
		return process_error.NewError("essential argument is missing", err)
	}

	return nil
}

func (i *InitTerraform) InputConfig() error {
	pr := logger.NewPrinter()
	// InputConfig
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
				return process_error.NewError("Please enter your project name with 'n' or 'p' or'project-name' or 'project' flag", nil)
			}
		}
	}
	pr.Start("Read muscle.init file")
	conf, err := loader.LoadConfig("muscle.init")
	if err != nil {
		pr.Error("muscle.init file is not exist. Please init")
		return process_error.NewError("muscle.init file is not exist", err)
	}
	i.Config["dir"] = conf["dir"]
	i.Config["repository-git-url"] = conf["repository-git-url"]
	pr.Done()

	return nil
}

func (i *InitTerraform) Run() error {
	// Run
	pr := logger.NewPrinter()
	git := git.NewGit(i.Config["project-name"])
	// 1. CLone Branch project
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

	// readfile dir + /project.conf file
	pr.Start("Check project.conf file")
	conf, err := loader.LoadConfig(i.Config["project-name"] + "/" + i.Config["project-name"] + "/project.conf")
	if err != nil {
		pr.Error("Please init the project")
		return process_error.NewError("Please init the project", err)
	}
	if _, ok := conf["project-type"]; ok && conf["project-type"] != "terraform" {
		pr.Error("This project is aleady set other type.")
		return process_error.NewError("This project is aleady set other type.", nil)
	}
	pr.Done()

	// 3. Config project.conf file
	pr.Start("Config project.conf file")
	conf["project-type"] = "terraform"
	conf["project-name"] = i.Config["project-name"]
	conf["deploy"] = "local"
	if err := writer.WriteConfig(i.Config["project-name"]+"/"+i.Config["project-name"]+"/project.conf", conf); err != nil {
		pr.Error("Config project.conf file error")
		return process_error.NewError("Config project.conf file error", err)
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
		pr.Error("Check your git status. Maybe no changes in your project")
		return process_error.NewError("Commit and Push failed.", err)
	}

	pr.Done()

	return nil
}
