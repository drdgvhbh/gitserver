package commit

type Contributor struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type LogData struct {
	Hash      string       `json:"hash,omitempty"`
	Summary   string       `json:"summary,omitempty"`
	Author    *Contributor `json:"author,omitempty"`
	Committer *Contributor `json:"committer,omitempty"`
}

