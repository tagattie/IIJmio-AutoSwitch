package main

import (
	"fmt"
	"net/smtp"
)

func sendMail(message string) error {
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

	smtpServerAddr := config.Mail.SmtpServer + ":" + config.Mail.SmtpPort
	if err := smtp.SendMail(
		smtpServerAddr,
		auth,
		config.Mail.FromAddr,
		config.Mail.ToAddrs,
		([]byte)(message)); err != nil {
			fmt.Printf("Sending mail error: ", err)
			return err
		}
	return nil
}
