package repository

type IHub interface {
	OnConnect(client IClient)
	OnDisconnect(client IClient)
	GetNumberOfActiveClients() int
	RegisterClient(client IClient)
	DeRegisterClient(client IClient)
	Run()
	BroadCast(message interface{}, ignoreClient string)
}

var AppHub IHub

func SetAppHub(hub IHub) {
	AppHub = hub
}

func Run() {
	go AppHub.Run()
}

func RegisterClient(client IClient) {
	AppHub.RegisterClient(client)
}

func DeRegisterClient(client IClient) {
	AppHub.DeRegisterClient(client)
}

func GetNumberOfActiveClients() int {
	return AppHub.GetNumberOfActiveClients()
}
