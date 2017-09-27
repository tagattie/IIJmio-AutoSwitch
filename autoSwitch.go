package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configFileDir  = "/usr/local/etc"
	configFileName = "autoSwitch.json"
	couponEndpoint = "https://api.iijmio.jp/mobile/d/v2/coupon/"
	packetEndpoint = "https://api.iijmio.jp/mobile/d/v2/log/packet/"
	authUrl        = "https://api.iijmio.jp/mobile/d/v1/authorization/?response_type=token&client_id=nWmKQvVQbEfM11PzENM&state=auth-request&redirect_uri=jp.or.iij4u.rr.tagattie.autoswitch"
	authUrlEncoded = "https://api.iijmio.jp/mobile/d/v1/authorization/%3Fresponse_type%3Dtoken%26client_id%3DnWmKQvVQbEfM11PzENM%26state%3Dauth-request%26redirect_uri%3Djp.or.iij4u.rr.tagattie.autoswitch"
)

type mioconf struct {
	DeveloperId    string `json:"developerId"`
	AccessToken    string `json:"accessToken"`
	MaxDailyAmount int    `json:"maxDailyAmount"`
}
type mailconf struct {
	Enabled    bool     `json:"enabled"`
	SmtpServer string   `json:"smtpServer"`
	SmtpPort   string   `json:"smtpPort"`
	ToAddrs    []string `json:"toAddrs"`
	FromAddr   string   `json:"fromAddr"`
	Auth       bool     `json:"auth"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
}
type slackconf struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

type configuration struct {
	Mio   mioconf   `json:"mio"`
	Mail  mailconf  `json:"mail"`
	Slack slackconf `json:"slack"`
}

var (
	configdir *string
	debug     bool
	force     bool
	silent    bool
)
var config configuration

func main() {
	// Command-line options
	configdir = flag.String("c", configFileDir, "config file directory")
	d := flag.Bool("d", false, "debug mode (verbose output)")
	f := flag.Bool("f", false, "force send change request (even when no change)")
	s := flag.Bool("s", true, "silent mode (no output except error msgs)")
	flag.Parse()
	debug = *d
	force = *f
	silent = *s

	// Configuration file
	cf := filepath.Join(*configdir, "/", configFileName)
	if _, err := os.Stat(cf); err != nil {
		fmt.Println("No such file or directory: ", cf)
		fmt.Println("Trying to read config file in current working dir.")
		cwd, _ := os.Getwd()
		cf = filepath.Join(cwd, "/", configFileName)
	}
	file, err := ioutil.ReadFile(cf)
	if err != nil {
		fmt.Println("Config file read error: ", err)
		os.Exit(1)
	}

	// Configuration
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Println("JSON unmarshal error: ", err)
		os.Exit(1)
	}
	if debug == true {
		fmt.Printf("%s\n", "Configuration: ")
		fmt.Printf("%+v\n\n", config)
	}
	if config.Mio.MaxDailyAmount <= 0 {
		fmt.Printf("WARNING: Max daily amount is less than or equal to 0. Coupon use will be always set OFF.\n")
	}

	// Get packet data from server
	packetBytes, err := getData("packet")
	if err != nil {
		fmt.Println("Packet data request error: ", err)
		os.Exit(1)
	}

	// Decode package data
	packetData, err := decodePacketDataJSON(packetBytes)
	if err != nil {
		fmt.Println("JSON data decode error: ", err)
		subjectReason := packetData.ReturnCode
		if config.Mail.Enabled == true {
			if err := sendMail(subjectReason); err != nil {
				fmt.Println("Sending mail error: ", err)
			}
		}
		if config.Slack.Enabled == true {
			if err = sendSlack(subjectReason); err != nil {
				fmt.Println("Sending slack error: ", err)
			}
		}
		os.Exit(1)
	}
	if debug == true {
		fmt.Printf("%s\n", "Packet data JSON:")
		fmt.Printf("%+v\n\n", *packetData)
	}

	// Extract latest packet usage data
	latestPacketData := extractLatestPacketData(packetData)
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Latest packet usage (in MB):")
		fmt.Printf("%+v\n\n", latestPacketData)
	}

	// Get coupon status and amount data from server
	couponBytes, err := getData("coupon")
	if err != nil {
		fmt.Println("Coupon data request error: ", err)
		os.Exit(1)
	}

	// Decode coupon data
	couponData, err := decodeCouponDataJSON(couponBytes)
	if err != nil {
		fmt.Println("JSON data decode error: ", err)
		os.Exit(1)
	}
	if debug == true {
		fmt.Printf("%s\n", "Coupon data JSON:")
		fmt.Printf("%+v\n\n", *couponData)
	}

	// Extract current coupon state and availability
	couponState, couponAmount := getCouponStateAndAmount(couponData)
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Current coupon state and amount:")
		fmt.Printf("%+v %+v\n\n", couponState, couponAmount)
	}

	// Make coupon state change request based on
	// latest packet data and current coupon state, and amount
	couponReqInfo := make(map[string]bool)
	for k, _ := range latestPacketData {
		// latestPacketData[k][0]: Packet data amount with coupon
		// latestPacketData[k][1]: Packet data amount without coupon
		if latestPacketData[k][0] >= config.Mio.MaxDailyAmount &&
			couponState[k] == true {
			couponReqInfo[k] = false
		} else if latestPacketData[k][0] < config.Mio.MaxDailyAmount &&
			couponState[k] == false {
			// Only when there is still coupon amount available
			if couponAmount > 0 {
				couponReqInfo[k] = true
			}
		}
	}
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Coupon status change request:")
		fmt.Printf("%+v\n\n", couponReqInfo)
	}
	// If no need to make change request, exit here
	if len(couponReqInfo) == 0 && !force {
		return
	}

	// Encode coupon request
	couponReqBytes, err := encodeCouponReqJSON(couponReqInfo)
	if err != nil {
		fmt.Println("JSON data encode error: ", err)
		os.Exit(1)
	}
	if debug == true {
		fmt.Printf("%s\n", "Coupon request JSON:")
		fmt.Printf("%+v\n\n", string(couponReqBytes))
	}

	// Send coupon request to server
	couponRespBytes, err := putCouponRequest(couponReqBytes)
	if err != nil {
		fmt.Println("Coupon change request error: ", err)
		os.Exit(1)
	}

	// Decode coupon change response
	couponResp, err := decodeCouponRespJSON(couponRespBytes)
	if err != nil {
		fmt.Println("JSON data decode error: ", err)
		os.Exit(1)
	}
	if silent == false || debug == true {
		fmt.Printf("%s\n", "Counpon change response JSON:")
		fmt.Printf("%+v\n\n", *couponResp)
	}

	return
}
