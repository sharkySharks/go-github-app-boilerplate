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
	"strings"
)

var (
	ProjectName string
	Repo        string
)

func ProjectNamePrompt(project string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		_, _ = fmt.Fprint(os.Stderr, project+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func RepoPrompt(repo string) string {
	return repo
}

func main() {
	project := ProjectNamePrompt("What is your project name?")
	fmt.Printf("this is your project: %s\n", project)

	//repo := RepoPrompt("Provide the full url of the repo where this app will live. [Use the copy button in the github ui.]")
}
