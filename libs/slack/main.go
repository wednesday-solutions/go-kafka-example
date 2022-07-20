package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

type Alert struct {
	Token       string
	ChannelId   string
	MessageText string
}

type Slack interface {
	SendMessage() error
}

func (alert *Alert) SendMessage() error {
	api := slack.New(alert.Token)
	// slack.OptionDebug(true)
	message, _, _, err := api.SendMessage(alert.ChannelId,
		slack.MsgOptionText(alert.MessageText, false))
	fmt.Println("\nSlack::SendMessage::", message)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
