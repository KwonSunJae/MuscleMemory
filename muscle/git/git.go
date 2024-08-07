package git

type Github interface {
	Init(repositoryURL string, userEmail string, userName string) error
	Clone() error
	Commit() error
	Fetch() error
	Push() error
	Merge() error
	Rebase() error
	CreateBranch() error
	DeleteBranch() error
	Checkout() error
	BranchList() ([]string, error)
}

type Git struct {
	// contains filtered or unexported fields
}
