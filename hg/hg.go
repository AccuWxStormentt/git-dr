package hg

import (
	"fmt"
	"log"

	"bitbucket.org/accuweather/git-dr/cmd"
	"github.com/spf13/viper"
)

// Clone clones the hg repo at the specified url
func Clone(url string) {
	hg("clone %s", url)
}

// Update updates the hg repo at the specified path
func Update(name string) {
	username := viper.GetString("USERNAME")
	password := viper.GetString("APP_PASSWORD")
	hg("pull -u --config auth.x.prefix=* --config auth.x.username=%s --config auth.x.password=%s", username, password)
}

// hg executes the specified hg command
func hg(s string, args ...interface{}) {
	cmdString := s
	if len(args) > 0 {
		cmdString = fmt.Sprintf(s, args...)
	}
	log.Printf("[HG] %s\n", cmdString)
	cmd.Run("hg", cmdString)
}
