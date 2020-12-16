package esockets

import (
	"fmt"
	"time"

	sr "github.com/TR-SLimey/E/stringres"
)

/* List of all available esockets. Appended to by calling
esocket.register(). */
var Available = make(map[string]*Esocket)

/* What error codes returned by E mean */
var ReturnErrCodes = [...]string{
	"generic error",               // 0
	"event cancelled",             // 1
	"unable to route event",       // 2
	"malformed event",             // 3
	"destination reported error",  // 4
	"client already registered",   // 5
	"invalid data type for event", // 6
}

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
	runlevel int

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
	run func(es *Esocket)

	/* */
	runFlag bool

	/* */
	InQueue chan map[string]string

	/* */
	OutQueue chan map[string]string

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

/* */
func (es *Esocket) readInQueue(timeout time.Duration) (map[string]string, error) {
	if timeout >= 0 {
		select {
		// Timeout (usually used to not block Esocket mainloop indefinitely)
		case <-time.After(timeout * time.Millisecond):
			return nil, nil

		case data := <-es.InQueue:
			return data, nil
		}
	} else {
		// If the timeout is negative, don't time out
		data := <-es.InQueue
		return data, nil
	}
}

/* */
func (es *Esocket) writeOutQueue(data map[string]string, timeout time.Duration) error {
	if timeout >= 0 {
		select {
		// Timeout (writing should never really time out unless something is
		// wrong with the E mainloop so this can be used to catch such occurences)
		case <-time.After(timeout * time.Millisecond):
			return fmt.Errorf(sr.ESOCKET_OUT_QUEUE_WRITE_TIMEOUT, es.ID)

		case es.OutQueue <- data:
			return nil
		}
	} else {
		// If the timeout is negative, don't time out
		es.OutQueue <- data
		return nil
	}
}

/* Check whether the runlevel of the esocket is as expected */
func (es *Esocket) CheckRunlevel(expected int) error {
	// Ensure the expected runlevel is valid to prevent further errors
	if expected < 0 || expected > len(Runlevels)-1 {
		return fmt.Errorf(sr.INVALID_EXPECTED_RUNLEVEL, expected)
	}
	// Actually check runlevel
	if es.runlevel != expected {
		return fmt.Errorf(sr.UNEXPECTED_RUNLEVEL_ERR, Runlevels[es.runlevel], Runlevels[expected])
	}
	return nil
}

/* Register the esocket in the `Available` map to allow it
to be listed and used. This should be called just after an
esocket is initialised in its file. */
func (es *Esocket) register() {
	Available[es.ID] = es
}
