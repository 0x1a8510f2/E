package esockets

func init() {
	var esocket = Esocket{
		ID: "defaultHttp",
		BindPorts: []int{8080},
	}
	esocket.register()
}