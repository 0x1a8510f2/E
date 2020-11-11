package esockets

/* List of all available esockets. Appended to by calling
esocket.register() */
var Available = make(map[string]Esocket)

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

	/* Configuration supported or required by this esocket.
	TODO - Find appropriate structure for this*/
	Config int
}

/* Register the esocket in the `Available` map to allow it
to be listed and used. This should be called after an
esocket is defined */
func (es Esocket) register() {
	Available[es.ID] = es
}
