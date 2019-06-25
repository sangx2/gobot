package gobot

import (
	"github.com/sangx2/gobot/model"
)

type Gobot struct {
	bot model.Bot

	recvPostChan chan *model.Post
}

func NewGobot(bot model.Bot, recvPostChan chan *model.Post) *Gobot {
	return &Gobot{bot: bot, recvPostChan: recvPostChan}
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
