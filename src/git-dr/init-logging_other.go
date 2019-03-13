// +build windows nacl plan9

package main

import (
	"log"

	"github.com/pmaseberg/lumberjack"
	"github.com/spf13/viper"
)

func initLog() {
	var logStdOut = viper.GetBool("logging.stdout")
	var logDir = viper.GetString("logging.dir")
	var logMaxSize = viper.GetInt("logging.max_size")
	var logMaxBackups = viper.GetInt("logging.max_backups")
	var logMaxAge = viper.GetInt("logging.max_age")

	if !logStdOut {
		logSink := &lumberjack.Logger{
			Filename:   logDir + "/" + processName + ".log",
			MaxSize:    logMaxSize,
			MaxBackups: logMaxBackups,
			MaxAge:     logMaxAge,
		}

		log.SetOutput(logSink)
	}
}
