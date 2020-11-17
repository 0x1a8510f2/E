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
		ConfDir           string `yaml:"confDir"`
		FatalInitFailures bool   `yaml:"fatalInitFailures"`
	}
}
