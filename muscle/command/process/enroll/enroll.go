package enroll

import (
	process_error "muscle/command/error"
	"muscle/command/git"
	"muscle/logger"
	"muscle/util/crypt"
	"muscle/util/loader"
	"muscle/util/terraform"
	"os"
)

type Enroll interface {
	// Run the enroll process
	CheckLock() error
	Enroll() error
	UnLock() error
}

type EnrollImpl struct {
	// contains filtered or unexported fields
	Config map[string]string
}

func NewEnrollProcessor(config map[string]string) (Enroll, error) {
	return &EnrollImpl{Config: config}, nil
}

func (e *EnrollImpl) CheckLock() error {
	// CheckLock
	git := git.NewGit(e.Config["project-name"])
	pr := logger.NewPrinter()

	// git fetch and pull
	pr.Start("Fetch and Pull")
	if err := git.Fetch(); err != nil {
		pr.Error("Fetch Error")
		return process_error.NewError("Fetch Error", err)
	}

	if err := git.Pull(); err != nil {
		pr.Error("Pull Error")
		return process_error.NewError("Pull Error", err)
	}
	pr.Done()

	// check .lock file
	pr.Start("Check .lock file exist")
	conf, err := loader.LoadConfig(e.Config["project-name"] + "/" + e.Config["project-name"] + "/.lock")
	if err != nil {
		pr.Error("Load .lock file error")
		return process_error.NewError("Load .lock file error", err)
	}
	if crypt.CompareOwner(e.Config["owner"], conf["owner"]) {
		pr.Error("This project is locked by other :" + conf["owner"])
		return process_error.NewError("This project is locked by "+conf["owner"], nil)
	}
	pr.Done()

	return nil
}

func (e *EnrollImpl) Enroll() error {
	// Enroll
	pr := logger.NewPrinter()
	git := git.NewGit(e.Config["project-name"])

	// add all and commit
	pr.Start("Add all and Commit")
	if err := git.AddAll(); err != nil {
		pr.Error("Add all error")
		return process_error.NewError("Add all error", err)
	}

	if err := git.Commit("Enroll by muscle process"); err != nil {
		pr.Error("Commit error")
		return process_error.NewError("Commit error", err)
	}

	if err := git.PushBranch(e.Config["project-name"]); err != nil {
		pr.Error("Push error")
		return process_error.NewError("Push error", err)
	}

	pr.Done()

	return nil
}

func (e *EnrollImpl) UnLock() error {
	// UnLock
	pr := logger.NewPrinter()
	git := git.NewGit(e.Config["project-name"])

	// load project.conf file
	pr.Start("Load project.conf file")
	conf, err := loader.LoadConfig(e.Config["project-name"] + "/" + e.Config["project-name"] + "/project.conf")
	if err != nil {
		pr.Error("Load project.conf file error")
		return process_error.NewError("Load project.conf file error", err)
	}
	pr.Done()

	pr.Start("Check Deploy type")
	if conf["deploy"] != "local" {
		pr.Done()
		return nil
	}
	pr.Done()

	if conf["project-type"] == "terraform" {
		pr.Start("Terraform work")
		tf := terraform.NewTerraform(e.Config["project-name"] + "/" + e.Config["project-name"])
		if err := tf.TerraformInit(); err != nil {
			pr.Error("Terraform Init error")
			return process_error.NewError("Terraform Init error", err)
		}
		if err := tf.TerraformPlan(); err != nil {
			pr.Error("Terraform Plan error")
			return process_error.NewError("Terraform Plan error", err)
		}

		if err := tf.TerraformApply(); err != nil {
			pr.Error("Terraform Apply error")
			return process_error.NewError("Terraform Apply error", err)
		}
		pr.Done()
	} else {
		pr.Error("Not Supported yet..")
		return process_error.NewError("Not Supported yet..", nil)
	}

	// remove .lock file
	pr.Start("Remove .lock file")
	if err := os.Remove(e.Config["project-name"] + "/.lock"); err != nil {
		pr.Error("Remove .lock file error")
		return process_error.NewError("Remove .lock file error", err)
	}
	pr.Done()

	// git add all and commit, push
	pr.Start("Add all and Commit")
	if err := git.AddAll(); err != nil {
		pr.Error("Add all error")
		return process_error.NewError("Add all error", err)
	}

	if err := git.Commit("Release Lock by muscle process"); err != nil {
		pr.Error("Commit error")
		return process_error.NewError("Commit error", err)
	}

	if err := git.PushBranch(e.Config["project-name"]); err != nil {
		pr.Error("Push error")
		return process_error.NewError("Push error", err)
	}
	pr.Done()

	return nil
}
