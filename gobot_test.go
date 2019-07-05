package gobot

import (
	"testing"
	"time"

	"github.com/sangx2/gobot/model"
)

func TestGobot4Mattermost(t *testing.T) {
	// fix me
	url := ""
	username := ""
	password := ""
	team := ""
	channel := ""

	mattermost := model.NewMattermost(url, username, password, team, channel)
	e := mattermost.IsValid()
	if e != nil {
		t.Errorf("NewMattermost error : %s", e)
	}

	recvPostChan := make(chan *model.Post, 100)
	gobot := NewGobot(mattermost, recvPostChan)

	e = gobot.Login()
	if e != nil {
		t.Errorf("gobot.Login error : %s", e)
	}
	defer gobot.Logout()

	gobot.Start()
	defer gobot.Shutdown()

	post := model.NewPost(model.MESSENGER_MATTERMOST, channel, "start TestGobot4Mattermost")
	e = gobot.SendPost(post)
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				post := model.NewPost(model.MESSENGER_MATTERMOST, recvPost.Channel, "[echo] "+recvPost.Message)
				e = gobot.SendPost(post)
				if e != nil {
					t.Errorf("gobot.SendPost error : %s", e)
				}
			case <-done:
				t.Log("done")
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)

	done <- 1

	post = model.NewPost(model.MESSENGER_MATTERMOST, channel, "end TestGobot4Mattermost")
	e = gobot.SendPost(post)
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}
}

func TestGobot4Telegram(t *testing.T) {
	// fix me
	token := ""
	chatID := 0
	channel := ""

	telegram := model.NewTelegram(token, int64(chatID), channel)
	e := telegram.IsValid()
	if e != nil {
		t.Errorf("NewTelegram error : %s", e)
	}

	recvPostChan := make(chan *model.Post, 100)
	gobot := NewGobot(telegram, recvPostChan)

	e = gobot.Login()
	if e != nil {
		t.Errorf("gobot.Login error : %s", e)
	}
	defer gobot.Logout()

	gobot.Start()
	defer gobot.Shutdown()

	post := model.NewPost(model.MESSENGER_TELEGRAM, channel, "start TestGobot4Telegram")
	e = gobot.SendPost(post)
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				post := model.NewPost(model.MESSENGER_TELEGRAM, recvPost.Channel, "[echo:"+recvPost.Channel+"] "+recvPost.Message)
				e = gobot.SendPost(post)
				if e != nil {
					t.Errorf("gobot.SendPost error : %s", e)
				}
			case <-done:
				t.Log("done")
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)

	done <- 1

	post = model.NewPost(model.MESSENGER_TELEGRAM, channel, "end TestGobot4Telegram")
	e = gobot.SendPost(post)
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}
}
