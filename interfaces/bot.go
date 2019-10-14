package interfaces

import "github.com/sangx2/gobot/model"

// Bot 봇 인터페이스
type Bot interface {
	Login() error
	Start()
	GetRecvPostChanChan() chan chan *model.Post
	SendPost(*model.Post) error
	Logout()
	Shutdown()
}
