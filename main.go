package main

import (
	"log"
	"os"
	"strings"

	"bitbucket.org/accuweather/git-dr/api"
	"bitbucket.org/accuweather/git-dr/cmd"
	"bitbucket.org/accuweather/git-dr/git"
	"bitbucket.org/accuweather/git-dr/hg"
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

// initConfig loads the viper configuration from git-dr.[json|yml|toml]
func initConfig() {
	viper.SetConfigName("git-dr")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("unable to read config file: %s", err)
	}

	// true: log to stdout, false: log to LOG_DIR
	viper.SetDefault("LOG_STD", true)
	// directory to save log files in
	viper.SetDefault("LOG_DIR", os.TempDir())
	// max file size of log in megabytes
	viper.SetDefault("LOG_MAX_SIZE", 10)
	// max number of old log files to retain
	viper.SetDefault("LOG_MAX_BACKUPS", 3)
	// max number of days to retain old log files
	viper.SetDefault("LOG_MAX_AGE", 365)

}

// initDirectories creates the output directory & cd's into it
func initDirectories() {
	outPath := viper.GetString("OUTPUT_PATH")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Printf("[MKDIR] %s\n", outPath)
		err = os.MkdirAll(outPath, 0755)
		if err != nil {

			log.Panicf("error making directory: %s", err)
		}
	}

	cmd.Chdir(outPath)
}

func main() {
	// create api client
	auth := newAuth()
	client := api.Client{auth}

	// create & execute api request
	req := client.NewRequest("https://api.bitbucket.org/2.0/repositories/%s", "accuweather")
	pages, err := req.Do()
	if err != nil {
		log.Panicf("unable to make request: %s \n", err)
	}

	// process api response
	for hasNext := true; hasNext; hasNext = pages.Next() {
		p := pages.Current()

		repositories := p["values"].([]interface{})
		for _, v := range repositories {
			name, scmType, cloneLink := getRepoInfo(v)

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
				log.Printf("[PULL] %s\n", name)

				cmd.Chdir(name)
				switch scmType {
				case "hg":
					hg.Update(cloneLink)
				case "git":
					git.Update(cloneLink)
				}
				cmd.Chdir("../")

				if err != nil {
					log.Panicf("error changing directory: %s", err)
				}
			}
		}
	}
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
		if link["name"].(string) == "ssh" {
			cloneLink = link["href"].(string)
		}
	}
	// nevermind, the nightmare is over. now we have a "cloneLink" which is a string

	return
}
