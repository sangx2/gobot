package messenger

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// TELEGRAM :
const TELEGRAM = "telegram"

// Telegram :
type Telegram struct {
	token   string
	chatID  int64
	channel string

	postChanChan chan chan *Post
	done         chan int

	botAPI        *tgbotapi.BotAPI
	config        tgbotapi.UpdateConfig
	updateChannel tgbotapi.UpdatesChannel
}

// NewTelegram :
func NewTelegram(token string, chatID int64, channel string) *Telegram {
	return &Telegram{token: token, chatID: chatID, channel: channel, postChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

// Login :
func (t *Telegram) Login() error {
	t.config = tgbotapi.NewUpdate(0)
	t.config.Timeout = 60

	bot, e := tgbotapi.NewBotAPI(t.token)
	if e != nil {
		return e
	}
	t.botAPI = bot

	t.updateChannel, e = t.botAPI.GetUpdatesChan(t.config)
	if e != nil {
		return e
	}

	return nil
}

// GetPostChanChan :
func (t *Telegram) GetPostChanChan() chan chan *Post {
	return t.postChanChan
}

// Start :
func (t *Telegram) Start() {
	go func() {
		postChan := <-t.postChanChan
		for {
			select {
			case req := <-t.updateChannel:
				postChan <- NewPost(TELEGRAM, t.channel, req.Message.Text)
			case <-t.done:
				break
			}
		}
	}()
}

// Send :
func (t Telegram) Send(post *Post) error {
	// check messenger & channel
	switch post.Messenger {
	case TELEGRAM:
		if strings.Compare(t.channel, post.Channel) != 0 {
			return nil
		}
	default:
		return nil
	}

	telegramPost := tgbotapi.NewMessage(t.chatID, post.Message)

	_, e := t.botAPI.Send(telegramPost)
	if e != nil {
		return e
	}

	return nil
}

// Logout :
func (t Telegram) Logout() {
	t.botAPI.StopReceivingUpdates()
}

// Shutdown :
func (t Telegram) Shutdown() {
	t.done <- 1
}
