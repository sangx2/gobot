package model

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	token   string
	chatID  int64
	Channel string

	postChanChan chan chan *Post
	done         chan int

	botAPI         *tgbotapi.BotAPI
	config         tgbotapi.UpdateConfig
	updatesChannel tgbotapi.UpdatesChannel
}

func NewTelegram(token string, chatID int64, channel string) *Telegram {
	return &Telegram{token: token, chatID: chatID, Channel: channel, postChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

func (t *Telegram) IsValid() error {
	if len(t.token) == 0 {
		return errors.New("token length is 0")
	}
	if len(t.Channel) == 0 {
		return errors.New("channel length is 0")
	}
	return nil
}

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

func (t *Telegram) GetPostChanChan() chan chan *Post {
	return t.postChanChan
}

func (t *Telegram) Start() {
	go func() {
		postChan := <-t.postChanChan
		for {
			select {
			case update := <-t.updatesChannel:
				if t.chatID == 0 {
					t.chatID = update.Message.Chat.ID
				}
				postChan <- NewPost(MESSENGER_TELEGRAM, t.Channel, update.Message.Text)
			case <-t.done:
				break
			}
		}
	}()
}

func (t Telegram) Send(post *Post) error {
	switch post.Messenger {
	case MESSENGER_TELEGRAM:
		if strings.Compare(t.Channel, post.Channel) != 0 {
			return nil
		}
	default:
		return nil
	}

	if t.chatID == 0 {
		return errors.New("chatID is 0. if you want to send a post, craete model.NewTelegram with chatID or write a message once in telegram")
	} else {
		fmt.Println("chatID:", t.chatID)
	}

	telegramPost := tgbotapi.NewMessage(t.chatID, post.Message)

	_, e := t.botAPI.Send(telegramPost)
	if e != nil {
		return e
	}

	return nil
}

func (t Telegram) Logout() {
	t.botAPI.StopReceivingUpdates()
}

func (t Telegram) Shutdown() {
	t.done <- 1
}
