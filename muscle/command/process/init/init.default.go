package init

import (
	"fmt"
	process_error "muscle/command/error"
	git "muscle/command/git"
	systemCMD "muscle/command/system"
	"muscle/logger"
	"os"
	"strings"
)

type InitDefault struct {
	Config map[string]string
}

func (i *InitDefault) CheckArgValidate() error {
	// CheckArgValidate
	var essentialArgList = []string{
		// Essential Arguments
		"repository-git-url",
	}

	// Check Essential Arguments
	for _, essentialArg := range essentialArgList {
		if _, ok := i.Config[essentialArg]; !ok {

			return process_error.NewError(fmt.Sprintf("essential argument '%s' is missing", essentialArg), nil)
		}
	}

	return nil
}

func (i *InitDefault) InputConfig() error {
	// InputConfig

	// init시 필요한항목

	// var initItems = [] string{
	// 	"repository-git-url",
	// }

	pr := logger.NewPrinter()

	if _, ok := i.Config["input"]; ok { // if input is nil, load from 'muscle.ini' file

		repositoryGitURL := pr.Ask("Enter your repository git url")
		i.Config["repository-git-url"] = repositoryGitURL

	} else {
		// Read 'muscle.ini' file
		if _, ok := i.Config["repository-git-url"]; !ok {
			if _, ok := i.Config["r"]; ok {
				i.Config["repository-git-url"] = i.Config["r"]
			} else if _, ok := i.Config["repo"]; ok {
				i.Config["repository-git-url"] = i.Config["repo"]
			} else {
				return process_error.NewError("Please enter your repository git url with 'r' or 'repo' flag", nil)
			}
		}
	}

	return nil
}

func (i *InitDefault) Run() error {

	// Run
	pr := logger.NewPrinter()

	cmd := systemCMD.NewCommandSystemExecutor()
	dir := strings.Split(strings.Split(i.Config["repository-git-url"], "/")[1], ".")[0]
	git := git.NewGit(dir)
	// 4. Git Clone
	pr.Start("git clone process")

	i.Config["dir"] = dir
	if err := git.Clone(i.Config["repository-git-url"]); err != nil {
		// if Already Exist Repository, delete and clone
		pr.Error("check your repository is already exist")
		return process_error.NewError("git clone error: You should have to Delete Dir.%v", err)
	}
	pr.Done()

	// Check muscle.init file exist
	pr.Start("Check muscle.init file exist")
	if _, err := os.Stat(dir + "/muscle.init"); err == nil {
		ans := pr.Ask("muscle.init file is already exist. Do you want to delete it? (y/n)")
		if ans == "y" {

			if err := cmd.Execute("rm", dir+"/muscle.init"); err != nil {
				pr.Error("check muscle.init file is exist")
				return process_error.NewError("rm muscle.init error", err)
			}

		} else {
			pr.Done()
			return nil
		}
	}
	pr.Done()

	// 5. Generate muscle.init file at clone repository
	pr.Start("Start generate muscle.init file")
	if err := cmd.Execute("touch", dir+"/muscle.init"); err != nil {
		return process_error.NewError("touch muscle.init error: %v", err)
	}
	pr.Done()

	// 6. Write config to muscle.init file
	pr.Start("Start write config to muscle.init file")
	file, err := os.OpenFile(dir+"/muscle.init", os.O_WRONLY, 0644)
	if err != nil {
		pr.Error("check muscle.init file is exist")
		return process_error.NewError("failed to open muscle.init file", err)
	}
	for key, value := range i.Config {
		if _, err := file.WriteString(key + "=" + value + "\n"); err != nil {
			pr.Error("check muscle.init file is exist")
			return process_error.NewError("failed to write config to muscle.init file", err)
		}
	}
	if err := file.Close(); err != nil {
		pr.Error("check muscle.init file is exist")
		return process_error.NewError("failed to close muscle.init file", err)
	}
	pr.Done()

	// 7. Git Add, Commit, Push
	pr.Start("Start git add, commit, push process")
	if err := git.AddAll(); err != nil {
		pr.Error("check git add")
		return process_error.NewError("git add error", err)
	}
	if err := git.Commit("Create muscle.init file"); err != nil {
		pr.Error("check git commit")
		return process_error.NewError("git commit error", err)
	}
	if err := git.Push(); err != nil {
		pr.Error("check git push")
		return process_error.NewError("git push error", err)
	}
	pr.Done()

	// 8. Copy dir/muscle.init file to muscle.init file
	pr.Start("Copy muscle.init file to muscle.init")
	if err := cmd.Execute("cp", dir+"/muscle.init", "muscle.init"); err != nil {
		pr.Error("check muscle.init file is exist")
		return process_error.NewError("cp muscle.init error", err)
	}
	pr.Done()

	// 10.Create blank branch then commit and push
	pr.Start("Create blank branch")
	if err := git.NewBlankBranch("blank"); err != nil {
		pr.Error("check git branch")
		return process_error.NewError("git branch error", err)
	}
	if err := git.AddAll(); err != nil {
		pr.Error("check git add")
		return process_error.NewError("git add error", err)
	}
	if err := git.Commit("Create blank branch"); err != nil {
		pr.Error("check git commit")
		return process_error.NewError("git commit error", err)
	}
	if err := git.PushBranch("blank"); err != nil {
		pr.Error("check git push")
		return process_error.NewError("git push error", err)
	}

	pr.Done()

	return nil
}
