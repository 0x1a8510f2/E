package esockets

import (
	"fmt"
	"time"

	"github.com/TR-SLimey/E/esockets/mx"
)

func init() {
	// Create the esocket as a local variable
	var esocket = Esocket{
		ID: "matrix",
		onInit: func(es *Esocket, confLocation string) error {
			es.runlevel = 1
			/*err := confmgr.GetEsocketConfig(confLocation, &es.Config)
			if err != nil {
				return fmt.Errorf("Error reading esocket (%s) config: %s", es.ID, err.Error())
			}*/
			return nil
		},
		onDeinit: func(es *Esocket) error {
			es.runlevel = 0
			return nil
		},
		onStart: func(es *Esocket) error {
			es.runFlag = true
			go es.run(es)
			es.runlevel = 2
			return nil
		},
		onStop: func(es *Esocket) error {
			es.runFlag = false
			fmt.Println("Waiting for mainloop to exit")
			for es.runlevel != 1 {
				time.Sleep(5 * time.Microsecond)
			}
			return nil
		},
		run: func(es *Esocket) {
			for es.runFlag {
				time.Sleep(1 * time.Second)
				es.RecvQueue <- map[string]string{"recv": "works!"}
			}
			fmt.Println("Matrix Esocket Has Exit")
			es.runlevel = 1
		},
		Config: mx.Config{},
	}
	// Register the esocket so that it can be listed and used
	esocket.register()
}
