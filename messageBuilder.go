package main

import (
	"fmt"
	"strings"
)

func buildToHeader() string {
	var toStr string
	for _, v := range config.Mail.ToAddrs {
		toStr += "To: " + v + "\r\n"
	}
	return toStr
}

func buildErrorMessage(subject string) string {
	toStr := buildToHeader()
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

func buildReportMessage(latestPacketData map[string][]int,
	couponState map[string]bool,
	couponAmount int,
	couponReqInfo map[string]bool) string {
	toStr := buildToHeader()
	subjectStr := "Subject: IIJmio AutoSwitch: Coupon status changed\r\n\r\n"

	message := "Your coupon status changed as follows:\r\n\r\n"

	message += "- Latest Packet Usage (MB)\r\n"
	message += "Id          WithCoupon WithoutCoupon\r\n"
	message += "--------------------------------------\r\n"
	for k, v := range latestPacketData {
		msg := fmt.Sprintf("%s %10d %12d\r\n", k, v[0], v[1])
		message += msg
	}
	message += "\r\n"
	message += "- Coupon Amount (MB)\r\n"
	message += "Amount\r\n"
	message += "------\r\n"
	message += fmt.Sprintf("%6d\r\n\r\n", couponAmount)
	message += "- Coupon Status\r\n"
	message += "Id          Status\r\n"
	message += "------------------\r\n"
	for k, v := range couponState {
		flag := "OFF"
		if v {
			flag = "ON"
		}
		msg := fmt.Sprintf("%s %s", k, flag)
		if v, ok := couponReqInfo[k]; ok {
			flag = "OFF"
			if v {
				flag = "ON"
			}
			msg += fmt.Sprintf(" -> %s", flag)
		}
		msg += "\r\n"
		message += msg
	}
	body := message + "\r\n"

	return toStr + subjectStr + body
}
