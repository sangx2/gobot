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

	e = gobot.SendMessage("start TestGobot4Mattermost")
	if e != nil {
		t.Errorf("gobot.SendMessage error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendMessage("[echo] "+recvPost.Message, recvPost.Param)
				if e != nil {
					t.Errorf("gobot.SendMessage error : %s", e)
				}
			case <-done:
				t.Log("done")
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)

	done <- 1

	e = gobot.SendMessage("end TestGobot4Mattermost")
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

	e = gobot.SendMessage("start TestGobot4Telegram")
	if e != nil {
		t.Errorf("gobot.SendMessage error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendMessage("[echo:"+recvPost.Channel+"] "+recvPost.Message, recvPost.Param)
				if e != nil {
					t.Errorf("gobot.SendMessage error : %s", e)
				}
			case <-done:
				t.Log("done")
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)

	done <- 1

	e = gobot.SendMessage("end TestGobot4Telegram")
	if e != nil {
		t.Errorf("gobot.SendMessage error : %s", e)
	}
}
