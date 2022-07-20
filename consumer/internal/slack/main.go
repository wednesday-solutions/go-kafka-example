package slack

import (
	"fmt"
	"libs/slack"
	"os"
)

func SendMessage(message string) {
	alert := slack.Alert{
		Token:     os.Getenv("SLACK_TOKEN"),
		ChannelId: os.Getenv("SLACK_CHANNEL_ID"),
	}
	alert.MessageText = message
	if err := alert.SendMessage(); err != nil {
		fmt.Printf("\n Failed To Send Slack Message, %s \n", err)
	}
}
