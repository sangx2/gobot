package model

type Bot interface {
	Login() error
	GetRecvPostChanChan() chan chan *Post
	Start()
	Send(*Post) error
	Logout()
	Shutdown()
}
