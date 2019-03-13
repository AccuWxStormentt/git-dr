package bitbucket

import (
	"log"

	"git-dr/api"

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
