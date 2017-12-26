package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	configFileDir  = "/usr/local/etc"
	configFileName = "mioswitch.json"
	logFilePath    = "/var/tmp/mioswitch.log"
	couponEndpoint = "https://api.iijmio.jp/mobile/d/v2/coupon/"
	packetEndpoint = "https://api.iijmio.jp/mobile/d/v2/log/packet/"
	authUrl        = "https://api.iijmio.jp/mobile/d/v1/authorization/?response_type=token&client_id=nWmKQvVQbEfM11PzENM&state=auth-request&redirect_uri=jp.or.iij4u.rr.tagattie.autoswitch"
	authUrlEncoded = "https://api.iijmio.jp/mobile/d/v1/authorization/%3Fresponse_type%3Dtoken%26client_id%3DnWmKQvVQbEfM11PzENM%26state%3Dauth-request%26redirect_uri%3Djp.or.iij4u.rr.tagattie.autoswitch"
)

type mioconf struct {
	DeveloperId    string `json:"developerId"`
	AccessToken    string `json:"accessToken"`
	MaxDailyAmount int    `json:"maxDailyAmount"`
	StartingAmount int    `json:"startingAmount"`
}
type switchconf struct {
	SwitchMethod int `json:"switchMethod"`
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
	Mio    mioconf    `json:"mio"`
	Switch switchconf `json:"switch"`
	Mail   mailconf   `json:"mail"`
	Slack  slackconf  `json:"slack"`
}

var (
	configdir *string
	debug     bool
	force     bool
	silent    bool
)
var config configuration

func main() {
	// Setup logger
	logfile, err := os.OpenFile(
		logFilePath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644)
	if err != nil {
		log.Fatalln("Cannot open log file: ", logFilePath)
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Command-line options
	log.Println("Parsing command-line options:")
	configdir = flag.String("c", configFileDir, "config file directory")
	d := flag.Bool("d", false, "debug mode (verbose output)")
	f := flag.Bool("f", false, "force send change request (even when no change)")
	s := flag.Bool("s", true, "silent mode (no output except error msgs)")
	flag.Parse()
	debug = *d
	force = *f
	silent = *s
	log.Printf("debug = %+v, force = %+v, silent = %+v\n", debug, force, silent)

	// Configuration file
	log.Println("Reading configuration file:")
	cf := filepath.Join(*configdir, "/", configFileName)
	if _, err := os.Stat(cf); err != nil {
		log.Println("No such file or directory: ", cf)
		log.Println("Try to read config file in the current working dir.")
		cwd, _ := os.Getwd()
		cf = filepath.Join(cwd, "/", configFileName)
	}
	file, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatalln("Configuration file read error: ", err)
	}

	// Configuration
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalln("Configuration JSON unmarshal error: ", err)
	}
	if config.Mio.MaxDailyAmount <= 0 {
		log.Fatalln("Max daily amount must be a positive number.")
	}
	log.Println("Configuration:")
	log.Printf("%+v\n", config)

	// Get packet data from server
	packetBytes, err := getData("packet")
	if err != nil {
		log.Fatalln("Packet data request error: ", err)
	}

	// Decode packet data
	packetData, err := decodePacketDataJSON(packetBytes)
	if err != nil {
		log.Println("JSON data decode error: ", err)
		subjectReason := packetData.ReturnCode
		if config.Mail.Enabled {
			message := buildErrorMessage(subjectReason)
			if err := sendMail(message); err != nil {
				log.Println("Sending mail error: ", err)
			}
		}
		if config.Slack.Enabled {
			if err = sendSlack(subjectReason); err != nil {
				log.Println("Sending slack error: ", err)
			}
		}
		os.Exit(1)
	}
	log.Println("Packet data JSON:")
	log.Printf("%+v\n", *packetData)

	// Get coupon status and amount data from server
	couponBytes, err := getData("coupon")
	if err != nil {
		log.Fatalln("Coupon data request error: ", err)
	}

	// Decode coupon data
	couponData, err := decodeCouponDataJSON(couponBytes)
	if err != nil {
		log.Fatalln("JSON data decode error: ", err)
	}
	log.Println("Coupon data JSON:")
	log.Printf("%+v\n", *couponData)

	var latestPacketData map[string][]int
	var couponState map[string]bool
	var couponAmount int
	var couponReqInfo map[string]bool
	switch config.Switch.SwitchMethod {
	case 0:
		latestPacketData, couponState, couponAmount, couponReqInfo =
			couponChangeByIdBasedCouponUsage(packetData, couponData)
	case 1:
		latestPacketData, couponState, couponAmount, couponReqInfo =
			couponChangeByRemainingCouponAmount(packetData, couponData)
	}
	// If no need to make change request, exit here
	if len(couponReqInfo) == 0 && !force {
		return
	}

	// Encode coupon request
	couponReqBytes, err := encodeCouponReqJSON(couponReqInfo)
	if err != nil {
		log.Fatalln("JSON data encode error: ", err)
	}
	log.Println("Coupon request JSON:")
	log.Printf("%+v\n", string(couponReqBytes))

	// Send coupon request to server
	couponRespBytes, err := putCouponRequest(couponReqBytes)
	if err != nil {
		log.Fatalln("Coupon change request error: ", err)
	}

	// Decode coupon change response
	couponResp, err := decodeCouponRespJSON(couponRespBytes)
	if err != nil {
		log.Fatalln("JSON data decode error: ", err)
	}
	log.Println("Counpon change response JSON:")
	log.Printf("%+v\n", *couponResp)

	// Send coupon change report mail
	if config.Mail.Enabled {
		message := buildReportMessage(latestPacketData,
			couponState,
			couponAmount,
			couponReqInfo)
		if err := sendMail(message); err != nil {
			log.Fatalln("Sending mail error: ", err)
		}
	}

	return
}
