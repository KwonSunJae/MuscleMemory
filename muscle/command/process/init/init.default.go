package init

import (
	"fmt"
	"io/ioutil"
	systemCMD "muscle/command/system"
	"os"
	"os/exec"
	"strings"
)

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

	// init시 필요한항목

	// var initItems = [] string{
	// 	"ssh-keygen",
	// 	"user-email",
	// 	"user-name",
	// 	"repository-git-url",
	// }

	if _, ok := i.Config["input"]; ok { // if input is nil, load from 'muscle.ini' file
		fmt.Print("Do you want to enroll new ssh-key for git? (y/n): ")
		var temp string
		fmt.Scanln(&temp)

		if temp == "y" {
			i.Config["ssh-keygen"] = "true"
			var userEmail string
			fmt.Print("Enter your email: ")
			fmt.Scanln(&userEmail)
			i.Config["user-email"] = userEmail

		} else {
			fmt.Println("Skip ssh-keygen")
		}

		fmt.Print("Have you already set git in your host? (y/n): ")
		fmt.Scanln(&temp)

		if temp == "n" {
			i.Config["git-set"] = "true"
			var userEmail, userName, repositoryGitURL string
			fmt.Print("Enter your email: ")
			fmt.Scanln(&userEmail)
			fmt.Print("Enter your name: ")
			fmt.Scanln(&userName)
			fmt.Print("Enter your repository git url (ex: git@): ")
			fmt.Scanln(&repositoryGitURL)
			i.Config["user-email"] = userEmail
			i.Config["user-name"] = userName
			i.Config["repository-git-url"] = repositoryGitURL
		}

	} else {
		// Read 'muscle.ini' file
		file, err := os.Open("muscle.ini")
		if err != nil {
			return err
		}
		defer file.Close()

	}

	return nil
}

func (i *InitDefault) Run() error {

	// Run
	cmd := systemCMD.NewCommandSystemExecutor()

	fmt.Println("Start init process")

	if i.Config["ssh-keygen"] == "true" {
		fmt.Println("Start ssh-keygen process")
		// 1. SSH 키 생성
		err := cmd.Execute("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", i.Config["user-email"], "-f", os.Getenv("HOME")+"/.ssh/id_rsa", "-N", "")
		stdout, err := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", i.Config["user-email"], "-f", os.Getenv("HOME")+"/.ssh/muscle.pub", "-N", "").CombinedOutput()
		fmt.Println(string(stdout))
		if err != nil {
			return fmt.Errorf("ssh-keygen error: %v", err)
		}
		// 2. 생성한 SSH 키 출력
		pubKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa.pub"
		pubKey, err := ioutil.ReadFile(pubKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read SSH public key: %v", err)
		}
		fmt.Printf("Generated SSH Key:\n%s\n", string(pubKey))
		fmt.Println("Please add the above SSH key to your Github account: https://github.com/settings/keys")
	}

	if i.Config["git-set"] == "true" {
		fmt.Println("Start git-set process")
		// 3. Git 설정
		if err := cmd.Execute("git", "config", "--global", "user.email", i.Config["user-email"]); err != nil {
			return fmt.Errorf("git config --global user.email error: %v", err)
		}
		if err := cmd.Execute("git", "config", "--global", "user.name", i.Config["user-name"]); err != nil {
			return fmt.Errorf("git config --global user.name error: %v", err)
		}
	}

	// 4. Git Clone
	fmt.Println("Start git clone process")
	if err := cmd.Execute("git", "clone", i.Config["repository-git-url"]); err != nil {
		return fmt.Errorf("git clone error: %v", err)
	}

	// 5. Generate muscle.init file at clone repository
	fmt.Println("Start generate muscle.init file")
	dir := strings.Split(strings.Split(i.Config["repository-git-url"], "/")[1], ".")[0]
	if err := cmd.Execute("touch", dir+"/muscle.init"); err != nil {
		return fmt.Errorf("touch muscle.init error: %v", err)
	}

	// 6. Write config to muscle.init file
	fmt.Println("Start write config to muscle.init file")
	file, err := os.OpenFile(dir+"/muscle.init", os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open muscle.init file: %v", err)
	}
	for key, value := range i.Config {
		if _, err := file.WriteString(key + "=" + value + "\n"); err != nil {
			return fmt.Errorf("failed to write config to muscle.init file: %v", err)
		}
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close muscle.init file: %v", err)
	}

	// 7. Git Add, Commit, Push
	fmt.Println("Start git push muscle.init file")
	if err := cmd.Execute("git", "add", "."); err != nil {
		return fmt.Errorf("git add error: %v", err)
	}
	if err := cmd.Execute("git", "commit", "-m", "Add muscle.init file"); err != nil {
		return fmt.Errorf("git commit error: %v", err)
	}
	if err := cmd.Execute("git", "push", "origin", "master"); err != nil {
		return fmt.Errorf("git push error: %v", err)
	}

	// 8. Git Create Branch blank, then puplish
	fmt.Println("Start git create branch blank")
	if err := cmd.Execute("git", "checkout", "-b", "blank"); err != nil {
		return fmt.Errorf("git checkout -b blank error: %v", err)
	}
	if err := cmd.Execute("git", "push", "origin", "blank"); err != nil {
		return fmt.Errorf("git push origin blank error: %v", err)
	}

	return nil
}
