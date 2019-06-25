# Gobot
This is bot client package using mattermost, telegram

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
    "time"
    "fmt"

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
	}

	gobot := gobot.NewGobot(telegram, 100)

	e = gobot.Start()
	if e != nil {
		fmt.Printf("gobot.Start error : %s", e)
    }
    defer gobot.Shutdown()

	post := model.NewPost(model.MESSENGER_TELEGRAM, channel, "start TestGobot4Telegram")
	e = gobot.SendPost(post)
	if e != nil {
		fmt.Printf("gobot.SendPost error : %s", e)
	}

	done := make(chan int, 1)
	recvPostChan := gobot.GetRecvPostChan()

    // echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				post := model.NewPost(model.MESSENGER_TELEGRAM, recvPost.Channel, "[echo:"+recvPost.Channel+"] "+recvPost.Message)
				e = gobot.SendPost(post)
				if e != nil {
					fmt.Printf("gobot.SendPost error : %s", e)
				}
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()

    time.Sleep(time.Second * 5)
    
    done <- 1

	post = model.NewPost(model.MESSENGER_TELEGRAM, channel, "end TestGobot4Telegram")
	e = gobot.SendPost(post)
	if e != nil {
		fmt.Printf("gobot.SendPost error : %s", e)
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

	gobot := gobot.NewGobot(mattermost, 100)

	e = gobot.Start()
	if e != nil {
        fmt.Printf("gobot.Start error : %s", e)
        return
    }
    defer gobot.Shutdown()

	post := model.NewPost(model.MESSENGER_MATTERMOST, channel, "start TestGobot4Mattermost")
	e = gobot.SendPost(post)
	if e != nil {
        fmt.Printf("gobot.SendPost error : %s", e)
        return
	}

	done := make(chan int, 1)
	recvPostChan := gobot.GetRecvPostChan()

    // echo
	go func() {
		for {
			select {
			case recvPost := <-recvPostChan:
				post := model.NewPost(model.MESSENGER_MATTERMOST, recvPost.Channel, "[echo] "+recvPost.Message)
				e = gobot.SendPost(post)
				if e != nil {
                    fmt.Printf("gobot.SendPost error : %s", e)
                    return
				}
			case <-done:
				fmt.Println("done")
				return
			}
		}
    }()
    
    time.Sleep(time.Second * 5)

	done <- 1

	post = model.NewPost(model.MESSENGER_MATTERMOST, channel, "end TestGobot4Mattermost")
	e = gobot.SendPost(post)
	if e != nil {
		fmt.Printf("gobot.SendPost error : %s", e)
	}
}
```
