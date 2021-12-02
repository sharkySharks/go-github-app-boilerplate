package github

import (
	"context"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v37/github"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

type Client struct {
	Ctx       context.Context
	GitHub    *github.Client
	RepoOwner string
}

type Repo struct {
	Name     string
	Branches []string
}

// InitClient initializes a github client to respond back to github after events are received
func InitClient(installationId int64, githubAppID int64, githubPrivateKey []byte) (*github.Client, error) {
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(
		tr,
		githubAppID,
		int64(installationId),
		githubPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	c := github.NewClient(&http.Client{Transport: itr})

	return c, nil
}

// CreateLabelInRepo created a label in a provided repo
func CreateLabelInRepo(ghc Client, label *github.Label, repo string) error {
	_, _, err := ghc.GitHub.Issues.CreateLabel(ghc.Ctx, ghc.RepoOwner, repo, label)
	if err != nil {
		// retry creating the label before throwing an error
		if _, _, retryError := ghc.GitHub.Issues.CreateLabel(ghc.Ctx, ghc.RepoOwner, repo, label); retryError != nil {
			log.Error("Error creating new label: ", err)
			return err
		}
	}
	return nil
}

// CreateLabelsInRepos creates a list of labels in a set of repos
func CreateLabelsInRepos(ghc Client, repos []Repo, labels []*github.Label) error {
	for _, repo := range repos {
		for _, label := range labels {
			_, _, err := ghc.GitHub.Issues.GetLabel(ghc.Ctx, ghc.RepoOwner, repo.Name, *label.Name)
			//create label if it does not exist
			if err != nil {
				log.Printf("Label %v does not exist, creating it in repo %v", label, repo)
				err := CreateLabelInRepo(ghc, label, repo.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// AddLabelToPR adds a label to a pull request
func AddLabelToPR(ghc Client, pull *github.PullRequest, label *github.Label) error {
	if _, _, err := ghc.GitHub.Issues.AddLabelsToIssue(
		ghc.Ctx,
		ghc.RepoOwner,
		*pull.Head.Repo.Name,
		int(*pull.Number),
		[]string{*label.Name},
	); err != nil {
		log.Error("Error adding label to pr, retrying... ", err)
		if _, _, retryError := ghc.GitHub.Issues.AddLabelsToIssue(
			ghc.Ctx,
			ghc.RepoOwner,
			*pull.Head.Repo.Name,
			int(*pull.Number),
			[]string{*label.Name},
		); retryError != nil {
			log.Error("Error retrying adding label: ", retryError)
			return err
		}
	}
	return nil
}

// RemoveLabelsFromPR removes a set of labels from a pull request
func RemoveLabelsFromPR(ghc Client, pull *github.PullRequest, labels []*github.Label) error {
	if pull.Labels != nil {
		//for each one see if it has labels
		for j := 0; j < len(pull.Labels); j++ {
			//if the label equals the labels in pr, then remove it
			if ContainsLabelWithName(labels, pull.Labels[j]) {
				if res, err := ghc.GitHub.Issues.RemoveLabelForIssue(
					ghc.Ctx,
					ghc.RepoOwner,
					*pull.Head.Repo.Name,
					int(*pull.Number),
					*pull.Labels[j].Name,
				); err != nil || (res.StatusCode != 200 && res.StatusCode != 404) {
					log.Error("Error removing label from pr, retrying... ", err)
					if retry, retryError := ghc.GitHub.Issues.RemoveLabelForIssue(
						ghc.Ctx,
						ghc.RepoOwner,
						*pull.Head.Repo.Name,
						int(*pull.Number),
						*pull.Labels[j].Name,
					); retryError != nil || (retry.StatusCode != 200 && retry.StatusCode != 404) {
						log.Error("Error retrying removing label: ", retryError)
						return retryError
					}
				}
			}
		}
	}
	return nil
}

// RemoveOldLabelsFromPR removes labels that match a particular regex pattern from a pull request
func RemoveOldLabelsFromPR(ghc Client, pull *github.PullRequest, regexPattern string) error {
	var s string
	var b []byte

	if pull.Labels != nil {
		//for each one see if it has labels
		for _, label := range pull.Labels {
			re := regexp.MustCompile(regexPattern)
			s = *label.Name
			b = []byte(s)
			if string(re.Find(b)) != "" {
				if res, err := ghc.GitHub.Issues.RemoveLabelForIssue(
					ghc.Ctx,
					ghc.RepoOwner,
					*pull.Head.Repo.Name,
					int(*pull.Number),
					*label.Name,
				); err != nil || (res.StatusCode != 200 && res.StatusCode != 404) {
					log.Error("Error removing label from pr, retrying... ", err)
					if retry, retryError := ghc.GitHub.Issues.RemoveLabelForIssue(
						ghc.Ctx,
						ghc.RepoOwner,
						*pull.Head.Repo.Name,
						int(*pull.Number),
						*label.Name,
					); retryError != nil || (retry.StatusCode != 200 && retry.StatusCode != 404) {
						log.Error("Error retrying removing label: ", retryError)
						return retryError
					}
				}
			}
		}
	}
	return nil
}

// GetAllPullRequests returns a list of all pull requests from repos/branches
// default behavior is to get open pull requests
func GetAllPullRequests(ghc Client, repos []Repo) ([]*github.PullRequest, error) {
	var pulls []*github.PullRequest
	for _, repo := range repos {
		for _, branch := range repo.Branches {
			opt := github.PullRequestListOptions{
				"open", "", branch, "created", "asc", github.ListOptions{},
			}
			p, _, err := ghc.GitHub.PullRequests.List(ghc.Ctx, ghc.RepoOwner, repo.Name, &opt)
			if err != nil {
				log.Error("Error getting PR list: ", err)
				return nil, err
			}
			for _, pull := range p {
				pulls = append(pulls, pull)
			}
		}
	}
	return sortPullRequestByTime(pulls), nil
}

// GetOnePullRequest returns a pull request by pr number for a provided repo
func GetOnePullRequest(ghc Client, repo string, pr_number int) (*github.PullRequest, error) {
	p, _, err := ghc.GitHub.PullRequests.Get(ghc.Ctx, ghc.RepoOwner, repo, pr_number)
	if err != nil {
		log.Error("Error getting one PR: ", err)
		return nil, err
	}

	return p, nil
}

// ListBranchCommits returns a slice of commits for a given repo/branch
func ListBranchCommits(ghc Client, repo string, branch string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{
		SHA: branch,
		// Path:   "",
		// Author: "",
		// Since:  time.Date(2013, time.August, 1, 0, 0, 0, 0, time.UTC),
		// Until:  time.Date(2013, time.September, 3, 0, 0, 0, 0, time.UTC),
	}

	commitInfo, _, err := ghc.GitHub.Repositories.ListCommits(ghc.Ctx, ghc.RepoOwner, repo, opt)
	if err != nil {
		log.Error("Repositories.ListCommits returned error: %v", err)
		return nil, err
	}

	return commitInfo, nil

}

// CompareCommits returns the number of commits between two commits
func CompareCommits(ghc Client, repo string, base string, head string) (int, error) {
	compareInfo, _, err := ghc.GitHub.Repositories.CompareCommits(ghc.Ctx, ghc.RepoOwner, repo, base, head, &github.ListOptions{})
	if err != nil {
		log.Error("Repositories.CompareCommits returned error: %v", err)
		return -1, err
	}

	return *compareInfo.BehindBy, nil
}

// MergeBranch updates a given branch with a given 'head' commit SHA
func MergeBranch(ghc Client, repo string, branch string, head string) (*github.RepositoryCommit, error) {
	var commitMessage = "The branch was brought up to date with head."
	input := &github.RepositoryMergeRequest{
		Base:          &branch,
		Head:          &head,
		CommitMessage: &commitMessage,
	}

	commit, _, err := ghc.GitHub.Repositories.Merge(ghc.Ctx, ghc.RepoOwner, repo, input)
	if err != nil {
		log.Error("Repositories.Merge returned error: %v", err)
		return nil, err
	}

	return commit, nil
}

// ListStatuses returns a list of pull requests statuses that are associated with a particular ref
func ListStatuses(ghc Client, repo string, ref string) ([]*github.RepoStatus, error) {
	listStatuses, _, err := ghc.GitHub.Repositories.ListStatuses(ghc.Ctx, ghc.RepoOwner, repo, ref, nil)
	if err != nil {
		log.Error("Repositories.ListStatuses returned error: %v", err)
		return nil, err
	}
	return listStatuses, nil
}

// ListCollaborators returns github users that are listed as collaborators on a repository
func ListCollaborators(ghc Client, repo string) ([]*github.User, error) {
	collaborators, _, err := ghc.GitHub.Repositories.ListCollaborators(ghc.Ctx, ghc.RepoOwner, repo, nil)
	if err != nil {
		log.Error("Repositories.ListCollaborators returned error: %v", err)
		return nil, err
	}
	return collaborators, nil
}

// CreateComment adds a comment to a github issue/pr
func CreateComment(ghc Client, repo string, number int, commentString string) (*github.IssueComment, error) {
	input := &github.IssueComment{
		Body: &commentString,
	}

	retComment, _, err := ghc.GitHub.Issues.CreateComment(ghc.Ctx, ghc.RepoOwner, repo, number, input)
	if err != nil {
		log.Error("Issues.CreateComment returned error: %v", err)
		return nil, err
	}
	return retComment, nil
}
