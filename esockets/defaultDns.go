package esockets

import (
	"fmt"
	"time"

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
				return fmt.Errorf("Error reading esocket (%s) config: %s", es.ID, err.Error())
			}*/
			return nil
		},
		onDeinit: func(es *Esocket) error {
			es.Runlevel = 0
			return nil
		},
		onStart: func(es *Esocket) error {
			es.runFlag = true
			go es.Run(es)
			es.Runlevel = 2
			return nil
		},
		onStop: func(es *Esocket) error {
			es.runFlag = false
			fmt.Println("Waiting for mainloop to exit")
			for es.Runlevel != 1 {
				time.Sleep(5 * time.Microsecond)
			}
			return nil
		},
		Run: func(es *Esocket) {
			for es.runFlag {
				time.Sleep(1 * time.Second)
				fmt.Println("Default DNS Esocket Still Running")
			}
			fmt.Println("Default DNS Esocket Has Exit")
			es.Runlevel = 1
		},
		Config: defaultDns.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
