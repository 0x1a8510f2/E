package esockets

import (
	"github.com/TR-SLimey/E/esockets/defaultHttp"
)

func init() {
	// Create the esocket as a local variable
	var esocket = Esocket{
		ID: "defaultHttp",
		onInit: func(es *Esocket) {
			println(es.ID)
		},
		Config: defaultHttp.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
