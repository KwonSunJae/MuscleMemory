package git

type Github interface {
	Clone() error
	Commit() error
	Fetch() error
	Push() error
	Merge() error
	Rebase() error
	CreateBranch() error
	DeleteBranch() error
	Checkout() error
}

type Git struct {
	// contains filtered or unexported fields
}
