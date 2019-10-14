package interfaces

import "github.com/sangx2/gobot/model"

// Bot 봇 인터페이스
type Bot interface {
	Login() error
	GetRecvPostChanChan() chan chan *model.Post
	Start()
	SendMessage(string, interface{}) error
	Logout()
	Shutdown()
}
