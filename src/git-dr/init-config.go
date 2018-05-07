package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// initConfig loads the viper configuration from git-dr.[json|yml|toml]
func initConfig() {
	viper.SetConfigName("git-dr")

	viper.AddConfigPath("/etc/git-dr/")
	viper.AddConfigPath("$HOME/.git-dr/")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to read config file: %s", err)
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
