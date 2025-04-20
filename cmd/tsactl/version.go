package main

import "fmt"

var (
	version = "0.1.0-dev"
	commit  = "none"
)

func getVersion() string {
	return fmt.Sprintf("tsactl v%s commit %s", version, commit)
}
