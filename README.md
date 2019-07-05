# Gobot
This is bot package using mattermost, telegram

### install

```bash
go get -u github.com/sangx2/gobot
```

## Getting started

### Gobot for telegram
chatID is not necessary value, but you can't send a post.

If you want to send a post, craete model.NewTelegram with chatID or write a message once in telegram

```go
package main

import (
	"fmt"
	"time"

	"github.com/sangx2/gobot"
	"github.com/sangx2/gobot/model"
)

func main() {
	// fix me
	token := ""
	chatID := 0
	channel := ""

	telegram := model.NewTelegram(token, int64(chatID), channel)
	e := telegram.IsValid()
	if e != nil {
		fmt.Printf("NewTelegram error : %s", e)
		return
	}

	recvPostChan := make(chan *model.Post, 100)
	gobot := gobot.NewGobot(telegram, recvPostChan)

	e = gobot.Login()
	if e != nil {
		fmt.Printlnf("gobot.Login error : %s", e)
		return
	}
	defer gobot.Logout()

	gobot.Start()
	defer gobot.Shutdown()

	e = gobot.SendMessage("start TestGobot4Telegram")
	if e != nil {
		fmt.Printf("gobot.SendMessage error : %s", e)
		return
	}

	done := make(chan int, 1)

	// echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendMessage("[echo:" + recvPost.Channel + "] " + recvPost.Message)
				if e != nil {
					fmt.Printf("gobot.SendMessage error : %s", e)
				}
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()

	time.Sleep(time.Second * 5)

	done <- 1

	e = gobot.SendMessage("end TestGobot4Telegram")
	if e != nil {
		fmt.Printf("gobot.SendMessage error : %s", e)
	}
}
```

### Gobot for mattermost

```go
package main

import (
    "fmt"
    "time"

    "github.com/sangx2/gobot"
	"github.com/sangx2/gobot/model"
)

func main() {
	// fix me
	url := ""
	username := ""
	password := ""
	team := ""
	channel := ""

	mattermost := model.NewMattermost(url, username, password, team, channel)
	e := mattermost.IsValid()
	if e != nil {
        fmt.Printf("NewMattermost error : %s", e)
        return
	}

	recvPostChan := make(chan *model.Post, 100)
	gobot := gobot.NewGobot(mattermost, recvPostChan)

	e = gobot.Login()
	if e != nil {
		fmt.Printf("gobot.Login error : %s", e)
		return
	}
	defer gobot.Logout()

	gobot.Start()
    defer gobot.Shutdown()

	e = gobot.SendMessage("start TestGobot4Mattermost")
	if e != nil {
        fmt.Printf("gobot.SendMessage error : %s", e)
        return
	}

	done := make(chan int, 1)

    // echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				e = gobot.SendMessage("[echo] " + recvPost.Message)
				if e != nil {
                    fmt.Printf("gobot.SendMessage error : %s", e)
				}
			case <-done:
				fmt.Println("done")
				return
			}
		}
    }()
    
    time.Sleep(time.Second * 5)

	done <- 1

	e = gobot.SendMessage("end TestGobot4Mattermost")
	if e != nil {
		fmt.Printf("gobot.SendMessage error : %s", e)
	}
}
```