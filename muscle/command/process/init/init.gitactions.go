package init

import (
	process_error "muscle/command/error"
	git "muscle/command/git"
	"muscle/generator"
	"muscle/logger"
	"muscle/util/checker"
	"muscle/util/loader"
	"muscle/util/writer"
	"os"
)

type InitGitActions struct {
	Config map[string]string
}

func (i *InitGitActions) CheckArgValidate() error {
	// CheckArgValidate

	if err := checker.CheckArgValidate(i.Config, []string{
		"project-name",
		"dir",
		"repository-git-url",
		"useremail",
		"username",
	}); err != nil {
		return process_error.NewError("essential argument is missing", err)
	}

	return nil
}

func (i *InitGitActions) InputConfig() error {
	// InputConfig
	pr := logger.NewPrinter()
	// InputConfig
	if _, ok := i.Config["input"]; ok {
		projectName := pr.Ask("Please enter your project name")
		i.Config["project-name"] = projectName

		public := pr.Ask("Is your Provider is public? (y/n)")
		if public == "y" {
			i.Config["public"] = "true"
			useremail := pr.Ask("Please enter your github email")
			i.Config["useremail"] = useremail
			username := pr.Ask("Please enter your github username")
			i.Config["username"] = username

		} else {
			return process_error.NewError("Not Supported yet", nil)
			// pem := pr.Ask("If you access to server, Use pem file? (y/n)")
			// if pem == "y" {
			// 	i.Config["pem"] = "true"
			// }
			// pw := pr.Ask("If you access to server, Use password? (y/n)")
			// if pw == "y" {
			// 	i.Config["pw"] = "true"
			// }
		}

	} else {
		if _, ok := i.Config["project-name"]; !ok {
			if _, ok := i.Config["project"]; ok {
				i.Config["project-name"] = i.Config["project"]
			} else if _, ok := i.Config["p"]; ok {
				i.Config["project-name"] = i.Config["p"]
			} else if _, ok := i.Config["n"]; ok {
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
	i.Config["repository-git-url"] = conf["repository-git-url"]
	pr.Done()

	pr.Start("Read project.conf file")
	conf, err = loader.LoadConfig(i.Config["project-name"] + "/" + i.Config["project-name"] + "/project.conf")
	if err != nil {
		pr.Error("project.conf file is not exist. Please init")
		return process_error.NewError("project.conf file is not exist", err)
	}
	i.Config["project-type"] = conf["project-type"]
	i.Config["dir"] = conf["dir"]
	pr.Done()

	pr.Start("Overwrite project conf file")
	i.Config["deploy"] = "gitactions"
	if err := writer.WriteConfig(i.Config["project-name"]+"/"+i.Config["project-name"]+"/project.conf", i.Config); err != nil {
		pr.Error("overwrite project conf file error")
		return process_error.NewError("overwrite project conf file error", err)
	}
	pr.Done()

	return nil
}

func (i *InitGitActions) Run() error {
	// Run
	pr := logger.NewPrinter()
	git := git.NewGit(i.Config["project-name"])

	// check project-name dir  exist
	pr.Start(("Check project-name dir exist"))
	if _, err := os.Stat(i.Config["project-name"]); os.IsNotExist(err) {
		// then clone repository
		if err := git.CloneBranch(i.Config["repository-git-url"], i.Config["project-name"]); err != nil {
			pr.Error("git clone error")
			return process_error.NewError("git clone error", err)
		}
	}
	pr.Done()

	// check .github/workflows dir exist
	pr.Start("Check .github/workflows dir exist")
	if _, err := os.Stat(i.Config["project-name"] + "/.github/workflows"); os.IsNotExist(err) {
		// then create .github/workflows dir
		if err := os.MkdirAll(i.Config["project-name"]+"/.github/workflows", 0755); err != nil {
			pr.Error("create .github/workflows dir error")
			return process_error.NewError("create .github/workflows dir error", err)
		}
	}
	pr.Done()
	actionfile := i.Config["project-name"] + "/.github/workflows/" + i.Config["project-name"] + ".yml"
	// create github actions file
	pr.Start("Touch main.yml file")
	if _, err := os.Create(actionfile); err != nil {
		pr.Error("create main.yml file error")
		return process_error.NewError("create main.yml file error", err)
	}
	pr.Done()

	// generate github actions file
	pr.Start("Generate github actions file")
	gen, err := generator.NewGenerator("gitactions", i.Config)
	if err != nil {
		pr.Error("generator error")
		return process_error.NewError("generator error", err)
	}
	if err := gen.Generate(actionfile); err != nil {
		pr.Error("generate github actions file error")
		return process_error.NewError("generate github actions file error", err)
	}
	pr.Done()

	// git add, commit, push
	pr.Start("Git add, commit, push")
	if err := git.AddAll(); err != nil {
		pr.Error("git add error")
		return process_error.NewError("git add error", err)
	}
	if err := git.Commit("Add github actions file"); err != nil {
		pr.Error("git commit error")
		return process_error.NewError("git commit error", err)
	}
	if err := git.PushBranch(i.Config["project-name"]); err != nil {
		pr.Error("git push error")
		return process_error.NewError("git push error", err)
	}
	pr.Done()

	return nil
}
