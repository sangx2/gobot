package model

type Bot interface {
	Login() error
	GetPostChanChan() chan chan *Post
	Start()
	Send(*Post) error
	Logout()
	Shutdown()
}
