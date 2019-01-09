package messenger

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

// MATTERMOST :
const MATTERMOST = "mattermost"

// Mattermost :
type Mattermost struct {
	url      string
	username string
	password string
	team     string
	channel  string

	postChanChan chan chan *Post
	done         chan int

	client          *model.Client4
	botUser         *model.User
	botChannel      *model.Channel
	webSocketClient *model.WebSocketClient
}

// NewMattermost :
func NewMattermost(url string, username string, password string, team string, channel string) *Mattermost {
	return &Mattermost{url: url, username: username, password: password, team: team, channel: channel, postChanChan: make(chan chan *Post, 1), done: make(chan int, 1)}
}

// Login :
func (m *Mattermost) Login() error {
	m.client = model.NewAPIv4Client(m.url)

	if _, resp := m.client.GetOldClientConfig(""); resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	user, resp := m.client.Login(m.username, m.password)
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botUser = user

	team, resp := m.client.GetTeamByName(m.team, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	channel, resp := m.client.GetChannelByName(m.channel, team.Id, "")
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	m.botChannel = channel

	u, _ := url.Parse(m.url)

	webSocketClient, e := model.NewWebSocketClient4("wss://"+u.Hostname(), m.client.AuthToken)
	if e != nil {
		return errors.New(e.Message)
	}
	m.webSocketClient = webSocketClient

	m.webSocketClient.Listen()

	return nil
}

// GetPostChanChan :
func (m *Mattermost) GetPostChanChan() chan chan *Post {
	return m.postChanChan
}

// Start :
func (m *Mattermost) Start() {
	go func() {
		postChan := <-m.postChanChan
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
					if req.UserId == m.botUser.Id {
						continue
					}
				}

				postChan <- NewPost(MATTERMOST, m.channel, req.Message)
			case <-m.done:
				break
			}
		}
	}()
}

// Send :
func (m Mattermost) Send(post *Post) error {
	// check messenger & channel
	switch post.Messenger {
	case MATTERMOST:
		if strings.Compare(m.channel, post.Channel) != 0 {
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

// Logout :
func (m Mattermost) Logout() {
	m.client.Logout()
}

// Shutdown :
func (m Mattermost) Shutdown() {
	m.done <- 1
}
