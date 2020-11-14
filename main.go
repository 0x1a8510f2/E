package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/TR-SLimey/E/confmgr"
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
	config confmgr.EConfigSkeleton
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
	config, err = confmgr.GetEConfig(configLocation)
	if err != nil {
		log.Fatalf("Error while getting configuration: %s", err.Error())
	}
}

func main() {

	// Log some information on start
	log.Printf("%s starting...", VersionString)
	log.Printf("Project URL: %s", ProjectUrl)
	log.Printf("%d esocket(s) available", len(esockets.Available))

	// Initialise esockets synchronously, and process errors if any
	for esName, es := range esockets.Available {
		log.Printf("Initialising `%s` esocket", esName)

		err := es.Init(config.Esockets.ConfDir + "/" + esName + ".yaml")
		if err != nil {
			if config.Esockets.FatalInitFailures {
				log.Fatalf("Error while initialising `%s` esocket: %s", esName, err.Error())
			} else {

				log.Printf("Error while initialising `%s` esocket. This esocket will be deinitialised. Error: %s", esName, err.Error())

				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err := es.Deinit()
				if err != nil {
					log.Printf("Deinitialising `%s` esocket failed with error: %s", esName, err.Error())
				}
			}
		} else {
			// Ensure that the esocket correctly reports as initialised
			err := es.CheckRunlevel(1)
			if err != nil {
				if config.Esockets.FatalInitFailures {
					log.Fatalf("Error while initialising `%s` esocket: %s", esName, err)
				} else {
					log.Printf("Error while initialising `%s` esocket. This esocket will be deinitialised. Error: %s", esName, err)

					// Attempt to deinitialise esocket to save resources. Failures are expected.
					err := es.Deinit()
					if err != nil {
						log.Printf("Deinitialising `%s` esocket failed with error: %s", esName, err.Error())
					}
				}
			}
		}
	}

}
