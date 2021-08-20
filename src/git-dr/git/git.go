package git

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"git-dr/cmd"
)

// Clone clones the repo at the url provided
func Clone(url string) {
	git("clone %s --bare", url)
}

func RemoteBranches() []string {
	//branchesCmd := fmt.Sprintf("-c \"%s\"", "for b in `git branch -r | grep -v -- '->'`; do git branch --track ${b##origin/} $b; done")
	branches := cmd.Run("git", "branch -r")

	var branchList []string
	scanner := bufio.NewScanner(strings.NewReader(branches))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "->") && strings.Contains(line, "origin") {
			branchList = append(branchList, strings.TrimSpace(line))
		}
	}

	return branchList
}

// Update cd's into the repository and updates the submodules and pulls any new changes
func Update(name string) {
	// checkout every branch to make sure we're tracking everything
	branches := RemoteBranches()
	for _, b := range branches {
		noRemote := strings.Replace(b, "origin/", "", 1)
		git("checkout %s", noRemote)
	}

	// some repos don't have a master
	for _, b := range branches {
		if b == "master" {
			git("checkout master")
		}
	}
	//git("submodule update --recursive --remote")
	//git("pull --all")
	git("fetch --all")
}

// Push pushes every branch to the remote
func Push() {
	git("push --all github")
}

// AddRemote adds a remote to the repo
func AddRemote(name string, username string, password string) {
	git("remote add github https://%s:%s@github.com/AccuWeather-Inc/%s.git", username, password, name)
}

// git executes the specified git command
func git(s string, args ...interface{}) {
	cmdString := s
	if len(args) > 0 {
		cmdString = fmt.Sprintf(s, args...)
	}
	log.Printf("[GIT] %s\n", cmdString)
	cmd.Run("git", cmdString)
}
