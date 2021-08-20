package github

import (
	"encoding/json"
	"fmt"
	"git-dr/api"
	"git-dr/cmd"
	"git-dr/git"
	"git-dr/repo"
	"log"
	"os"

	"github.com/spf13/viper"
)

type ghPage struct {
	Name    string
	Ssh_url string
	Size    float64
}

func Exists(name string) bool {
	auth := newAuth()
	client := api.Client{auth}
	r := client.NewRequest("https://api.github.com/repos/%s/%s", "AccuWeather-Inc", name)

	_, err := r.Get()
	if err != nil {
		return false
	}

	return true
}
func Create(repo repo.Repo) {
	auth := newAuth()
	client := api.Client{auth}
	r := client.NewRequest("https://api.github.com/orgs/%s/repos", "AccuWeather-Inc")

	resp, err := r.Post(map[string]string{
		"name":        repo.Name,
		"description": repo.Description,
		"has_wiki":    repo.HasWiki,
		"private":     "true",
		"team_id":     viper.GetString("github.team_id"),
	})

	log.Println(resp)

	if err != nil {
		log.Fatal(err)
	}
}

func BackupAll() {
	auth := newAuth()
	client := api.Client{auth}

	total := 0.0
	repoN := 1
	for i := 1; i < 9; i++ {
		r := client.NewRequest("https://api.github.com/orgs/%s/repos?per_page=100&page=%d", "AccuWeather-Inc", i)
		resp, err := r.Get()
		if err != nil {
			log.Fatal(err)
		}

		var pages []ghPage
		err = json.Unmarshal([]byte(resp), &pages)
		if err != nil {
			log.Fatal(err)
		}

		for _, page := range pages {
			size := page.Size
			total += size
			log.Printf("[%d %d] %s: %.0f\n", i, repoN, page.Ssh_url, size)
			repoN++

			if _, err := os.Stat(fmt.Sprintf("%s.git", page.Name)); os.IsNotExist(err) {
				git.Clone(page.Ssh_url)
			} else {
				continue
				cmd.Chdir(fmt.Sprintf("%s.git", page.Name))
				git.Update(page.Ssh_url)
				cmd.Chdir("../")

			}
		}

		if err != nil {
			panic(err)
		}

	}
	log.Printf("[TOTAL] %.0f KB\n", total)
}

func newAuth() api.Authenticator {
	username := viper.GetString("github.username")
	password := viper.GetString("github.password")
	return api.NewBasicAuth(username, password)
}
