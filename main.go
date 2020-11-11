package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/TR-SLimey/E/esockets"
	"github.com/smallfish/simpleyaml"
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

	// Filled by command line flags
	viewVersion          bool
	configLocation       string
	registrationLocation string
)

func init() {
	VersionString = fmt.Sprintf("%s %s %s [%s]", ProjectName, ReleaseVersion, BuildTime, VcsCommit)

	// Handle command-line flags
	flag.BoolVar(&viewVersion, "version", false, "Print version and exit")
	flag.StringVar(&configLocation, "config", "config.yaml", "The location of the configuration file (YAML format)")
	flag.StringVar(&registrationLocation, "registration", "none", "Where the registration file (YAML config to be placed on the homeserver) should be saved. Values other than `none` imply that the file should be re-/generated")
	flag.Parse()

	// Check, read and parse the config (TODO)
	data, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Fatalf("Could not open config file (%s) for reading! Failed with error: %s", configLocation, err)
	}
	_, err = simpleyaml.NewYaml(data)
	if err != nil {
		log.Fatalf("Could not parse config file (%s). Failed with error: %s", configLocation, err)
	}
}

func main() {

	log.Printf("Starting... <%s>", VersionString)
	println(len(esockets.Available))

}
