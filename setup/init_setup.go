package main

/*
   only meant to run once per project setup, when first initializing the application from the boilerplate app code
   this script will:
   - checkout a new git branch based on the project name
   - customize the application files with your project name
   - change go.mod module to your repo location
   - add your repo as a new git remote to push the code to
   - output next steps about what to do with the setup
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		os.Exit(1)
	}
	Repo = repo

	gitAddRemote := exec.Command(
		"git",
		"remote",
		"add",
		ProjectNameWithDashes,
		Repo,
	)
	err := gitAddRemote.Run()
	if err != nil {
		fmt.Printf("[git remote add %s %s] error: %v\nCheck if the remote already exists for this project. Remove if necessary.\n", ProjectNameWithDashes, Repo, err)
		os.Exit(1)
	}

	// checkout new git branch with ProjectNameWithDashes
	gitCheckoutBranch := exec.Command(
		"git",
		"checkout",
		"-b",
		ProjectNameWithDashes,
	)
	err = gitCheckoutBranch.Run()
	if err != nil {
		fmt.Printf("[git checkout] error: %v\nIf the branch already exists for this project decide to remove the branch or pick a different project name.\n", err)
		os.Exit(1)
	}

	// make all needed file changes
	err = updateFilesForProject()
	if err != nil {
		fmt.Printf("[updateFileForProject] error: %v\n", err)
		os.Exit(1)
	}

	gitAddAll := exec.Command(
		"git",
		"add",
		"--all",
	)
	err = gitAddAll.Run()
	if err != nil {
		fmt.Printf("[git add --all] error: %v\n", err)
		os.Exit(1)
	}

	gitCommit := exec.Command(
		"git",
		"commit",
		"-m",
		fmt.Sprintf("Initialization script: customize files for project %s\n", ProjectName),
	)
	err = gitCommit.Run()
	if err != nil {
		fmt.Printf("[git commit] error: %v\n", err)
		os.Exit(1)
	}

	// return to the default branch: serverless
	gitCheckoutServerless := exec.Command(
		"git",
		"checkout",
		"serverless",
	)
	err = gitCheckoutServerless.Run()
	if err != nil {
		fmt.Printf("[git checkout serverless] error: %v\n", err)
		os.Exit(1)
	}
	// print message to user on the changes made and next steps to push the branch to their new repo location
	fmt.Printf(`
	
	Project Setup Completed!

	Branch Created with above changes: %s
	
	To see all the changes that were made by this script, run the following git command:
		git diff serverless..%s
	
	TL;DR - files updated:
		- go.mod
		- package.json
		- README.md
		- serverless.yml
	
	Remote added: %s
	
	When you are ready, push these changes to your remote github repository and start hacking away at your new github app!
		git push %s %s

	`, ProjectNameWithDashes, ProjectNameWithDashes, Repo, ProjectNameWithDashes, ProjectNameWithDashes)

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

func updateFilesForProject() error {
	// go.mod : module line should be the Repo url
	var repoOwner, repoName string
	switch {
	case strings.Contains(Repo, "https://"):
		theSplit := strings.Split(Repo, "/")
		repoOwner = theSplit[3]
		repoName = strings.Split(theSplit[4], ".")[0]
	case strings.Contains(Repo, "git@"):
		theSplit := strings.Split(Repo, "/")
		repoOwner = strings.Split(theSplit[0], ":")[1]
		repoName = strings.Split(theSplit[1], ".")[0]
	default:
		return fmt.Errorf("[Repo] url is not in an expected format: %s", Repo)
	}
	err := sed(
		"module github.com/sharkysharks/go-github-app-boilerplate",
		fmt.Sprintf("module github.com/%s/%s", repoOwner, repoName),
		"go.mod",
	)
	if err != nil {
		return fmt.Errorf("[sed go.mod] err: %v", err)
	}

	// package.json : name should be ProjectName, version should be 0.1.0
	err = updateJsonFile("package.json")
	if err != nil {
		return fmt.Errorf("[sed package.json] err: %v", err)
	}

	// README.md : first heading should be ProjectName
	err = sed(
		"# go-github-app-boilerplate",
		fmt.Sprintf("# %s", ProjectName),
		"README.md",
	)
	if err != nil {
		return fmt.Errorf("[sed README.md] err: %v", err)
	}

	// serverless.yml : custom.projectName should be ProjectName
	err = sed(
		"go-github-app-boilerplate",
		ProjectNameWithDashes,
		"serverless.yml",
	)
	if err != nil {
		return fmt.Errorf("[sed serverless.yml] err: %v", err)
	}

	return nil
}

func sed(old, new, filePath string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileString := string(fileData)
	fileString = strings.ReplaceAll(fileString, old, new)
	fileData = []byte(fileString)

	err = ioutil.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func updateJsonFile(filePath string) error {
	byteVal, err := ioutil.ReadFile("package.json")
	if err != nil {
		return err
	}
	var data struct {
		Name            string                 `json:"name"`
		Version         string                 `json:"version"`
		Description     string                 `json:"description"`
		Scripts         map[string]interface{} `json:"scripts"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}
	err = json.Unmarshal(byteVal, &data)
	if err != nil {
		return err
	}

	data.Name = ProjectName
	data.Version = "0.1.0"

	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, []byte(d), 0644)
	if err != nil {
		return err
	}
	return nil
}
