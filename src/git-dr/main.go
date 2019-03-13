package main

import (
	"fmt"
	"git-dr/bitbucket"
	"git-dr/cmd"
	"git-dr/git"
	"git-dr/hg"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var processName string

func init() {
	split := strings.Split(os.Args[0], "/")
	processName = split[len(split)-1]

	initConfig()
	initLog()
	initDirectories()
}

func main() {
	log.Println("[BEGIN]")

	// create & execute api request
	pages, err := bitbucket.GetRepos()
	if err != nil {
		log.Fatalf("unable to get repos: %s", err)
	}

	repoCount := 0
	// process api response
	for hasNext := true; hasNext; hasNext = pages.Next() {
		p := pages.Current()

		repositories := p["values"].([]interface{})
		for _, v := range repositories {
			name, scmType, cloneLink := getRepoInfo(v)
			log.Printf("[REPO] %s\n", name)

			// if the repo hasn't been cloned, we need to clone it.
			// if the repo has already been cloned, we just need to update it.
			if _, err := os.Stat(name); os.IsNotExist(err) {
				log.Printf("[CLONE] %s\n", cloneLink)
				switch scmType {
				case "hg":
					hg.Clone(cloneLink)
				case "git":
					git.Clone(cloneLink)
				}
			} else {
				log.Printf("[UPDATE] %s\n", name)

				cmd.Chdir(name)
				switch scmType {
				case "hg":
					hg.Update(name)
				case "git":
					git.Update(name)
				}
				cmd.Chdir("../")
			}

			repoCount++
		}
	}
	log.Printf("[REPORT] Handled %d repos\n", repoCount)
	log.Println("[END]")
}

// getRepoInfo casts an interface and pulls out the repo name, scm type (git or hg) & a link to clone it
func getRepoInfo(v interface{}) (name, scmType, cloneLink string) {
	// go's type system is an unending nightmare
	repo := v.(map[string]interface{})
	name = repo["slug"].(string)
	scmType = repo["scm"].(string)

	links := repo["links"].(map[string]interface{})
	cloneLinks := links["clone"].([]interface{})
	for _, v := range cloneLinks {
		link := v.(map[string]interface{})
		if link["name"].(string) == "https" {
			username := viper.GetString("bitbucket.username")
			password := viper.GetString("bitbucket.app_password")
			combined := fmt.Sprintf("%s:%s", username, password)
			cloneLink = strings.Replace(link["href"].(string), username, combined, -1)
		}
	}
	// nevermind, the nightmare is over. we're good now

	return
}
