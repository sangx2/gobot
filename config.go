package gobot

import (
	"errors"
	"fmt"
	"net/url"
)

// MattermostBotSettings :
type MattermostBotSettings struct {
	Enable   bool
	URL      string
	Username string
	Password string
	Team     string
	Channel  string
}

// TelegramBotSettings :
type TelegramBotSettings struct {
	Enable  bool
	Token   string
	ChatID  int64
	Channel string
}

// Config :
type Config struct {
	RecvPostChanSize int
	SendPostChanSize int

	MattermostBotSettings []MattermostBotSettings
	TelegramBotSettings   []TelegramBotSettings
}

// IsValid :
func (c *Config) IsValid() error {
	if c.RecvPostChanSize <= 0 {
		return errors.New("Config.RecvPostChanSize is not vaild")
	}

	if c.SendPostChanSize <= 0 {
		return errors.New("Config.SendPostChanSize is not vaild")
	}

	for i, m := range c.MattermostBotSettings {
		if e := m.isValid(); e != nil {
			return fmt.Errorf("%dth %s", i+1, e.Error())
		}
	}

	for i, t := range c.TelegramBotSettings {
		if e := t.isValid(); e != nil {
			return fmt.Errorf("%dth %s", i+1, e.Error())
		}
	}
	return nil
}

func (ms *MattermostBotSettings) isValid() error {
	if ms.Enable {
		if len(ms.URL) == 0 {
			return errors.New("Config.MattermostBotSettings.URL is nil")
		}

		_, e := url.Parse(ms.URL)
		if e != nil {
			return e
		}

		if len(ms.Username) == 0 {
			return errors.New("Config.MattermostBotSettings.Username is nil")
		}
		if len(ms.Password) == 0 {
			return errors.New("Config.MattermostBotSettings.Pssword is nil")
		}
		if len(ms.Team) == 0 {
			return errors.New("Config.MattermostBotSettings.Team is nil")
		}
		if len(ms.Channel) == 0 {
			return errors.New("Config.MattermostBotSettings.Channel is nil")
		}
	}
	return nil
}

func (ts *TelegramBotSettings) isValid() error {
	if ts.Enable {
		if len(ts.Token) == 0 {
			return errors.New("Config.TelegramBotSettings.Token is nil")
		}
		if ts.ChatID == 0 {
			return errors.New("Config.TelegramBotSettings.ChatID is 0")
		}
	}
	return nil
}
