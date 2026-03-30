package types

type PRStatus string

const (
	StatusOpen   PRStatus = "open"
	StatusClosed PRStatus = "closed"
	StatusMerged PRStatus = "merged"
)

type PR struct {
	Number         int
	Title          string
	Org            string
	Repo           string
	Author         string
	Status         PRStatus
	Labels         []string
	URL            string
	ReviewDecision string
}

func (pr PR) RepoPath() string {
	return pr.Org + "/" + pr.Repo
}
