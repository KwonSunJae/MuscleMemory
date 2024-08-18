package git

import (
	systemCMD "muscle/command/system"
)

type Git interface {
	Checkout(branch string) error
	Clone(url string) error
	CloneBranch(url string, branch string) error
	AddAll() error
	Commit(message string) error
	NewBlankBranch(branch string) error
	Push() error
	PushBranch(branch string) error
	Fetch() error
	Pull() error
	Rebase() error
	Reset() error
	Branch() error
}

type GitImpl struct {
	// contains filtered or unexported fields
	repository string
	cmd        systemCMD.CommandSystem
}

func NewGit(repo string) Git {
	return &GitImpl{repository: repo, cmd: systemCMD.NewCommandSystemExecutor()}
}

func (g *GitImpl) CloneBranch(url string, branch string) error {
	// Clone
	g.repository = branch
	return g.cmd.Execute("git", "clone", "-b", branch, url, branch)
}

func (g *GitImpl) PushBranch(branch string) error {
	// Push
	return g.cmd.Execute("git", "-C", g.repository, "push", "origin", branch)
}
func (g *GitImpl) NewBlankBranch(branch string) error {
	// Checkout
	return g.cmd.Execute("git", "-C", g.repository, "checkout", "--orphan", branch)
}
func (g *GitImpl) Checkout(branch string) error {
	// Checkout
	return g.cmd.Execute("git", "-C", g.repository, "checkout", branch)
}

func (g *GitImpl) Clone(url string) error {
	// Clone

	return g.cmd.Execute("git", "clone", url)
}

func (g *GitImpl) AddAll() error {
	// Add
	return g.cmd.Execute("git", "-C", g.repository, "add", ".")
}

func (g *GitImpl) Commit(message string) error {
	// Commit
	return g.cmd.Execute("git", "-C", g.repository, "commit", "-m", message)
}

func (g *GitImpl) Push() error {
	// Push
	return g.cmd.Execute("git", "-C", g.repository, "push", "origin", "main")
}

func (g *GitImpl) Fetch() error {
	// Fetch
	return g.cmd.Execute("git", "-C", g.repository, "fetch")
}

func (g *GitImpl) Pull() error {

	// Pull
	return g.cmd.Execute("git", "-C", g.repository, "pull")
}

func (g *GitImpl) Rebase() error {

	// Rebase
	return g.cmd.Execute("git", "-C", g.repository, "rebase")
}

func (g *GitImpl) Reset() error {

	// Reset
	return g.cmd.Execute("git", "-C", g.repository, "reset")
}

func (g *GitImpl) Branch() error {

	// Branch
	return g.cmd.Execute("git", "-C", g.repository, "branch")
}
