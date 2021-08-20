package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func Chdir(path string) {
	log.Printf("[CHDIR] %s\n", path)
	err := os.Chdir(path)
	if err != nil {
		log.Printf("error changing directory: %s", err)
	}
}

func Run(base, cmdString string) string {
	split := strings.Split(cmdString, " ")
	cmd := exec.Command(base, split...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(out))
		log.Printf("error executing command: %s", err)
	}

	return string(out)
}
