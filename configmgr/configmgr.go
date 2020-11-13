package configmgr

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// The rough layout of the config file
type ConfigSkeleton struct {
	Matrix struct {
		AsId      string   `yaml:"asId"`
		Address   string   `yaml:"address"`
		BindAddrs []string `yaml:"bindAddrs"`
		BindPorts []int    `yaml:"bindPorts"`
		AsToken   string   `yaml:"asToken"`
		HsToken   string   `yaml:"hsToken"`
		Bot       struct {
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
		ConfDir string `yaml:"confDir"`
	}
}

func GetConfig(location string) (ConfigSkeleton, error) {
	var config ConfigSkeleton

	data, err := ioutil.ReadFile(location)
	if err != nil {
		return config, fmt.Errorf("Could not open config file (%s) for reading! Failed with error: %s", location, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("Could not parse config file (%s). Failed with error: %s", location, err)
	}

	return config, nil
}
