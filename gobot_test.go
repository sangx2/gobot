package gobot

import (
	"testing"
	"time"

	"github.com/sangx2/gobot/model"
)

const (
	POST_CHAN_SIZE = 100
)

func TestGobot4Mattermost(t *testing.T) {
	// fix me
	url := ""
	username := ""
	password := ""
	team := ""
	channel := ""

	mattermost := NewMattermost(url, username, password, team, channel)
	e := mattermost.IsValid()
	if e != nil {
		t.Errorf("NewMattermost error : %s", e)
	}

	recvPostChan := make(chan *model.Post, POST_CHAN_SIZE)
	gobot := NewGobot(mattermost, recvPostChan)

	e = gobot.Login()
	if e != nil {
		t.Errorf("gobot.Login error : %s", e)
	}
	defer gobot.Logout()

	gobot.Start()
	defer gobot.Shutdown()

	e = gobot.SendPost(model.NewPost("start TestGobot4Mattermost", nil))
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendPost(model.NewPost("[echo] "+recvPost.Message, recvPost.RootID))
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

	e = gobot.SendPost(model.NewPost("end TestGobot4Mattermost", nil))
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}
}

func TestGobot4Telegram(t *testing.T) {
	// fix me
	token := ""
	channel := ""

	telegram := NewTelegram(token, channel)
	e := telegram.IsValid()
	if e != nil {
		t.Errorf("NewTelegram error : %s", e)
	}

	recvPostChan := make(chan *model.Post, POST_CHAN_SIZE)
	gobot := NewGobot(telegram, recvPostChan)

	e = gobot.Login()
	if e != nil {
		t.Errorf("gobot.Login error : %s", e)
	}
	defer gobot.Logout()

	gobot.Start()
	defer gobot.Shutdown()

	e = gobot.SendPost(model.NewPost("start TestGobot4Telegram", nil))
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendPost(model.NewPost("[echo] "+recvPost.Message, recvPost.RootID))
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

	e = gobot.SendPost(model.NewPost("end TestGobot4Telegram", nil))
	if e != nil {
		t.Errorf("gobot.SendPost error : %s", e)
	}
}
