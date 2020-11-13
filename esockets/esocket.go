package esockets

/* List of all available esockets. Appended to by calling
esocket.register() */
var Available = make(map[string]*Esocket)

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

	/* This function should run anything necessary to set
	up the esocket, and exit quickly as it runs
	synchronously. */
	onInit func(es *Esocket)

	/* This function should run anything necessary to clean
	up after the esocket (including saving any data)
	and exit quickly as it runs synchronously. */
	onDeinit func(es *Esocket)

	/* This function should run in the background and handle
	all incoming and outgoing data. */
	onStart func(es *Esocket)

	/* This function should cleanly stop the onStart function. */
	onStop func(es *Esocket)

	/* Configuration supported or required by this esocket. */
	Config struct{}
}

/* Called when the esocket should prepare for receiving data.
Most likely on E's startup. */
func (es *Esocket) Init() {
	es.onInit(es)
}

/* Called when the esocket should clean up and exit. This
most likely means that E is exiting. */
func (es *Esocket) Deinit() {
	es.onDeinit(es)
}

/* Called when the esocket should start receiving
and outputting data. Runs asynchronously. */
func (es *Esocket) Start() {
	es.onStart(es)
}

/* Called when the esocket should stop receiving and
outputting data. However, the esocket should still hold on
to its data because it may be started again. */
func (es *Esocket) Stop() {
	es.onStop(es)
}

/* Register the esocket in the `Available` map to allow it
to be listed and used. This should be called just after an
esocket is initialised in its file. */
func (es *Esocket) register() {
	Available[es.ID] = es
}
