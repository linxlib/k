package main

import (
	"github.com/linxlib/k/cmd"
	"github.com/linxlib/k/utils/innerlog"
	"github.com/linxlib/k/utils/rundaemon"
	"os"
)

func main() {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				innerlog.Log.Error(err.Error())
			} else {
				panic(exception)
			}
		}
	}()
	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "init", "i":
		cmd.Initialize()
	case "build":
		cmd.Build()
	case "run":
		rundaemon.Run()
	default:
		rundaemon.Run()
	}
}
