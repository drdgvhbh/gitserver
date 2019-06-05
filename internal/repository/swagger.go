package repository

// swagger:parameters listCommits listReferences getCommit
type Params struct {
	// The directory of the repository
	//
	// in: path
	// required: true
	Directory string `json:"directory"`
}

// swagger:parameters getCommit
type CommitIDParams struct {
	// The hash of the commit
	//
	// in: path
	// required: true
	Hash string `json:"hash"`
}
