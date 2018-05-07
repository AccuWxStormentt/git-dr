package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/accuweather/git-dr/src/git-dr/api"
	"bitbucket.org/accuweather/git-dr/src/git-dr/cmd"
	"bitbucket.org/accuweather/git-dr/src/git-dr/git"
	"bitbucket.org/accuweather/git-dr/src/git-dr/hg"
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
	// create api client
	auth := newAuth()
	client := api.Client{auth}

	// create & execute api request
	req := client.NewRequest("https://api.bitbucket.org/2.0/repositories/%s", "accuweather")
	pages, err := req.Do()
	if err != nil {
		log.Fatalf("unable to make request: %s", err)
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

// newAuth creates a new api.Authenticator. Right now it only creates a BasicAuth object. Later, if it needs to, it can support OAuth2 or whatever else Bitbucket wants to do
func newAuth() api.Authenticator {
	username := viper.GetString("USERNAME")
	app_password := viper.GetString("APP_PASSWORD")
	return api.NewBasicAuth(username, app_password)
}

// getRepoInfo casts an interface and pulls out the repo name, scm type (git or hg) & a link to clone it
func getRepoInfo(v interface{}) (name, scmType, cloneLink string) {
	// go's type system is an unending nightmare
	repo := v.(map[string]interface{})
	name = repo["name"].(string)
	scmType = repo["scm"].(string)

	links := repo["links"].(map[string]interface{})
	cloneLinks := links["clone"].([]interface{})
	for _, v := range cloneLinks {
		link := v.(map[string]interface{})
		if link["name"].(string) == "https" {
			username := viper.GetString("USERNAME")
			password := viper.GetString("APP_PASSWORD")
			combined := fmt.Sprintf("%s:%s", username, password)
			cloneLink = strings.Replace(link["href"].(string), username, combined, -1)
		}
	}
	// nevermind, the nightmare is over. we're good now

	return
}
