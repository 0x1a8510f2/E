package main

import (
	"fmt"
	//"flag"

	//"github.com/spf13/viper"
	"github.com/TR-SLimey/E/esockets"
	//	"maunium.net/go/mautrix"
)

var (
	// Basic info (static)
	ProjectName = "E"
	ProjectUrl  = "https://github.com/TR-SLimey/E"
	// Incremented on release
	ReleaseVersion = "pre-alpha"
	// Filled at build time
	VcsCommit = "unknown_commit"
	BuildTime = "unknown_build_time"
	// Filled by init
	VersionString = "unknown_version_string"
)

func init() {
	VersionString = fmt.Sprintf("%s %s %s [%s]", ProjectName, ReleaseVersion, BuildTime, VcsCommit)

	// Config handling (TODO)
}

func main() {

	fmt.Println(VersionString)
	println(len(esockets.Available))

}
