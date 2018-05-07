// +build !windows,!nacl,!plan9

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pmaseberg/lumberjack"
	"github.com/spf13/viper"
)

func initLog() {
	var logStdOut = viper.GetBool("LOG_STD")
	var logDir = viper.GetString("LOG_DIR")
	var logMaxSize = viper.GetInt("LOG_MAX_SIZE")
	var logMaxBackups = viper.GetInt("LOG_MAX_BACKUPS")
	var logMaxAge = viper.GetInt("LOG_MAX_AGE")

	if !logStdOut {
		logSink := &lumberjack.Logger{
			Filename:   logDir + "/" + processName + ".log",
			MaxSize:    logMaxSize,
			MaxBackups: logMaxBackups,
			MaxAge:     logMaxAge,
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGUSR1)

		go func() {
			for {
				<-sigCh
				logSink.Rotate()
			}
		}()

		log.SetOutput(logSink)
	}
}
