package confmgr

import (
	"fmt"
	"io/ioutil"

	"github.com/TR-SLimey/E/confmgr/template"
	yaml "github.com/TR-SLimey/E/shim/yaml"
	sr "github.com/TR-SLimey/E/stringres"
)

func GetEConfig(location string) (template.EConfig, error) {
	var config template.EConfig

	data, err := ioutil.ReadFile(location)
	if err != nil {
		return config, fmt.Errorf(sr.CONFIG_FILE_OPEN_ERR, location, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf(sr.CONFIG_FILE_PARSE_ERR, location, err)
	}

	return config, nil
}

func GetEsocketConfig(location string, confVarPtr *struct{}) error {
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return fmt.Errorf(sr.CONFIG_FILE_OPEN_ERR, location, err)
	}

	err = yaml.Unmarshal(data, confVarPtr)
	if err != nil {
		return fmt.Errorf(sr.CONFIG_FILE_PARSE_ERR, location, err)
	}

	return nil
}
