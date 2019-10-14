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

// Login 메신저 로그인
func (g *Gobot) Login() error {
	e := g.bot.Login()
	if e != nil {
		return e
	}
	return nil
}

// Logout 메신저 로그아웃
func (g *Gobot) Logout() {
	g.bot.Logout()
}

// Start 봇 서버 시작
func (g *Gobot) Start() {
	g.bot.Start()

	g.bot.GetRecvPostChanChan() <- g.recvPostChan
}

// SendMessage 사용자에게 메시지 전달
func (g *Gobot) SendMessage(message string, param interface{}) error {
	e := g.bot.SendMessage(message, param)
	if e != nil {
		return e
	}
	return nil
}

// Shutdown 봇서버 종료
func (g *Gobot) Shutdown() {
	g.bot.Shutdown()
}
