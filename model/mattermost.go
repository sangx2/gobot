package model

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

// Mattermost mattermost 봇 구조체
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

// NewMattermost mattermost 봇 생성
func NewMattermost(url string, username string, password string, team string, channel string) *Mattermost {
	return &Mattermost{URL: url, username: username, password: password, Team: team, Channel: channel, recvPostChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

// IsValid mattermost 객체의 유효성 검사
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

// Login mattermost 로그인
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

// GetRecvPostChanChan 메시지를 전달할 채널를 위한 chan chan
func (m *Mattermost) GetRecvPostChanChan() chan chan *Post {
	return m.recvPostChanChan
}

// Start mattermost 봇 시작
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

				postChan <- NewPost(m.Channel, req.Message, nil)
			case <-m.done:
				break
			}
		}
	}()
}

// SendMessage mattermost 사용자에게 메시지 전달
func (m Mattermost) SendMessage(message string) error {
	mattermostPost := &model.Post{}
	mattermostPost.ChannelId = m.botChannel.Id
	mattermostPost.Message = message

	if _, resp := m.client.CreatePost(mattermostPost); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

// Logout mattermost 로그아웃
func (m Mattermost) Logout() {
	m.client.Logout()
}

// Shutdown mattermost 봇 종료
func (m Mattermost) Shutdown() {
	m.done <- 1
}
