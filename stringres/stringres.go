package stringres

const (
	UNKNOWN_COMMIT         = "unknown_commit"
	UNKNOWN_BUILD_TIME     = "unknown_build_time"
	UNKNOWN_VERSION_STRING = "unknown_version_string"

	VERSION_STRING = "%s %s %s [%s]"

	FLAG_HELP_VERSION      = "Print version and exit"
	FLAG_HELP_ESOCKETS     = "Print a space-delimeted list of available esockets and exit"
	FLAG_HELP_CONFIG       = "The location of the configuration file (YAML format)"
	FLAG_HELP_REGISTRATION = "Where the registration file (YAML config to be placed on the homeserver) should be saved. Values other than `none` imply that the file should be re-/generated"

	CONFIG_GET_ERR = "Error while getting configuration: %s"

	STARTING_WITH_VERSION_STRING = "%s starting..."
	PROJECT_URL                  = "Project URL: %s"
	ESOCKETS_AVAILABLE_COUNT     = "%d esocket(s) available"

	MASTER_ESOCKET_NOT_FOUND_ERR = "The master esocket `%s` was not found; cannot continue."

	ESOCKET_ACTION_INITIALISING   = "initialising"
	ESOCKET_ACTION_STARTING       = "starting"
	ESOCKET_ACTION_STOPPING       = "stopping"
	ESOCKET_ACTION_DEINITIALISING = "deinitialising"
	ESOCKET_ACTION_SKIPPING       = "skipping"

	ESOCKET_GENERIC_ACTION_DESCRIPTION = "%s `%s` esocket"
	ESOCKET_ERR_FATAL                  = "Fatal error while %s `%s` esocket: %s"
	ESOCKET_ERR_NON_FATAL              = "Error while %s `%s` esocket. This esocket will be deinitialised and ignored. Error: %s"
	ESOCKET_ERR_GENERIC                = "Error while %s `%s` esocket: %s"
	INITIALISING_IS_STARTING_ERR       = "A terrible mistake has occured and initialising is the same as starting :O"
	STOPPING_IS_DEINITIALISING_ERR     = "A terrible mistake has occured and stopping is the same as deinitialising :O"

	CONFIG_FILE_OPEN_ERR  = "Could not open config file (%s) for reading! Failed with error: %s"
	CONFIG_FILE_PARSE_ERR = "Could not parse config file (%s). Failed with error: %s"

	INVALID_EXPECTED_RUNLEVEL = "The expected runlevel is invalid."
	UNEXPECTED_RUNLEVEL_ERR   = "Esocket reports as `%s` but `%s` was expected."

	MATRIX_SOCKET_INIT      = "Initialising Matrix socket"
	MATRIX_SOCKET_START     = "Starting Matrix socket"
	MATRIX_SOCKET_INIT_ERR  = "A fatal error has occured while initialising the Matrix Socket: %s"
	MATRIX_SOCKET_START_ERR = "A fatal error has occured while starting the Matrix Socket: %s"

	NO_SRC_ESOCKET_ERR                          = "An event was received but the origin esocket could not be determined. Dropping the event"
	MALFORMED_DATA_FROM_ESOCKET_ERR             = "An event was received from esocket `%s` but it was malformed. Dropping the event"
	CLIENT_ID_ALREADY_REGISTERED_REJECTION_WARN = "Esocket `%s` attempted to register client ID `%s` but the request was rejected because the ID is already registered by esocket `%s`"
	CLIENT_ID_ALREADY_REGISTERED_OVERWRITE_WARN = "Esocket `%s` re-registered client ID `%s` from esocket `%s`"
	NO_ESOCKET_MAPPING_FOR_CLIENT_ERR           = "No esocket mapping was found for client `%s` so the message with ID `%s` could not be routed."

	CLEAN_EXIT_ON_SIGNAL = "Signal received (%s). Exiting cleanly (re-send to force exit)..."
	CLEAN_EXIT_TRIGGERED = "Exit requested by application. Exiting cleanly..."
	CLEAN_EXIT_DONE      = "Exit"
	FORCE_EXIT           = "Received follow-up signal. Forcing exit!"
)
