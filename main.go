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
	conftemplate "github.com/TR-SLimey/E/confmgr/template"
	"github.com/TR-SLimey/E/esockets"
	"github.com/TR-SLimey/E/matrixsocket"
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
	config conftemplate.EConfig

	// Channel for handling exit signals
	exitSignalChan = make(chan os.Signal)
)

func triggerCleanExit() {
	// Send a signal down the exit signal channel to trigger cleanExit
	exitSignalChan <- syscall.SIGUSR2
}

func setupCleanExit() {
	// Wait for signal while running in the background
	sig := <-exitSignalChan

	if sig.String() == "user defined signal 2" {
		log.Infof(sr.CLEAN_EXIT_TRIGGERED)
	} else {
		log.Infof(sr.CLEAN_EXIT_ON_SIGNAL, sig.String())
	}

	// Handle follow-up signals to allow force-exit
	go func() {
		<-exitSignalChan
		log.Fatalf(sr.FORCE_EXIT)
	}()

	// Ensure that the action strings haven't been tweaked to be
	// the same because that breaks some of the logic
	if sr.ESOCKET_ACTION_INITIALISING == sr.ESOCKET_ACTION_STARTING {
		log.Fatalf(sr.STOPPING_IS_DEINITIALISING_ERR)
	}

	// Stop and deinitialise running esockets
	for _, action := range [2]string{sr.ESOCKET_ACTION_STOPPING, sr.ESOCKET_ACTION_DEINITIALISING} {

		// For each esocket being stopped or deinitialised depending on $action
		for _, es := range esockets.Available {
			log.Infof("%s `%s` esocket", strings.Title(action), es.ID)

			var err error
			if action == sr.ESOCKET_ACTION_STOPPING {
				err = es.CheckRunlevel(2)
			} else {
				err = es.CheckRunlevel(1)
			}
			// If err is nil, the current esocket is to have $action performed on it
			if err == nil {
				if action == sr.ESOCKET_ACTION_STOPPING {
					err = es.Stop()
				} else {
					err = es.Deinit()
				}

				if err != nil {
					log.Errorf(sr.ESOCKET_ERR_GENERIC, action, es.ID, err.Error())
				}
			}
		}
	}

	log.Infof(sr.CLEAN_EXIT_DONE)
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

	// Process command-line flags which instantly exit to save unnecessary run-time
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
		syscall.SIGTERM, // kill XXXX
	)
	go setupCleanExit()

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

	/* Initialise and start esockets
	Both initialisation and starting of esockets are essentially
	the same code so putting it in a loop and running it twice
	makes sense. */

	// First, ensure that the strings haven't been tweaked to be
	// the same because that breaks some of the logic
	if sr.ESOCKET_ACTION_INITIALISING == sr.ESOCKET_ACTION_STARTING {
		log.Fatalf(sr.INITIALISING_IS_STARTING_ERR)
	}
	// Run the actual loop
	for _, action := range [2]string{sr.ESOCKET_ACTION_INITIALISING, sr.ESOCKET_ACTION_STARTING} {

		// For each esocket being initialised or started depending on $action
		for _, es := range esockets.Available {
			log.Infof("%s `%s` esocket", strings.Title(action), es.ID)

			var err error
			if action == sr.ESOCKET_ACTION_INITIALISING {
				err = es.Init(config.Esockets.ConfDir + "/" + es.ID + ".yaml")
			} else {
				err = es.Start()
			}

			if err == nil {
				// Ensure that the esocket reports the correct runlevel
				if action == sr.ESOCKET_ACTION_INITIALISING {
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
				log.Errorf(sr.ESOCKET_ERR_FATAL, action, es.ID, err.Error())
				triggerCleanExit()
			} else {
				log.Warnf(sr.ESOCKET_ERR_NON_FATAL, action, es.ID, err)

				if action == sr.ESOCKET_ACTION_STARTING {
					err = es.Stop()
					if err != nil {
						log.Errorf(sr.ESOCKET_ERR_GENERIC, sr.ESOCKET_ACTION_STOPPING, es.ID, err.Error())
					}
				}
				// Attempt to deinitialise esocket to save resources. Failures are expected.
				err = es.Deinit()
				if err != nil {
					log.Errorf(sr.ESOCKET_ERR_NON_FATAL, sr.ESOCKET_ACTION_DEINITIALISING, es.ID, err.Error())
				}
			}
		}
	}

	// Init and start the E<->Matrix interface
	log.Infof(sr.MATRIX_SOCKET_INIT)
	err := matrixsocket.Init(config.Matrix.RegFilePath)
	if err != nil {
		log.Errorf(sr.MATRIX_SOCKET_INIT_ERR, err.Error())
		triggerCleanExit()
	}
	log.Infof(sr.MATRIX_SOCKET_START)
	err = matrixsocket.Start()
	if err != nil {
		log.Errorf(sr.MATRIX_SOCKET_START_ERR, err.Error())
		triggerCleanExit()
	}

	// Pass data between Matrix and the Esockets
	// This is an infinite loop which can only end
	// when a signal is received or if a panic occurs
	for {
		for _, es := range esockets.Available {
			select {
			case data := <-es.CtrlChannel:
				fmt.Println(data)
			case data := <-es.DataChannel:
				fmt.Println(data)
			}
		}
	}
}
