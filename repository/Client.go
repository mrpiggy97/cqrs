package repository

type IClient interface {
	Write()
	GetId() string
	CloseConnection()
	Send(data []byte)
}
