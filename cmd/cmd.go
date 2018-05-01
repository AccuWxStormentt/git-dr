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
		log.Fatalf("error changing directory: %s", err)
	}
}

func Run(base, cmdString string) {
	split := strings.Split(cmdString, " ")
	cmd := exec.Command(base, split...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(out))
		log.Fatalf("error executing command: %s", err)
	}
}
