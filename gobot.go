package gobot

import (
	"github.com/sangx2/gobot/interfaces"
	"github.com/sangx2/gobot/model"
)

// Gobot 봇 구조체
type Gobot struct {
	bot interfaces.Bot

	recvPostChan chan *model.Post
}

// NewGobot 봇 구조체 생성
func NewGobot(bot interfaces.Bot, recvPostChan chan *model.Post) *Gobot {
	return &Gobot{bot: bot, recvPostChan: recvPostChan}
}

// Login 봇 로그인
func (g *Gobot) Login() error {
	e := g.bot.Login()
	if e != nil {
		return e
	}
	return nil
}

// Logout 봇 로그아웃
func (g *Gobot) Logout() {
	g.bot.Logout()
}

// Start 봇 시작
func (g *Gobot) Start() {
	g.bot.Start()

	g.bot.GetRecvPostChanChan() <- g.recvPostChan
}

// SendPost 사용자에게 메시지 전달
func (g *Gobot) SendPost(post *model.Post) error {
	e := g.bot.SendPost(post)
	if e != nil {
		return e
	}
	return nil
}

// Shutdown 봇 종료
func (g *Gobot) Shutdown() {
	g.bot.Shutdown()
}
