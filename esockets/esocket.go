package esockets

// List (slice) of all Esockets available in the module
var Available []Esocket

type Esocket struct {
	// Human readable ID used to refer to this esocket
	ID string
	// Array of ports which this esocket listens on
	BindPorts []int
}

/* Register the esocket in the Available slice to allow it
to be listed and used. This should be executed when an
esocket is defined */
func (es Esocket) register() {
	Available = append(Available, es)
}