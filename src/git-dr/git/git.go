package git

import (
	"fmt"
	"log"

	"bitbucket.org/accuweather/git-dr/src/git-dr/cmd"
)

// Clone clones the repo at the url provided
func Clone(url string) {
	git("clone %s", url)
}

// Update cd's into the repository and updates the submodules and pulls any new changes
func Update(name string) {
	git("submodule update --recursive --remote")
	git("pull --all")
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
