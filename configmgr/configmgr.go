package configmgr

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// The rough layout of the config file
type ConfigSkeleton struct {
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
