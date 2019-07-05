package model

type Bot interface {
	Login() error
	GetRecvPostChanChan() chan chan *Post
	Start()
	SendMessage(string) error
	Logout()
	Shutdown()
}
