package template

type EConfig struct {
	Matrix struct {
		AsId        string   `yaml:"asId"`
		Address     string   `yaml:"address"`
		BindAddrs   []string `yaml:"bindAddrs"`
		BindPorts   []int    `yaml:"bindPorts"`
		AsToken     string   `yaml:"asToken"`
		HsToken     string   `yaml:"hsToken"`
		RegFilePath string   `yaml:"regFilePath"`
		Bot         struct {
			Username        string   `yaml:"username"`
			Displayname     string   `yaml:"displayname"`
			AvatarUrl       string   `yaml:"avatarUrl"`
			Sudoers         []string `yaml:"sudoers"`
			EnabledCommands []string `yaml:"enabledCommands"`
			NosudoCommands  []string `yaml:"nosudoCommands"`
		} `yaml:"bot"`
		Homeserver struct {
			Address      string `yaml:"address"`
			MxidSuffix   string `yaml:"mxidSuffix"`
			Provisioning struct {
				Path         string `yaml:"path"`
				SharedSecret string `yaml:"sharedSecret"`
			} `yaml:"provisioning"`
		} `yaml:"homeserver"`
		ManagedUsers struct {
			UsernameTemplate    string `yaml:"usernameTemplate"`
			DisplaynameTemplate string `yaml:"displaynameTemplate"`
		} `yaml:"managedUsers"`
	} `yaml:"matrix"`
	Esockets struct {
		ConfDir                       string `yaml:"confDir"`
		FatalInitFailures             bool   `yaml:"fatalInitFailures"`
		AllowClientIdLocationOverride bool   `yaml:"allowClientIdLocationOverride"`
	}
}

func (config *EConfig) SetDefaults() {
	config.Matrix.AsId = "E"
	config.Matrix.Address = "https://E.example.com:12345"
	config.Matrix.BindAddrs = []string{"0.0.0.0"}
	config.Matrix.BindPorts = []int{8080}
	config.Matrix.AsToken = ""
	config.Matrix.HsToken = ""
	config.Matrix.RegFilePath = "registration.yaml"
	config.Matrix.Bot.Username = "E"
	config.Matrix.Bot.Displayname = "E"
	config.Matrix.Bot.AvatarUrl = ""
	config.Matrix.Bot.Sudoers = []string{}
	config.Matrix.Bot.EnabledCommands = []string{}
	config.Matrix.Bot.NosudoCommands = []string{}
	config.Matrix.Homeserver.Address = "https://matrix.example.com"
	config.Matrix.Homeserver.MxidSuffix = "example.com"
	config.Matrix.Homeserver.Provisioning.Path = "/_matrix/provision/v1"
	config.Matrix.Homeserver.Provisioning.SharedSecret = "disable"
	config.Matrix.ManagedUsers.UsernameTemplate = "e_{{ConnectionId}}"
	config.Matrix.ManagedUsers.DisplaynameTemplate = "e_{{ConnectionId}}"
	config.Esockets.ConfDir = "./esocket-conf"
	config.Esockets.FatalInitFailures = true
	config.Esockets.AllowClientIdLocationOverride = false
}
