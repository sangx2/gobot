package gobot

import (
	"github.com/sangx2/gobot/interfaces"
	"github.com/sangx2/gobot/messenger"
)

// Gobot :
type Gobot struct {
	recvPostChan chan *messenger.Post
	recvDone     chan int

	SendPostChan chan *messenger.Post
	sendDone     chan int

	botList []interfaces.Bot

	logger interfaces.Logger
}

// NewGobot :
func NewGobot(config Config, logger interfaces.Logger) *Gobot {
	g := &Gobot{logger: logger}

	g.recvPostChan = make(chan *messenger.Post, config.RecvPostChanSize)
	g.recvDone = make(chan int, 1)

	g.SendPostChan = make(chan *messenger.Post, config.SendPostChanSize)
	g.sendDone = make(chan int, 1)

	for _, t := range config.TelegramBotSettings {
		if t.Enable {
			g.botList = append(g.botList, messenger.NewTelegram(t.Token, t.ChatID, t.Channel))
		}
	}

	for _, m := range config.MattermostBotSettings {
		if m.Enable {
			g.botList = append(g.botList, messenger.NewMattermost(m.URL, m.Username, m.Password, m.Team, m.Channel))
		}
	}

	return g
}

// StartGobot :
func (g *Gobot) StartGobot(f func(*messenger.Post) []*messenger.Post) error {
	for _, b := range g.botList {
		e := b.Login()
		if e != nil {
			return e
		}
		b.Start()

		b.GetPostChanChan() <- g.recvPostChan
	}

	go func() {
		for {
			select {
			case recvPost := <-g.recvPostChan:
				postList := f(recvPost)

				for _, post := range postList {
					g.SendPostChan <- post
				}
			case <-g.recvDone:
				break
			}
		}
	}()

	go func() {
		for {
			select {
			case post := <-g.SendPostChan:
				// broadcast
				for _, b := range g.botList {
					e := b.Send(post)
					if e != nil {
						g.logger.Error(e.Error())
					}
				}
			case <-g.sendDone:
				break
			}
		}
	}()

	return nil
}

// ShutdownGobot :
func (g *Gobot) ShutdownGobot() {
	for _, b := range g.botList {
		b.Shutdown()
	}

	g.recvDone <- 1
	g.sendDone <- 1
}
