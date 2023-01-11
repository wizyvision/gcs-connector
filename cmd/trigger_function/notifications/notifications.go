package notifications

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type Message struct {
	Pretext string
	Text    string
}

func SendSlackMessage(message slack.MsgOption) {
	// Set slack env config
	SLACK_BEARER_TOKEN := os.Getenv("SLACK_BEARER_TOKEN")
	SLACK_CHANNEL_ID := os.Getenv("SLACK_CHANNEL_ID")

	if SLACK_BEARER_TOKEN == "" || SLACK_CHANNEL_ID == "" {
		return
	}

	slackClient := slack.New(SLACK_BEARER_TOKEN)

	// Executes sending message to slack
	_, _, err := slackClient.PostMessage(SLACK_CHANNEL_ID, message)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func sendSlackMsgDefault(message Message, status string) {
	SendSlackMessage(slack.MsgOptionAttachments(
		slack.Attachment{
			Color:   StatusColor(status),
			Pretext: message.Pretext,
			Text:    message.Text,
		},
	))
}

func SendSlackErrorMessage(message Message) {
	sendSlackMsgDefault(message, "error")
}

func SendSlackSuccessMessage(message Message) {
	sendSlackMsgDefault(message, "success")
}
