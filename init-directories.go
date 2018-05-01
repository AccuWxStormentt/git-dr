package main

import (
	"log"
	"os"

	"bitbucket.org/accuweather/git-dr/cmd"
	"github.com/spf13/viper"
)

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
