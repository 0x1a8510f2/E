package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/TR-SLimey/E/confmgr"
	"github.com/TR-SLimey/E/esockets"
	log "github.com/TR-SLimey/E/shim/log"
	sr "github.com/TR-SLimey/E/stringres"
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
	VcsCommit = sr.UNKNOWN_COMMIT
	BuildTime = sr.UNKNOWN_BUILD_TIME
	// Filled by init
	VersionString = sr.UNKNOWN_VERSION_STRING

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

	log.Infof(sr.CLEAN_EXIT)

	// Handle follow-up signals to allow force-exit
	go func() {
		<-exitSignalChan
		log.Fatalf(sr.FORCE_EXIT)
	}()

	// Stop running esockets
	for _, es := range esockets.Available {
		err := es.CheckRunlevel(2)
		if err == nil {

			err := es.Stop()
			if err != nil {
				log.Errorf(sr.ESOCKET_STOP_ERR_NON_FATAL, es.ID, err.Error())
			}
		}
	}

	// Deinitialise initialised esockets
	os.Exit(0)
}

func init() {
	VersionString = fmt.Sprintf(sr.VERSION_STRING, ProjectName, ReleaseVersion, BuildTime, VcsCommit)

	// Handle command-line flags
	flag.BoolVar(&viewVersion, "version", false, sr.FLAG_HELP_VERSION)
	flag.BoolVar(&printEsockets, "esockets", false, sr.FLAG_HELP_ESOCKETS)
	flag.StringVar(&configLocation, "config", "config.yaml", sr.FLAG_HELP_CONFIG)
	flag.StringVar(&registrationLocation, "registration", "none", sr.FLAG_HELP_REGISTRATION)
	flag.Parse()

	// Process command-line flags which end the program to save unnecessary run-time
	if viewVersion {
		fmt.Printf("%s\n", VersionString)
		os.Exit(0)
	} else if printEsockets {
		fmt.Printf("%+v\n", reflect.ValueOf(esockets.Available).MapKeys())
		os.Exit(0)
	}

	// Create logger
	log.Init(os.Stdout)

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
		log.Fatalf(sr.CONFIG_GET_ERR, err.Error())
	}
}

func main() {

	// Log some information on start
	log.Infof(sr.STARTING_WITH_VERSION_STRING, VersionString)
	log.Infof(sr.PROJECT_URL, ProjectUrl)
	log.Infof(sr.ESOCKETS_AVAILABLE_COUNT, len(esockets.Available))

	// Initialise and start esockets
	// Both initialisation and starting of esockets are essentially
	// the same code so putting it in a loop and running it twice
	// makes sense
	for _, action := range [2]string{sr.ESOCKET_INITIALISING, sr.ESOCKET_STARTING} {

		// For each esocket being initialised or started depending on $action
		for _, es := range esockets.Available {
			log.Infof("%s `%s` esocket", strings.Title(action), es.ID)

			var err error
			if action == sr.ESOCKET_INITIALISING {
				err = es.Init(config.Esockets.ConfDir + "/" + es.ID + ".yaml")
			} else {
				err = es.Start()
			}

			if err == nil {
				// Ensure that the esocket reports the correct runlevel
				if action == sr.ESOCKET_INITIALISING {
					err = es.CheckRunlevel(1)
				} else {
					err = es.CheckRunlevel(2)
				}
				if err == nil {
					// No errors have occured so move on to next esocket
					continue
				}
			}

			// We haven't hit the continue above so an error has occured
			if config.Esockets.FatalInitFailures {
				if action == sr.ESOCKET_INITIALISING {
					log.Errorf(sr.ESOCKET_INIT_ERR_FATAL, es.ID, err.Error())
				} else {
					log.Errorf(sr.ESOCKET_START_ERR_FATAL, es.ID, err.Error())
				}
				triggerCleanExit()
			} else {
				if action == sr.ESOCKET_INITIALISING {
					log.Warnf(sr.ESOCKET_INIT_ERR_NON_FATAL, es.ID, err)
				} else {
					log.Warnf(sr.ESOCKET_START_ERR_NON_FATAL, es.ID, err)
				}

				if action == sr.ESOCKET_INITIALISING {
					err = es.Stop()
					if err != nil {
						log.Errorf(sr.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
					}
				}
				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err = es.Deinit()
				if err != nil {
					log.Errorf(sr.ESOCKET_DEINIT_ERR_NON_FATAL, es.ID, err.Error())
				}
			}
		}
	}
}
