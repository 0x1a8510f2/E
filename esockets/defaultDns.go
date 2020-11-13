package esockets

func init() {
	var esocket = Esocket{
		ID: "defaultDns",
	}
	esocket.onInit = func(es *Esocket) {
		println(es.ID)
	}
	// Register the esocket so that it can be used by E
	esocket.register()
}
