package generator

import (
	"html/template"
	format "muscle/format/gitactions"
	"muscle/util/checker"
	"os"
)

type GeneratorGitActions struct {
	Config map[string]string
}

func (g *GeneratorGitActions) CheckConfig() error {
	// CheckConfig
	if err := checker.CheckArgValidate(g.Config, []string{
		"project-type",
		"public",
	}); err != nil {
		return err
	}

	return nil
}

func (g *GeneratorGitActions) Generate(filepath string) error {

	// Generate
	// Generate GitActions yaml file
	var script string

	if g.Config["project-type"] == "terraform" {
		if g.Config["public"] == "true" {
			script = format.GitActionsTerraformPublicTemplate
		} else {
			if g.Config["pem"] == "true" {
				script = format.GitAcitonsTerraformTemplateWithKey
			} else {
				script = format.GitAcitonsTerraformTemplateWithPW
			}
		}
	} else if g.Config["project-type"] == "ansible" {
		if g.Config["public"] == "true" {
			script = format.GitActionsAnsiblePublicTemplate
		} else {
			if g.Config["pem"] == "true" {
				script = format.GitActionsAnsibleTemplateWithKey
			} else {
				script = format.GitActionsAnsibleTemplateWithPW
			}
		}
	}

	g.Config["github_token"] = "${{ secrets.GITHUB_TOKEN }}"
	g.Config["commit_message"] = "${{ github.event.head_commit.message }}"
	g.Config["result"] = "${{ steps.check_commit_message.outputs.deploy }}"
	g.Config["result_output"] = "${{ needs.check_commit_message.outputs.deploy }}"
	g.Config["branch"] = g.Config["project-name"]
	tmpl, err := template.New("yaml").Parse(script)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	result := tmpl.Execute(file, g.Config)
	defer file.Close()

	return result

}
