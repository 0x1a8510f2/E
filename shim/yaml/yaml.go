package log

import (
	"gopkg.in/yaml.v3"
)

func Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}

func Unmarshal(in []byte, out interface{}) (err error) {
	return yaml.Unmarshal(in, out)
}
