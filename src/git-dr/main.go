package main

import (
	"git-dr/bitbucket"
	"git-dr/cmd"
	"git-dr/git"
	"git-dr/github"
	"git-dr/hg"
	"git-dr/repo"
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

	github.BackupAll()
	os.Exit(0)
}

func clone(repo repo.Repo) {
	log.Printf("[CLONE] %s\n", repo.CloneLink)
	switch repo.ScmType {
	case "hg":
		hg.Clone(repo.CloneLink)
	case "git":
		git.Clone(repo.CloneLink)
	}
}

func update(repo repo.Repo) {
	log.Printf("[UPDATE] %s\n", repo.Name)

	cmd.Chdir(repo.Name)
	switch repo.ScmType {
	case "hg":
		hg.Update(repo.Name)
	case "git":
		git.Update(repo.Name)
	}
	cmd.Chdir("../")
}

func migrate(repo repo.Repo) {
	cmd.Chdir(repo.Name)
	log.Printf("[MIGRATE] %s\n", repo.Name)
	if !github.Exists(repo.Name) {
		github.Create(repo)
		git.AddRemote(repo.Name, viper.GetString("github.username"), viper.GetString("github.password"))
		git.Push()
	}
	cmd.Chdir("../")
}

func contains(s []string, v string) bool {
	for _, cur := range s {
		if cur == v {
			return true
		}
	}

	return false
}
