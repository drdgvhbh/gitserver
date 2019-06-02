package repository

// swagger:parameters listCommits listReferences
type Params struct {
	// The directory of the repository
	//
	// in: path
	// required: true
	Directory string `json:"directory"`
}
