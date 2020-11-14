package esockets

import (
	"fmt"

	"github.com/TR-SLimey/E/configmgr"
	"github.com/TR-SLimey/E/esockets/defaultHttp"
)

func init() {
	// Create the esocket as a local variable
	var esocket = Esocket{
		ID: "defaultHttp",
		onInit: func(es *Esocket, confLocation string) error {
			err := configmgr.GetEsocketConfig(confLocation, &es.Config)
			if err != nil {
				return fmt.Errorf("Error reading esocket (%s) config: %s", es.ID, err.Error())
			}
			return nil
		},
		Config: defaultHttp.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
