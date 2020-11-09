package esockets

func init() {
	var esocket = Esocket{
		ID: "defaultHttp",
		BindPort: 1234,
	}
	Available = append(Available, esocket)
}