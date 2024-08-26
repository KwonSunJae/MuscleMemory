package terraform

import (
	systemCMD "muscle/command/system"
)

type Terraform interface {
	TerraformInit() error
	TerraformPlan() error
	TerraformApply() error
}

type TerraformImpl struct {
	cmd     systemCMD.CommandSystem
	workdir string
}

func NewTerraform(workdir string) Terraform {
	return &TerraformImpl{cmd: systemCMD.NewCommandSystemExecutor(), workdir: workdir}
}

func (t *TerraformImpl) TerraformPlan() error {
	// Terraform Plan
	if err := t.cmd.Execute("terraform", "plan", "-chdir="+t.workdir); err != nil {
		return err
	}
	return nil
}

func (t *TerraformImpl) TerraformApply() error {
	// Terraform Apply
	if err := t.cmd.Execute("terraform", "apply", "-auto-approve", "-chdir="+t.workdir); err != nil {
		return err
	}
	return nil
}

func (t *TerraformImpl) TerraformInit() error {
	// Terraform Init
	if err := t.cmd.Execute("terraform", "init", "-chdir="+t.workdir); err != nil {
		return err
	}
	return nil
}
