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
