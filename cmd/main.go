package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mospany/minicni/pkg/args"
	"github.com/mospany/minicni/pkg/handler"
	"log"
)

const (
	IPStore = "/tmp/reserved_ips"
	LogFile = "/var/log/minicni.log"
)

func setupLogger() {
	logFileLocation, _ := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	log.SetOutput(logFileLocation)
}

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
}

func main() {
	setupLogger()

	cmd, cmdArgs, err := args.GetArgsFromEnv()
	if err != nil {
		log.Fatalf("getting cmd arguments with error: %v", err)
		os.Exit(1)
	}

	fh := handler.NewFileHandler(IPStore)

	switch cmd {
	case "ADD":
		err = fh.HandleAdd(cmdArgs)
	case "DEL":
		err = fh.HandleDel(cmdArgs)
	case "CHECK":
		err = fh.HandleCheck(cmdArgs)
	case "VERSION":
		err = fh.HandleVersion(cmdArgs)
	default:
		err = fmt.Errorf("unknown CNI_COMMAND: %s", cmd)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to handle CNI_COMMAND %q: %v", cmd, err)
		os.Exit(1)
	}
}
