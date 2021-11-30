package github

import (
	"github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
	"sort"
)

//Print all the prs and the count
//Debug method
func PrintPullRequests(pulls []*github.PullRequest) {
	//print all pull requests
	for _, p := range pulls {
		log.Info("Created: ", *p.CreatedAt, "\nRepo name: ", *p.Head.Repo.Name, "\nPR name: ", *p.Title, "\nPR #URL", *p.URL, "\n")
	}
	log.Info("\nNumber of prs: %d\n", len(pulls))
}

func RemovePullRequestsWithLabel(pulls []*github.PullRequest, label *github.Label) []*github.PullRequest {
	var y []*github.PullRequest
	for i := 0; i < len(pulls); i++ {
		if ContainsLabelWithName(pulls[i].Labels, label) == false {
			y = append(y, pulls[i])
		}
	}
	return y
}

func GetPullRequestsWithLabel(pulls []*github.PullRequest, label *github.Label) []*github.PullRequest {
	var prs []*github.PullRequest
	for i := 0; i < len(pulls); i++ {
		if ContainsLabelWithName(pulls[i].Labels, label) == true {
			prs = append(prs, pulls[i])
		}
	}
	return prs
}

func RemovePRByIndex(s []*github.PullRequest, index int) []*github.PullRequest {
	return append(s[:index], s[index+1:]...)
}

// sortPullRequestByTime sorts by oldest pr createdAt time
func sortPullRequestByTime(p []*github.PullRequest) []*github.PullRequest {
	sort.Slice(p, func(i, j int) bool {
		return p[i].CreatedAt.Before(*p[j].CreatedAt)
	})
	return p
}
