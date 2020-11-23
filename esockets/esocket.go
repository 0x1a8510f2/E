package esockets

import (
	"fmt"

	sr "github.com/TR-SLimey/E/stringres"
)

/* List of all available esockets. Appended to by calling
esocket.register(). */
var Available = make(map[string]*Esocket)

/* What runlevels translate to in human readable format */
var Runlevels = [...]string{
	"UNINITIALISED",
	"INITIALISED",
	"RUNNING",
}

/* The data structure of an Esocket. Contains everything
needed to communicate with a client in the protocol it
represents, and translate between that protocol and
Go data structures (to be transalted to Matrix and vice
versa) */
type Esocket struct {
	/* Human readable ID used to refer to this esocket.
	Can be any string but must be unique between esocket
	instances, else the last instance to register the
	ID will replace all others. The name of the file
	defining the esocket is a good choice for an ID */
	ID string

	/* The current runlevel of the esocket as an integer where:
		0 => Not initialised (or has been deinitialised)
		1 => Initialised but not running (or stopped after running)
		2 => Active/running
	This should be updated by the esocket like so:
		0 => The esocket does not expect to be deinitialised before exit
		1 => The esocket does not expect to be stopped, but should be deinitialised before exit
		2 => The esocket should be both stopped and deinitialised before exit
	*/
	Runlevel int

	/* This function should run anything necessary to set
	up the esocket, and exit quickly as it runs
	synchronously. On top of the esocket object, it should
	also accept the location of the config file which it
	should process. */
	onInit func(es *Esocket, confLocation string) error

	/* This function should run anything necessary to clean
	up after the esocket (including saving any data)
	and exit quickly as it runs synchronously. */
	onDeinit func(es *Esocket) error

	/* This function should run in the background and handle
	all incoming and outgoing data. */
	onStart func(es *Esocket) error

	/* This function should cleanly stop the onStart function. */
	onStop func(es *Esocket) error

	/* */
	Run func(es *Esocket)

	/* */
	runFlag bool

	/* */
	SendQueue chan map[string]string

	/* */
	RecvQueue chan map[string]string

	/* Configuration supported or required by this esocket. */
	Config struct{}
}

/* Called when the esocket should prepare for receiving data.
Most likely on E's startup. */
func (es *Esocket) Init(confLocation string) error {
	return es.onInit(es, confLocation)
}

/* Called when the esocket should clean up and exit. This
most likely means that E is exiting. */
func (es *Esocket) Deinit() error {
	return es.onDeinit(es)
}

/* Called when the esocket should start receiving
and outputting data. */
func (es *Esocket) Start() error {
	return es.onStart(es)
}

/* Called when the esocket should stop receiving and
outputting data. However, the esocket should still hold on
to its data because it may be started again. */
func (es *Esocket) Stop() error {
	return es.onStop(es)
}

/* Check whether the runlevel of the esocket is as expected */
func (es *Esocket) CheckRunlevel(expected int) error {
	// Ensure the expected runlevel is valid to prevent further errors
	if expected < 0 || expected > len(Runlevels)-1 {
		return fmt.Errorf(sr.INVALID_EXPECTED_RUNLEVEL)
	}
	// Actually check runlevel
	if es.Runlevel != expected {
		return fmt.Errorf(sr.UNEXPECTED_RUNLEVEL_ERR, Runlevels[es.Runlevel], Runlevels[expected])
	}
	return nil
}

/* Register the esocket in the `Available` map to allow it
to be listed and used. This should be called just after an
esocket is initialised in its file. */
func (es *Esocket) register() {
	Available[es.ID] = es
}
