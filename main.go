package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

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

	// Channel for handling exit signals
	exitSignalChan = make(chan os.Signal, 1)
)

func triggerCleanExit() {
	// Send a signal down the exit signal channel to trigger doCleanExit
	exitSignalChan <- syscall.SIGHUP
}

func doCleanExit() {
	// Wait for signal while running in the background
	<-exitSignalChan

	log.Println(strings.CLEAN_EXIT)

	// Handle follow-up signals to allow force-exit
	go func() {
		<-exitSignalChan
		log.Fatalf(strings.FORCE_EXIT)
	}()

	// Stop running esockets
	for _, es := range esockets.Available {
		err := es.CheckRunlevel(2)
		if err == nil {

			err := es.Stop()
			if err != nil {
				log.Printf(strings.ESOCKET_STOP_ERR_NON_FATAL, es.ID, err.Error())
			}
		}
	}

	// Deinitialise initialised esockets
	os.Exit(0)
}

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

	// Register signal handler to exit gracefully
	signal.Notify(
		exitSignalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)
	go doCleanExit()

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
	for _, es := range esockets.Available {
		log.Printf(strings.ESOCKET_INIT, es.ID)

		err := es.Init(config.Esockets.ConfDir + "/" + es.ID + ".yaml")
		if err == nil {
			// Ensure that the esocket correctly reports as initialised
			err := es.CheckRunlevel(1)
			if err != nil {
				if config.Esockets.FatalInitFailures {
					log.Fatalf(strings.ESOCKET_INIT_ERR_FATAL, es.ID, err.Error())
				} else {
					log.Printf(strings.ESOCKET_INIT_ERR_NON_FATAL, es.ID, err)

					// Attempt to deinitialise esocket to save resources. Failures are expected.
					err := es.Deinit()
					if err != nil {
						log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
					}
				}
			}
		} else {
			if config.Esockets.FatalInitFailures {
				log.Fatalf(strings.ESOCKET_INIT_ERR_FATAL, es.ID, err.Error())
			} else {

				log.Printf(strings.ESOCKET_INIT_ERR_NON_FATAL, es.ID, err.Error())

				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err := es.Deinit()
				if err != nil {
					log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
				}
			}
		}
	}

	// Asynchronously start the esockets, and process errors if any
	for _, es := range esockets.Available {
		log.Printf(strings.ESOCKET_START, es.ID)

		err := es.Init(config.Esockets.ConfDir + "/" + es.ID + ".yaml")
		if err == nil {
			// Ensure that the esocket correctly reports as initialised
			err := es.CheckRunlevel(1)
			if err != nil {
				if config.Esockets.FatalInitFailures {
					log.Fatalf(strings.ESOCKET_START_ERR_FATAL, es.ID, err.Error())
				} else {
					log.Printf(strings.ESOCKET_START_ERR_NON_FATAL, es.ID, err)

					// Attempt to deinitialise esocket to save resources. Failures are expected.
					err := es.Deinit()
					if err != nil {
						log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
					}
				}
			}
		} else {
			if config.Esockets.FatalInitFailures {
				log.Fatalf(strings.ESOCKET_START_ERR_FATAL, es.ID, err.Error())
			} else {

				log.Printf(strings.ESOCKET_START_ERR_NON_FATAL, es.ID, err.Error())

				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err := es.Deinit()
				if err != nil {
					log.Printf(strings.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
				}
			}
		}
	}
}
