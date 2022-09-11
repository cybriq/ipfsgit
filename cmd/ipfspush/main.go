package main

import (
	"log"
	"os"
	"os/exec"
)

var (
	Info *log.Logger
	// Warning *log.Logger
	Error *log.Logger
)

func initLog() {
	infoHandle, errorHandle := os.Stderr, os.Stderr
	Info = log.New(infoHandle, "INFO ",
		log.Ldate|log.Ltime|log.Lshortfile)
	// Warning = log.New(warningHandle, "WARNING ",
	// 	log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLog()

	cmd := exec.Command("ipfs")
	err := cmd.Run()
	if err != nil {
		Error.Fatal(err,
			"\n\nYou need to install kubo to use this tool"+
				" (and it needs to be running as a service "+
				"to push your repo)\n\n")

	}

}
