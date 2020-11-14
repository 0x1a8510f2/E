package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/TR-SLimey/E/confmgr"
	"github.com/TR-SLimey/E/esockets"
	"github.com/TR-SLimey/E/strings"
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
	VcsCommit = strings.UNKNOWN_COMMIT
	BuildTime = strings.UNKNOWN_BUILD_TIME
	// Filled by init
	VersionString = strings.UNKNOWN_VERSION_STRING

	// Filled by command line flags
	viewVersion          bool
	printEsockets        bool
	configLocation       string
	registrationLocation string

	// Filled when config is read
	config confmgr.EConfigSkeleton
)

func init() {
	VersionString = fmt.Sprintf(strings.VERSION_STRING, ProjectName, ReleaseVersion, BuildTime, VcsCommit)

	// Handle command-line flags
	flag.BoolVar(&viewVersion, "version", false, strings.FLAG_HELP_VERSION)
	flag.BoolVar(&printEsockets, "esockets", false, strings.FLAG_HELP_ESOCKETS)
	flag.StringVar(&configLocation, "config", "config.yaml", strings.FLAG_HELP_CONFIG)
	flag.StringVar(&registrationLocation, "registration", "none", strings.FLAG_HELP_REGISTRATION)
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
		log.Fatalf(strings.CONFIG_GET_ERR, err.Error())
	}
}

func main() {

	// Log some information on start
	log.Printf(strings.STARTING_WITH_VERSION_STRING, VersionString)
	log.Printf(strings.PROJECT_URL, ProjectUrl)
	log.Printf(strings.ESOCKETS_AVAILABLE_COUNT, len(esockets.Available))

	// Initialise esockets synchronously, and process errors if any
	for esName, es := range esockets.Available {
		log.Printf(strings.ESOCKET_INIT, esName)

		err := es.Init(config.Esockets.ConfDir + "/" + esName + ".yaml")
		if err != nil {
			if config.Esockets.FatalInitFailures {
				log.Fatalf(strings.ESOCKET_INIT_ERR_FATAL, esName, err.Error())
			} else {

				log.Printf(strings.ESOCKET_INIT_ERR_NON_FATAL, esName, err.Error())

				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err := es.Deinit()
				if err != nil {
					log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, esName, err.Error())
				}
			}
		} else {
			// Ensure that the esocket correctly reports as initialised
			err := es.CheckRunlevel(1)
			if err != nil {
				if config.Esockets.FatalInitFailures {
					log.Fatalf(strings.ESOCKET_INIT_ERR_FATAL, esName, err.Error())
				} else {
					log.Printf(strings.ESOCKET_INIT_ERR_NON_FATAL, esName, err)

					// Attempt to deinitialise esocket to save resources. Failures are expected.
					err := es.Deinit()
					if err != nil {
						log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, esName, err.Error())
					}
				}
			}
		}
	}

}
