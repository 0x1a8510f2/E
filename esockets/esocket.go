package esockets

// List (slice) of all Esockets available in the module
var Available []Esocket

type Esocket struct {
	ID string
	BindPort int
}