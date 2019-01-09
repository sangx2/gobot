package interfaces

import "github.com/sangx2/gobot/messenger"

// Bot :
type Bot interface {
	Login() error
	GetPostChanChan() chan chan *messenger.Post
	Start()
	Send(*messenger.Post) error
	Logout()
	Shutdown()
}
