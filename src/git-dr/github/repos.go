package github

import (
	"git-dr/api"
	"log"

	"github.com/spf13/viper"
)

type Repo struct {
	Name    string
	Desc    string
	HasWiki string
}

func Create(repo *Repo) {
	auth := newAuth()
	client := api.Client{auth}
	r := client.NewRequest("https://api.github.com/orgs/%s/repos", "AccuWeather-Inc")

	resp, err := r.Post(map[string]string{
		"name":        repo.Name,
		"description": repo.Desc,
		"has_wiki":    repo.HasWiki,
		"private":     "true",
		"team_id":     viper.GetString("github.team_id"),
	})

	log.Println(resp)

	if err != nil {
		log.Fatal(err)
	}
}

func newAuth() api.Authenticator {
	username := viper.GetString("github.username")
	password := viper.GetString("github.password")
	return api.NewBasicAuth(username, password)
}
