package strings

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

	ESOCKET_INIT                 = "Initialising `%s` esocket"
	ESOCKET_INIT_ERR_FATAL       = "Error while initialising `%s` esocket: %s"
	ESOCKET_INIT_ERR_NON_FATAL   = "Error while initialising `%s` esocket. This esocket will be deinitialised. Error: %s"
	ESOCKET_DEINIT_ERR_FATAL     = "Deinitialising `%s` esocket failed with error: %s"
	ESOCKET_DEINIT_ERR_NON_FATAL = "Deinitialising `%s` esocket failed with error: %s"
	ESOCKET_START_ERR_FATAL      = ""
	ESOCKET_START_ERR_NON_FATAL  = ""
	ESOCKET_STOP_ERR_FATAL       = ""
	ESOCKET_STOP_ERR_NON_FATAL   = "Stopping `%s` esocket failed with error: %s"

	CONFIG_FILE_OPEN_ERR  = "Could not open config file (%s) for reading! Failed with error: %s"
	CONFIG_FILE_PARSE_ERR = "Could not parse config file (%s). Failed with error: %s"

	INVALID_EXPECTED_RUNLEVEL = "The expected runlevel is invalid."
	UNEXPECTED_RUNLEVEL_ERR   = "Esocket reports as `%s` but `%s` was expected."

	ESOCKET_CONFIG_READ_ERR = "Error reading esocket (%s) config: %s"

	CLEAN_EXIT = "Exiting Cleanly..."
	FORCE_EXIT = "Forcing Exit!"
)
