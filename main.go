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
	"github.com/TR-SLimey/E/mxsocket"
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
			var err error
			if action == sr.ESOCKET_ACTION_STOPPING {
				err = es.CheckRunlevel(2)
			} else {
				err = es.CheckRunlevel(1)
			}
			// If err is nil, the current esocket is to have $action performed on it
			if err == nil {
				log.Infof(sr.ESOCKET_GENERIC_ACTION_DESCRIPTION, strings.Title(action), es.ID)
				if action == sr.ESOCKET_ACTION_STOPPING {
					err = es.Stop()
				} else {
					err = es.Deinit()
				}

				if err != nil {
					log.Errorf(sr.ESOCKET_ERR_GENERIC, action, es.ID, err.Error())
				}
			} else {
				log.Infof(sr.ESOCKET_GENERIC_ACTION_DESCRIPTION, strings.Title(sr.ESOCKET_ACTION_SKIPPING)+" "+action, es.ID)
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

	// Create queue (channel) for receiving data from esockets
	esRecvQueue := make(chan map[string]string)

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
				// Inject received data queue channel
				es.RecvQueue = esRecvQueue
				// Init
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

	// Init and start the MxSocket (E<->Matrix interface)
	log.Infof(sr.MATRIX_SOCKET_INIT)
	ms := mxsocket.New()
	err := ms.Init(config.Matrix.RegFilePath)
	if err != nil {
		log.Errorf(sr.MATRIX_SOCKET_INIT_ERR, err.Error())
		triggerCleanExit()
	}
	log.Infof(sr.MATRIX_SOCKET_START)
	err = ms.Start()
	if err != nil {
		log.Errorf(sr.MATRIX_SOCKET_START_ERR, err.Error())
		triggerCleanExit()
	}

	// Maintain a mapping between the client ID and the Esocket they
	// are connected to
	esClientMap := make(map[string]string)

	// Pass data between Matrix and the Esockets
	// This is an infinite loop which can only end
	// when a signal is received or if a panic occurs
	for {
		select {
		case esdata := <-esRecvQueue:
			/*
				Data in the esRecvQueue is a map of strings in the following format:
					"src_esocket": <Esocket.ID>
					"src_client": An ID of the client to allow easy recognition. Should
										not change frequently as each new ID means a new Matrix room.
										The esocket is responsible for determining and keeping
										track of this, and possibly notifying the user of new IDs
										via the recv channel.
					"event_type": Which Matrix event type is to be used to send this data.
										Most often this is simply `m.room.message`.
					"in_main_room": Whether this event should be sent by the main E bot
										in its control room (such events will be marked with
										the name of the esocket they came from and client ID).
										This is useful for notifying about new clients connecting,
										for example, without the need to create a new room. The
										value here will be converted from a string to a boolean,
										using strconv.ParseBool.
					"data": The contents of the Matrix event. This should be a valid JSON string
										which will have little to no validation performed on it
										before being sent, so the ESocket is responsible for
										ensuring the data is valid.
					"ref": Any string reference which the esocket can use to match replies from E
										(like errors) to its original messages. This can also be left
										blank if the esocket doesn't wish to keep track of this.
				Except when the "type" key exists and is not set to "data". In this case,
				it will be interpreted as follows:
					- "client_id_reg" - the client ID in "src_client" will be registered
						to "src_esocket" assuming no other esocket has registered the ID
						or "allowClientIdLocationOverride" is true. If this is not the case,
						an error will be sent to the esocket's send queue.
					- "client_id_unreg" - the client ID in "src_client" will be unregistered
						if it is currently registered to "src_esocket", otherwise the command
						will be silently ignored.
				All other values of "type" will result in an error being sent to the esocket's
				send queue.
			*/

			// Handle non-data message types
			if msgtype, ok := esdata["type"]; ok && msgtype != "data" {
				if msgtype == "client_id_reg" {
					if _, ok := esClientMap[esdata["src_client"]]; ok {
						// The ID is already registered
						if config.Esockets.AllowClientIdLocationOverride {
							// ...and it should get overwritten
							log.Warnf(sr.CLIENT_ID_ALREADY_REGISTERED_OVERWRITE_WARN, esdata["src_esocket"], esdata["src_client"], esClientMap[esdata["src_client"]])
							esClientMap[esdata["src_client"]] = esdata["src_esocket"]
						} else {
							// ...and the registration should be rejected
							log.Warnf(sr.CLIENT_ID_ALREADY_REGISTERED_REJECTION_WARN, esdata["src_esocket"], esdata["src_client"], esClientMap[esdata["src_client"]])
							esockets.Available[esdata["src_esocket"]].SendQueue <- map[string]string{
								"type":        "error",
								"dst_esocket": esdata["src_esocket"],
								"dst_client":  esdata["src_client"],
								"ref":         esdata["ref"],
							}
						}
					}
				} else if msgtype == "client_id_unreg" && esClientMap[esdata["src_client"]] == esdata["src_esocket"] {
					delete(esClientMap, esdata["src_esocket"])
				}
			} else {
				// Pass data to Matrix Socket via channel
				ms.SendQueue <- esdata
			}

		case mxdata := <-ms.RecvQueue:
			/*
				Data in the Matrix socket's RecvQueue is a map of strings in the following format:
					"dst_client": The client this message is intended for. This should simply be
									a string ID as in src_client. If such ID does not exist or has
									disconnected, it is up to the esocket to decide what to do.
					"event_type": Which Matrix event was used to send this message. This should usually
									not be that useful to esockets since most simply forward the message
									regardless, but may be useful in some edge cases or for specific
									requirements.
					"data": The raw content of the Marix message. This should be parsed by the esocket
									itself as, while this makes writing an esocket more difficult,
									it makes them much more versatile and gives them more freedom and
									information.
			*/

			// Look up the esocket mapping of the destination client
			if esId, ok := esClientMap[mxdata["dst_client"]]; ok {
				// ...and send the message to the correct esocket
				esockets.Available[esId].SendQueue <- mxdata
			} else {
				log.Warnf("No esocket mapping was found for client `%s` so the message could not be routed.", mxdata["dst_client"])
			}
		}
	}
}
