package main

import (
	"fmt"
	"strings"

	"github.com/bluele/slack"
)

func sendSlack(message string) error {
	api := slack.New(config.Slack.Token)
	var opt slack.ChatPostMessageOpt
	opt.Username = "IIJmio Autoswitch"
	if err := api.ChatPostMessage(config.Slack.Channel, constructSlackMessage(message), &opt); err != nil {
		fmt.Println("Slack message send error: ", err)
		return err
	}

	return nil
}

func constructSlackMessage(subject string) string {
	var message string
	switch subject {
	case "Your application is not registered":
		message = "The configured developerId seems wrong.\n"
		message += "Please check your configuration.\n"
	case "User Authorization Failure":
		message = "Access token seems to have expired.\n"
		message += "Please acquire new access token at the following URL.\n\n"
		message += "<" + authUrlEncoded + ">" + "\n"
	default:
		msgStr := "An error occurred.\n"
		returnStr := fmt.Sprintf("Return code is %s.\n", subject)
		msgs := []string{msgStr, returnStr}
		message = strings.Join(msgs, "")
	}
	return message
}
