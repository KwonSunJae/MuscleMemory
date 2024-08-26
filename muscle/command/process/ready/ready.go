package ready

import (
	process_error "muscle/command/error"
	"muscle/command/git"
	systemCMD "muscle/command/system"
	"muscle/logger"
	"muscle/util/crypt"
	"muscle/util/loader"
	"muscle/util/writer"
	"os"
	"strconv"
	"time"
)

type Ready interface {
	LoadConfig() error
	Lock() error
	LoadRepository() error
	ReadyRepository() error
}

type ReadyImpl struct {
	// contains filtered or unexported fields
	Config map[string]string
}

func NewReadyProcessor(config map[string]string) (Ready, error) {
	return &ReadyImpl{Config: config}, nil
}

func (r *ReadyImpl) LoadConfig() error {
	// LoadConfig
	pr := logger.NewPrinter()
	pr.Start("Load Config")
	conf, err := loader.LoadConfig("muscle.init")
	if err != nil {
		pr.Error("muscle.init file is not exist. Please init")
		return process_error.NewError("muscle.init file is not exist", err)
	}
	r.Config["repository-git-url"] = conf["repository-git-url"]
	pr.Done()

	return nil
}
func (r *ReadyImpl) Lock() error {
	// Lock
	pr := logger.NewPrinter()
	cmd := systemCMD.NewCommandSystemExecutor()
	git := git.NewGit(r.Config["project-name"])

	projectDir := r.Config["project-name"] + "/" + r.Config["project-name"]

	// Check .lock file exist
	pr.Start("Check .lock file exist")
	if _, err := os.Stat(".lock"); err == nil {
		lockConf, err := loader.LoadConfig(projectDir + "/.lock")
		if err == nil {
			// ./lock file is not exist
			if _, ok := lockConf["expire"]; ok {

				expTimestamp, err := strconv.ParseInt(lockConf["expire"], 10, 64)
				if err != nil {
					return process_error.NewError("ParseInt", err)
				}

				currentTimestamp := time.Now().Unix()
				if currentTimestamp < expTimestamp && !crypt.CompareOwner(lockConf["owner"], r.Config["owner"]) {
					pr.Error("The project is locked by other user")
					return process_error.NewError("The project is locked", nil)
				}
			}

			// if expired, remove .lock file
			if err := cmd.Execute("rm", ".lock"); err != nil {
				return process_error.NewError("Remove .lock file", err)
			}
		}

	}
	pr.Done()

	// Create .lock file
	pr.Start("Create Lock file repository")
	if err := cmd.Execute("touch", ".lock"); err != nil {
		return process_error.NewError("Create Lock file repository", err)
	}
	pr.Done()

	// Write .lock file
	pr.Start("Write Lock file repository")
	currentTimestamp := time.Now().Unix()
	var duration int64 = 60 * 60 * 24
	var err error
	if _, ok := r.Config["lock-duration"]; ok {
		duration, err = strconv.ParseInt(r.Config["lock-duration"], 10, 64)
		if err != nil {
			return process_error.NewError("ParseInt", err)
		}
	}

	lockConf := map[string]string{
		"expire": strconv.Itoa(int(currentTimestamp + duration)),
		"owner":  r.Config["owner"],
	}
	if err := writer.WriteConfig(projectDir+"/.lock", lockConf); err != nil {
		return process_error.NewError("Write Lock file repository", err)
	}
	pr.Done()

	//Commit and Push
	pr.Start("Commit and Push")
	if err := git.AddAll(); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Commit and Push failed.", err)
	}
	if err := git.Commit("Lock project"); err != nil {
		pr.Error("Check your git status")
		return process_error.NewError("Commit and Push failed.", err)
	}
	if err := git.PushBranch(r.Config["project-name"]); err != nil {
		pr.Error("Maybe Other User Update Project. Try Again")
		if errs := git.ResetHard(); errs != nil {
			pr.Error("Please check your git status")
			return process_error.NewError("Commit and Push failed.", errs)
		}
		return process_error.NewError("Commit and Push failed.", err)
	}
	pr.Done()

	return nil
}

func (r *ReadyImpl) LoadRepository() error {
	// LoadRepository
	pr := logger.NewPrinter()
	git := git.NewGit(r.Config["project-name"])

	//Check Repository Exist
	pr.Start("Check Repository Exist")
	if _, err := os.Stat(r.Config["project-name"]); err != nil {
		// if not exist, git clone
		if err := git.CloneBranch(r.Config["repository-git-url"], r.Config["project-name"]); err != nil {
			pr.Error("git clone error")
			return process_error.NewError("git clone error", err)
		}
	} else {
		// if exist, git fetch
		if errs := git.Fetch(); errs != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", errs)
		}
		if errs := git.Pull(); errs != nil {
			pr.Error("Please init the project")
			return process_error.NewError("Please init the project", errs)
		}
	}
	pr.Done()

	return nil
}

func (r *ReadyImpl) ReadyRepository() error {
	// ReadyRepository
	pr := logger.NewPrinter()

	//Create Symbolic Link with project
	pr.Start("Create Symbolic Link")
	pwd, _ := os.Getwd()

	if err := os.Symlink(pwd+"/"+r.Config["project-name"]+"/"+r.Config["project-name"], r.Config["work-dir"]+"/"+r.Config["project-name"]); err != nil {
		pr.Error("Please check your work-dir, maybe dir is already exist.")
		return process_error.NewError("Create Symbolic Link Error", err)
	}

	pr.Done()

	return nil
}
