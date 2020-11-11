package main

type Config struct {
	Matrix struct {
		AsId       string
		Address    string
		BindAddr   string
		BindPort   int
		AsToken    string
		HsToken    string
		Sudoers    []string
		Homeserver struct {
			Address string
			Domain  string
		}
		Provisioning struct {
			Path         string
			SharedSecret string
		}
		Bot struct {
			Username    string
			Displayname string
			AvatarUrl   string
		}
		ManagedUsers struct {
			UsernameTemplate    string
			DisplaynameTemplate string
		}
	}
	Esockets struct {
		// TODO Esocket config
	}
}
