package model

const (
	MESSENGER_TELEGRAM   = "telegram"
	MESSENGER_MATTERMOST = "mattermost"
)

type Post struct {
	Messenger string
	Channel   string
	Message   string
}

func NewPost(messenger string, channel string, message string) *Post {
	return &Post{Messenger: messenger, Channel: channel, Message: message}
}
