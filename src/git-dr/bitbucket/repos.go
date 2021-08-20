package bitbucket

import (
	"fmt"
	"log"
	"strings"

	"git-dr/api"
	"git-dr/repo"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func GetRepos() (*PageCollection, error) {
	auth := newAuth()
	client := api.Client{auth}
	r := client.NewRequest("https://api.bitbucket.org/2.0/repositories/%s", "accuweather")

	pages := &PageCollection{}
	hasNext := true

	for hasNext {
		log.Printf("[API] %s %s\n", r.Method, r.Endpoint)

		respString, err := r.Do()
		if err != nil {
			return nil, errors.Wrap(err, "http request failed")
		}

		page, err := ParsePage(string(respString))
		if err != nil {
			return nil, errors.Wrap(err, "http response parsing failed")
		}

		pages.Add(page)

		r.Endpoint, hasNext = page["next"].(string)
	}

	return pages, nil
}

// newAuth creates a new api.Authenticator. Right now it only creates a BasicAuth object. Later, if it needs to, it can support OAuth2 or whatever else Bitbucket wants to do
func newAuth() api.Authenticator {
	username := viper.GetString("bitbucket.username")
	app_password := viper.GetString("bitbucket.app_password")
	return api.NewBasicAuth(username, app_password)
}

// getRepoInfo casts an interface and pulls out the repo name, scm type (git or hg) & a link to clone it
func GetRepoInfo(v interface{}) (repo repo.Repo) {
	// go's type system is an unending nightmare
	casted := v.(map[string]interface{})
	repo.Name = casted["slug"].(string)
	repo.Description = casted["description"].(string)
	repo.ScmType = casted["scm"].(string)

	links := casted["links"].(map[string]interface{})
	cloneLinks := links["clone"].([]interface{})
	for _, v := range cloneLinks {
		link := v.(map[string]interface{})
		if link["name"].(string) == "https" {
			username := viper.GetString("bitbucket.username")
			password := viper.GetString("bitbucket.app_password")
			combined := fmt.Sprintf("%s:%s", username, password)
			repo.CloneLink = strings.Replace(link["href"].(string), username, combined, -1)
		}
	}

	project := casted["project"].(map[string]interface{})
	project_name := project["key"].(string)
	repo.Project = project_name
	// nevermind, the nightmare is over. we're good now

	return
}
