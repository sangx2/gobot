package gobot

import (
	"errors"

	"github.com/sangx2/gobot/model"
	tgbotapi "gopkg.in/telegram-bot-api"
)

// Telegram telegram 봇
type Telegram struct {
	token   string
	chatID  int64
	Channel string

	recvPostChanChan chan chan *model.Post
	done             chan int

	botAPI         *tgbotapi.BotAPI
	config         tgbotapi.UpdateConfig
	updatesChannel tgbotapi.UpdatesChannel
}

// NewTelegram telegram 봇 생성
func NewTelegram(token string, channel string) *Telegram {
	return &Telegram{token: token, chatID: 0, Channel: channel,
		recvPostChanChan: make(chan chan *model.Post, 1), done: make(chan int, 1)}
}

// IsValid 객체의 유효성 검사
func (t *Telegram) IsValid() error {
	if len(t.token) == 0 {
		return errors.New("token length is 0")
	}
	if len(t.Channel) == 0 {
		return errors.New("channel length is 0")
	}
	return nil
}

// Login telegram 봇 로그인
func (t *Telegram) Login() error {
	t.config = tgbotapi.NewUpdate(0)
	t.config.Timeout = 60

	bot, e := tgbotapi.NewBotAPI(t.token)
	if e != nil {
		return e
	}
	t.botAPI = bot

	t.updatesChannel, e = t.botAPI.GetUpdatesChan(t.config)
	if e != nil {
		return e
	}

	t.updatesChannel.Clear()

	return nil
}

// GetRecvPostChanChan 메시지를 전달할 채널를 위한 chan chan
func (t *Telegram) GetRecvPostChanChan() chan chan *model.Post {
	return t.recvPostChanChan
}

// Start telegram 봇 시작
func (t *Telegram) Start() {
	go func() {
		postChan := <-t.recvPostChanChan
		for {
			select {
			case update := <-t.updatesChannel:
				if update.Message == nil {
					continue
				}

				if t.chatID == 0 {
					t.chatID = update.Message.Chat.ID
				}

				postChan <- model.NewPost(update.Message.Text, update.Message.MessageID)
			case <-t.done:
				break
			}
		}
	}()
}

// SendPost telegram 사용자에게 메시지 전달
func (t Telegram) SendPost(post *model.Post) error {
	if t.chatID == 0 {
		return errors.New("chatID is 0. if you want to send a post," +
			" craete model.NewTelegram with chatID or write a message once in telegram")
	}

	telegramPost := tgbotapi.NewMessage(t.chatID, post.Message)
	if rootID, ok := post.RootID.(int); ok {
		telegramPost.ReplyToMessageID = rootID
	}

	_, e := t.botAPI.Send(telegramPost)
	if e != nil {
		return e
	}

	return nil
}

// Logout telegram 봇 로그아웃
func (t Telegram) Logout() {
	t.botAPI.StopReceivingUpdates()
}

// Shutdown telegram 봇 종료
func (t Telegram) Shutdown() {
	t.done <- 1
}
