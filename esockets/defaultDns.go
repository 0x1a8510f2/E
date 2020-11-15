package esockets

import (
	"github.com/TR-SLimey/E/esockets/defaultDns"
)

func init() {
	// Create the esocket as a local variable
	var esocket = Esocket{
		ID: "defaultDns",
		onInit: func(es *Esocket, confLocation string) error {
			es.Runlevel = 1
			/*err := confmgr.GetEsocketConfig(confLocation, &es.Config)
			if err != nil {
				return fmt.Errorf(strings.ESOCKET_CONFIG_READ_ERR, es.ID, err.Error())
			}*/
			return nil
		},
		onDeinit: func(es *Esocket) error {
			es.Runlevel = 0
			return nil
		},
		onStart: func(es *Esocket) error {
			es.Runlevel = 2
			return nil
		},
		onStop: func(es *Esocket) error {
			es.Runlevel = 1
			return nil
		},
		Config: defaultDns.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
