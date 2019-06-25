package gobot

import (
	"github.com/sangx2/gobot/model"
)

type Gobot struct {
	bot model.Bot

	recvPostChan chan *model.Post
}

func NewGobot(bot model.Bot, recvPostChanSize int) *Gobot {
	return &Gobot{bot: bot, recvPostChan: make(chan *model.Post, recvPostChanSize)}
}

func (g *Gobot) Start() error {
	e := g.bot.Login()
	if e != nil {
		return e
	}

	g.bot.Start()

	g.bot.GetRecvPostChanChan() <- g.recvPostChan

	return nil
}

func (g *Gobot) GetRecvPostChan() chan *model.Post {
	return g.recvPostChan
}

func (g *Gobot) SendPost(post *model.Post) error {
	e := g.bot.Send(post)
	if e != nil {
		return e
	}

	return nil
}

func (g *Gobot) Shutdown() {
	g.bot.Logout()
	g.bot.Shutdown()
}
