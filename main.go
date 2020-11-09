package main

import (
	"fmt"
	"github.com/TR-SLimey/E/esockets"
//	"maunium.net/go/mautrix"
)

var (
	// Basic info (static)
	ProjectName = "E"
	ProjectUrl = "https://github.com/TR-SLimey/E"
	// Incremented on release
	ReleaseVersion = "pre-alpha"
	// Filled at build time
	VcsCommit = "unknown"
	BuildTime = "unknown"
	// Filled by init
	VersionString = "unknown"
)

func init() {
	VersionString = fmt.Sprintf("%s %s %s [%s]", ProjectName, ReleaseVersion, BuildTime, VcsCommit)
}

func main() {

	fmt.Println(VersionString)
	println(len(esockets.Available))

}
