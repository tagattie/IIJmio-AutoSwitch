package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendMail(subject string) error {
	var auth smtp.Auth
	if config.Mail.Auth == true {
		auth = smtp.PlainAuth(
			"",
			config.Mail.Username,
			config.Mail.Password,
			config.Mail.SmtpServer)
	} else {
		auth = nil
	}

	body := constructMessage(subject)

	smtpServerAddr := config.Mail.SmtpServer + ":" + config.Mail.SmtpPort
	if err := smtp.SendMail(
		smtpServerAddr,
		auth,
		config.Mail.FromAddr,
		config.Mail.ToAddrs,
		([]byte)(body)); err != nil {
		fmt.Printf("Sending mail error: ", err)
		return err
	}
	return nil
}

func constructMessage(subject string) string {
	var toStr string
	for _, v := range config.Mail.ToAddrs {
		toStr += "To: " + v + "\r\n"
	}
	subjectStr := "Subject: IIJmio AutoSwitch: " + subject + "\r\n\r\n"

	var message string
	switch subject {
	case "Your application is not registered":
		message = "The configured developerId seems wrong.\n"
		message += "Please check your configuration.\n"
	case "User Authorization Failure":
		message = "Access token seems to have expired.\n"
		message += "Please acquire new access token at the following URL.\n\n"
		message += authUrl + "\n"
	default:
		msgStr := "An error occurred.\n"
		returnStr := fmt.Sprintf("Return code is %s.\n", subject)
		msgs := []string{msgStr, returnStr}
		message = strings.Join(msgs, "")
	}
	body := message + "\r\n"

	return toStr + subjectStr + body
}
