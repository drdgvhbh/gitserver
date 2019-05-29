package reference

type Reference struct {
	// The hash of the commit this reference points to
	//
	// required: true
	// example: e38e2cde1fada4a738f2461b283e561bc767568b
	Hash string `json:"hash,omitempty"`

	// The name of the reference
	//
	// required: true
	// example: refs/heads/master
	Name string `json:"name,omitempty"`
}
