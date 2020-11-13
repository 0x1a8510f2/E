package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/TR-SLimey/E/configmgr"
	"github.com/TR-SLimey/E/esockets"
)

const (
	// Basic info (static)
	ProjectName = "E"
	ProjectUrl  = "https://github.com/TR-SLimey/E"
	// Incremented on release
	ReleaseVersion = "pre-alpha"
)

var (
	// Filled at build time
	VcsCommit = "unknown_commit"
	BuildTime = "unknown_build_time"
	// Filled by init
	VersionString = "unknown_version_string"

	// Filled by command line flags
	viewVersion          bool
	printEsockets        bool
	configLocation       string
	registrationLocation string

	// Filled when config is read
	config configmgr.EConfigSkeleton
)

func init() {
	VersionString = fmt.Sprintf("%s %s %s [%s]", ProjectName, ReleaseVersion, BuildTime, VcsCommit)

	// Handle command-line flags
	flag.BoolVar(&viewVersion, "version", false, "Print version and exit")
	flag.BoolVar(&printEsockets, "esockets", false, "Print a space-delimeted list of available esockets and exit")
	flag.StringVar(&configLocation, "config", "config.yaml", "The location of the configuration file (YAML format)")
	flag.StringVar(&registrationLocation, "registration", "none", "Where the registration file (YAML config to be placed on the homeserver) should be saved. Values other than `none` imply that the file should be re-/generated")
	flag.Parse()

	// Process command-line flags which end the program to save unnecessary run-time
	if viewVersion {
		fmt.Printf("%s\n", VersionString)
		os.Exit(0)
	} else if printEsockets {
		fmt.Printf("%+v\n", reflect.ValueOf(esockets.Available).MapKeys())
		os.Exit(0)
	}

	// Get the config (automatically check if it's readable and valid)
	var err error
	config, err = configmgr.GetEConfig(configLocation)
	if err != nil {
		log.Fatalf("Error while getting configuration: %s", err.Error())
	}
}

func main() {

	// Log some information on start
	log.Printf("%s starting...", VersionString)
	log.Printf("Project URL: %s", ProjectUrl)
	log.Printf("%d esockets available", len(esockets.Available))

	fmt.Printf("%v\n", config)

	for _, es := range esockets.Available {
		es.Init()
	}

}
