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
	log "github.com/TR-SLimey/E/shim/log"
	sr "github.com/TR-SLimey/E/stringres"
)

const (
	// Basic info (static)
	ProjectName = "E"
	ProjectUrl  = "https://github.com/TR-SLimey/E"
	// Incremented on release
	ReleaseVersion = "pre-alpha"
	// Name of the esocket which commands come from
	// This must exist else an error will occur
	// The esocket should also be designed with being a
	// master esocket in mind as they have slightly different
	// functionality.
	MasterEsocket = "matrix"
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

	// Ensure the master esocket is loaded and available
	if _, ok := esockets.Available[MasterEsocket]; !ok {
		log.Fatalf(sr.MASTER_ESOCKET_NOT_FOUND_ERR, MasterEsocket)
	}

	// Create queue (channel) for receiving data from esockets
	esOutQueue := make(chan map[string]string)

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
				// Inject received data queue
				es.OutQueue = esOutQueue
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

	// Maintain a mapping between the client ID and the Esocket they
	// are connected to
	esClientMap := make(map[string]string)

	// Pass data between Matrix and the Esockets
	// This is an infinite loop which can only end
	// when a signal is received or if a panic occurs
	for {
		/*
			esOutQueue can contain data from both the master esocket and the regular esockets, but
			there are differences between the two. Data coming from the master esocket should look
			as follows:
				"src_esocket": <Esocket.ID>
				"dst_client": The client this message is intended for. This should simply be
								a string ID as in src_client below. If such ID does not exist or has
								disconnected, it is up to the dst esocket to decide what to do.
				"event_type": Which Matrix event was used to send this message. This should usually
								not be that useful to esockets since most simply forward the message
								regardless, but may be useful in some edge cases or for specific
								requirements.
				"data": The raw content of the Marix message. This should be parsed by the esocket
								itself as, while this makes writing an esocket more difficult,
								it makes them much more versatile and gives them more freedom and
								information.
			Meanwhile, data from regular esockets looks like this:
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
									(like errors) to its original messages. This can be ignored
									in replies if the esocket does not wish to keep track of
									references, but should nevertheless contain some random ID.
			Except when the "type" key exists and is not set to "data". In this case,
			it will be interpreted as follows:
					- "client_id_reg" - the client ID in "src_client" will be registered
						to "src_esocket" assuming no other esocket has registered the ID
						or "allowClientIdLocationOverride" is true. If this is not the case,
						an error will be sent to the esocket's in queue. If the ID is already
						registered by this esocket, this will be ignored silently. No other
						data is required.
					- "client_id_unreg" - the client ID in "src_client" will be unregistered
						if it is currently registered to "src_esocket", otherwise the command
						will be silently ignored. No other data is required.
			All other values of "type" will result in an error being sent to the esocket's
			in queue.
		*/
		esdata := <-esOutQueue

		// Ensure a valid "return address" is included in the data
		if sourceEsId, ok := esdata["src_esocket"]; !ok {
			log.Warnf(sr.NO_SRC_ESOCKET_ERR)

		} else if sourceEsId == MasterEsocket {
			// Is master esocket

			// Ensure data is valid
			okTotal := true
			for _, key := range []string{"dst_client", "event_type", "data"} {
				if _, ok := esdata[key]; !ok {
					okTotal = false
				}
			}

			if !okTotal {
				log.Warnf(sr.MALFORMED_DATA_FROM_ESOCKET_ERR, sourceEsId)
			} else {
				// Look up the esocket mapping of the destination client
				if esId, ok := esClientMap[esdata["dst_client"]]; ok {
					// ...and send the message to the correct esocket
					esockets.Available[esId].InQueue <- esdata
				} else {
					log.Warnf(sr.NO_ESOCKET_MAPPING_FOR_CLIENT_ERR, esdata["dst_client"], esdata["ref"])
				}
			}

		} else {
			// Is regular esocket

			// Handle non-data message types
			if msgtype, ok := esdata["type"]; ok && msgtype != "data" {
				if msgtype == "client_id_reg" {
					// Registering client ID to esocket

					// Ensure data is valid
					okTotal := true
					for _, key := range []string{"src_client", "ref"} {
						if _, ok := esdata[key]; !ok {
							okTotal = false
						}
					}

					if !okTotal {
						// The data is invalid
						log.Warnf(sr.MALFORMED_DATA_FROM_ESOCKET_ERR, sourceEsId)
						esockets.Available[sourceEsId].InQueue <- map[string]string{
							"type":        "error",
							"err_code":    "3",
							"err_msg":     "",
							"dst_esocket": esdata["src_esocket"],
							"dst_client":  esdata["src_client"],
							"ref":         esdata["ref"],
						}
					} else if _, ok := esClientMap[esdata["src_client"]]; okTotal && ok {
						// The ID is already registered
						if config.Esockets.AllowClientIdLocationOverride {
							// ...and it should get overwritten
							log.Warnf(sr.CLIENT_ID_ALREADY_REGISTERED_OVERWRITE_WARN, sourceEsId, esdata["src_client"], esClientMap[esdata["src_client"]])
							esClientMap[esdata["src_client"]] = sourceEsId
						} else {
							// ...and the registration should be rejected
							log.Warnf(sr.CLIENT_ID_ALREADY_REGISTERED_REJECTION_WARN, sourceEsId, esdata["src_client"], esClientMap[esdata["src_client"]])
							esockets.Available[sourceEsId].InQueue <- map[string]string{
								"type":        "error",
								"err_code":    "5",
								"err_msg":     "",
								"dst_esocket": esdata["src_esocket"],
								"dst_client":  esdata["src_client"],
								"ref":         esdata["ref"],
							}
						}
					} else {
						// The client ID is not yet registered
						esClientMap[esdata["src_client"]] = sourceEsId
					}
				} else if msgtype == "client_id_unreg" {
					// Unregistering client ID from esocket

					// Ensure data is valid
					okTotal := true
					for _, key := range []string{"src_client", "ref"} {
						if _, ok := esdata[key]; !ok {
							okTotal = false
						}
					}

					if !okTotal {
						// The data is invalid
						log.Warnf(sr.MALFORMED_DATA_FROM_ESOCKET_ERR, sourceEsId)
						esockets.Available[sourceEsId].InQueue <- map[string]string{
							"type":        "error",
							"err_code":    "3",
							"err_msg":     "",
							"dst_esocket": esdata["src_esocket"],
							"dst_client":  esdata["src_client"],
							"ref":         esdata["ref"],
						}
					} else if esClientMap[esdata["src_client"]] == esdata["src_esocket"] {
						// The data is valid and the client belongs to the sending esocket
						// NOTE: The above check relies on on esockets being truthful about who they
						// are, but we assume esockets are trustworthy, since running an untrustworthy
						// esocket means arbitrary code execution anyway, which is a bigger concern
						// than clients being unregistered maliciously.
						delete(esClientMap, esdata["src_esocket"])
					}
				} else {
					// Invalid data type
					log.Warnf(sr.INVALID_DATA_TYPE_IN_CHANNEL_ERR, msgtype)
					esockets.Available[sourceEsId].InQueue <- map[string]string{
						"type":        "error",
						"err_code":    "6",
						"err_msg":     fmt.Sprintf(sr.INVALID_DATA_TYPE_IN_CHANNEL_RETURN_MSG, msgtype),
						"dst_esocket": esdata["src_esocket"],
						"dst_client":  esdata["src_client"],
						"ref":         esdata["ref"],
					}
				}
			} else {
				// Ensure data is valid
				okTotal := true
				for _, key := range []string{"src_client", "event_type", "in_main_room", "data", "ref"} {
					if _, ok := esdata[key]; !ok {
						okTotal = false
					}
				}

				if !okTotal {
					log.Warnf(sr.MALFORMED_DATA_FROM_ESOCKET_ERR, sourceEsId)
					esockets.Available[sourceEsId].InQueue <- map[string]string{
						"type":        "error",
						"err_code":    "3",
						"err_msg":     "",
						"dst_esocket": esdata["src_esocket"],
						"dst_client":  esdata["src_client"],
						"ref":         esdata["ref"],
					}
				} else {
					// Pass data to master esocket via its in queue
					esockets.Available[MasterEsocket].InQueue <- esdata
				}
			}
		}
	}
}
