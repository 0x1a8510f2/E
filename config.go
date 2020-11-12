package main

type Config struct {
	Matrix struct {
		AsId      string
		Address   string
		BindAddrs []string
		BindPorts []int
		AsToken   string
		HsToken   string
		Bot       struct {
			Username        string
			Displayname     string
			AvatarUrl       string
			Sudoers         []string
			EnabledCommands []string
			NosudoCommands  []string
		}
		Homeserver struct {
			Address      string
			MxidSuffix   string
			Provisioning struct {
				Path         string
				SharedSecret string
			}
		}
		ManagedUsers struct {
			UsernameTemplate    string
			DisplaynameTemplate string
		}
	}
	Esockets struct {
		ConfDir string
	}
}
