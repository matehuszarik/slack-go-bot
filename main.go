package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)

	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Printf("Message received: %v\n", msg)
			switch event := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connected %v\n", event.ConnectionCount)
			case *slack.MessageEvent:
				fmt.Printf("Message received %v\n", event)
				info := rtm.GetInfo()

				prefix := fmt.Sprintf("<@%s>", info.User.ID)

				if event.User != info.User.ID && strings.HasPrefix(event.Text, prefix) {
					user, _ := api.GetUserInfo(event.User)

					msg := fmt.Sprintf("Acknowledged %v", user.Name)

					rtm.SendMessage(rtm.NewOutgoingMessage(msg, event.Channel))
				}
			case *slack.RTMError:
				fmt.Printf("Error: %v\n", event)
			case *slack.InvalidAuthEvent:
				fmt.Print("INVALID CREDENTIALS")
			}
		}
	}
}
