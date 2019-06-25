package model

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

type Mattermost struct {
	URL      string
	username string
	password string
	Team     string
	Channel  string

	recvPostChanChan chan chan *Post
	done             chan int

	client          *model.Client4
	botUser         *model.User
	botChannel      *model.Channel
	webSocketClient *model.WebSocketClient
}

func NewMattermost(url string, username string, password string, team string, channel string) *Mattermost {
	return &Mattermost{URL: url, username: username, password: password, Team: team, Channel: channel, recvPostChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

func (m *Mattermost) IsValid() error {
	if len(m.URL) == 0 {
		return errors.New("url is nil")
	}

	_, e := url.Parse(m.URL)
	if e != nil {
		return e
	}

	if len(m.username) == 0 {
		return errors.New("username is nil")
	}
	if len(m.password) == 0 {
		return errors.New("password is nil")
	}
	if len(m.Team) == 0 {
		return errors.New("team is nil")
	}
	if len(m.Channel) == 0 {
		return errors.New("channel is nil")
	}

	return nil
}

func (m *Mattermost) Login() error {
	m.client = model.NewAPIv4Client(m.URL)

	if _, resp := m.client.GetOldClientConfig(""); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	user, resp := m.client.Login(m.username, m.password)
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botUser = user

	team, resp := m.client.GetTeamByName(m.Team, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	channel, resp := m.client.GetChannelByName(m.Channel, team.Id, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botChannel = channel

	u, _ := url.Parse(m.URL)

	webSocketClient, e := model.NewWebSocketClient4("wss://"+u.Hostname(), m.client.AuthToken)
	if e != nil {
		return errors.New(e.Message)
	}
	m.webSocketClient = webSocketClient

	m.webSocketClient.Listen()

	return nil
}

func (m *Mattermost) GetRecvPostChanChan() chan chan *Post {
	return m.recvPostChanChan
}

func (m *Mattermost) Start() {
	go func() {
		postChan := <-m.recvPostChanChan
		for {
			select {
			case eventChannel := <-m.webSocketClient.EventChannel:
				if eventChannel.Broadcast.ChannelId != m.botChannel.Id {
					continue
				}
				if eventChannel.Event != model.WEBSOCKET_EVENT_POSTED {
					continue
				}
				req := model.PostFromJson(strings.NewReader(eventChannel.Data["post"].(string)))
				if req != nil {
					if len(req.PendingPostId) == 0 {
						continue
					}
				}

				postChan <- NewPost(MESSENGER_MATTERMOST, m.Channel, req.Message)
			case <-m.done:
				break
			}
		}
	}()
}

func (m Mattermost) Send(post *Post) error {
	switch post.Messenger {
	case MESSENGER_MATTERMOST:
		if strings.Compare(m.Channel, post.Channel) != 0 {
			return nil
		}
	default:
		return nil
	}

	mattermostPost := &model.Post{}
	mattermostPost.ChannelId = m.botChannel.Id
	mattermostPost.Message = post.Message

	if _, resp := m.client.CreatePost(mattermostPost); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

func (m Mattermost) Logout() {
	m.client.Logout()
}

func (m Mattermost) Shutdown() {
	m.done <- 1
}
