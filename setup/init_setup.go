package main

/*
   only meant to run once per project setup, when first initializing the application from the boilerplate app code
   this script will:
   - checkout a new git branch based on the project name
   - customize the application files with your project name
   - change go.mod module to your repo location
   - add your repo as a new git remote to push the code to
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	ProjectName           string
	ProjectNameWithDashes string
	Repo                  string
)

func main() {
	// get project name
	project := stringPrompt("What is your project name?")
	ProjectName = project
	ProjectNameWithDashes = strings.Join(strings.Split(strings.ToLower(project), " "), "-")

	// get git repo url for git remote
	repo := stringPrompt("Provide the full url of the repo where this app will live. [Use the copy button in the github repo ui.]")
	//	verify that the repo entry matches the regex pattern for either https://github.com/user/repo.git or git@github.com:user/repo.git pattern - works with git.corp addresses as well
	isValidRepoDest, _ := regexp.MatchString(`^(https://|git@)[a-zA-Z.]+(:|/)[a-zA-Z-_]+/[a-zA-Z-_]+.git$`, repo)
	if !isValidRepoDest {
		fmt.Printf("The repo url is not in one of the following patterns: https://github.com/user/repo.git || git@github.com:user/repo.git\nyour entry: %s\n", repo)
		return
	}
	Repo = repo

	// checkout new git branch with ProjectNameWithDashes
	cmd := exec.Command(
		"git",
		"checkout",
		"-b",
		ProjectNameWithDashes,
	)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[git checkout] error: %v", err)
		return
	}
	// add git remote with repo
	cmd2 := exec.Command(
		"git",
		"remote",
		"add",
		ProjectNameWithDashes,
		Repo,
	)
	err = cmd2.Run()
	if err != nil {
		fmt.Printf("[git remote add %s %s] error: %v", ProjectNameWithDashes, Repo, err)
		return
	}

	// make all needed file changes
	err = updateFileForProject()
	if err != nil {
		fmt.Printf("[git remote add %s %s] error: %v", ProjectNameWithDashes, Repo, err)
		return
	}

	// git add + git commit
	cmd3 := exec.Command(
		"git",
		"add",
		"--all",
	)
	err = cmd3.Run()
	if err != nil {
		fmt.Printf("[git add --all] error: %v", err)
		return
	}
	cmd4 := exec.Command(
		"git",
		"commit",
		"-m",
		fmt.Sprintf("Initialization script: customize files for project %s", ProjectName),
	)
	err = cmd4.Run()
	if err != nil {
		fmt.Printf("[git commit] error: %v", err)
		return
	}
	// print message to user on the changes made and next steps to push the branch to their new repo location
	fmt.Printf(`
	
	Insert text here about files that were changed successfully and next steps to push to remote from new branch...

	`)

}

func stringPrompt(prompt string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		_, _ = fmt.Fprint(os.Stderr, prompt+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func updateFileForProject() error {
	// go.mod : module line should be the Repo url

	// package.json : name should be ProjectName

	// README.md : first heading should be ProjectName

	// serverless.yaml : custom.projectName should be ProjectName

	return nil
}
