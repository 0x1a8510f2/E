package esockets

import (
	"github.com/TR-SLimey/E/esockets/defaultDns"
)

func init() {
	// Create the esocket as a local variable
	var esocket = Esocket{
		ID: "defaultDns",
		onInit: func(es *Esocket) {
			println(es.ID)
		},
		Config: defaultDns.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
