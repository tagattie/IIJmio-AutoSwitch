# IIJmio-AutoSwitch (mioswitch)
[![Build Status](https://travis-ci.org/tagattie/IIJmio-AutoSwitch.svg?branch=master)](https://travis-ci.org/tagattie/IIJmio-AutoSwitch)
[![Go Report Card](https://goreportcard.com/badge/github.com/tagattie/IIJmio-AutoSwitch)](https://goreportcard.com/report/github.com/tagattie/IIJmio-AutoSwitch)

[日本語](README.ja.md)

Automatically disable IIJmio coupon use while packet usage is over preset amount.

## Background
mioswitch is a tool for users of [IIJmio Mobile](https://www.iijmio.jp/).

Don't you run out of your coupon when nearing month-end due to (over)use of movie and/or music streaming services? Doesn't one of the members of your "family share plan" consume coupon too much?

There is an android application called [Mio Mix](https://play.google.com/store/apps/details?id=com.itworks.miomix&hl=en), which automatically switches coupon use off when packet usage is over a preset amount of maximum daily use.

However, iOS doesn't seem to have that kind of application. In case of android, the application may not function well because of "background task kill" of the OS.

mioswitch runs on Unix-like OS and Windows-based computers and turns coupon use off when a preset daily limit is reached.

## Functions
mioswitch provides the following functions:

- Automatic on/off of coupon use
  - When coupon use is set on:
    - Turns coupon use off when a preset daily limit is reached.
  - When coupon use is set off:
    - Turns coupon use on when date changes (and coupon still remains).
- Notification by email (when enabled):
  - Send an email when coupon use is switched.
  - Send an email when authentication error occurs.
- Notification to Slack (when enabled):
  - Send a notification when authentication error occurs.

## Requirements
- Unix-like OS / Windows
- Go
- GNU Make

Tested with the following environments:

- FreeBSD 11.1-RELEASE 64bit (Go 1.9, GNU Make 4.2.1)
- Windows 10 64bit (Go 1.9, GNU Make 4.2.1)

## Build
Execute the following commands:

```sh
go get github.com/tagattie/IIJmio-AutoSwitch
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
make
```

An executable file will be created at `${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/mioswitch`. (In case of Windows, the executable filename will be `mioswitch.exe`.)

## Configuration
mioswitch uses [Coupon Switch API](https://www.iijmio.jp/hdd/coupon/mioponapi.jsp) provided by IIJmio. Prior to using the tool, the following configuration will be required.

### Acquisition of access token
First, you need to get an access token for authenticating accesses to the Coupon Switch API. IIJmio requires "authorization should be done using web browsers." So, follow the procedure below with with a web browser.

1. Access the following URL:

    <https://api.iijmio.jp/mobile/d/v1/authorization/?response_type=token&client_id=nWmKQvVQbEfM11PzENM&state=auth-request&redirect_uri=jp.or.iij4u.rr.tagattie.autoswitch>

1. Login with your mio ID and password.
1. A confirmation saying "使いすぎ防止オートスイッチが以下の機能の許可を求めています。..." will be displayed. If you are OK, click on the "許可する" button.
1. After clicking, the browser will be redirected to the following URL:

    <https://api.iijmio.jp/mobile/d/v1/authorization/jp.or.iij4u.rr.tagattie.autoswitch#access_token=YOUR_ACCESS_TOKEN&state=auth-request&token_type=Bearer&expires_in=7776000>

    The browser contents shows the error message:

        {"returnCode": "Requested resource is not found"}

    However, the access token has been successfully acquired and is included in the redirected URL (The part `access_token=YOUR_ACCESS_TOKEN`). The token is valid for 90 days. After token expiration, you will be required to follow this procedure again.

### Configuration file
mioswitch expects a configuration file named `mioswitch.json` exists in one of the following directories:

- `/usr/local/etc`
- Current working directory

The contents of the configuration file is as follows:

```json:mioswitch.json
{
  "mio": {
    "developerId":    "nWmKQvVQbEfM11PzENM",
    "accessToken":    "YOUR_ACCESS_TOKEN",
    "maxDailyAmount": 100,
    "startingAmount": 10000
  },
  "switch": {
    "switchMethod": 1
  },
  "mail": {
    "enabled":    false,
    "smtpServer": "smtp.example.com",
    "smtpPort":   "587",
    "toAddrs":    [
      "someone1@example.com",
      "someone2@example.com"
    ],
    "fromAddr":   "autoswitch@example.com",
    "auth":       true,
    "username":   "authusername",
    "password":   "authpassword"
  },
  "slack": {
    "enabled": false,
    "token":   "slacktoken",
    "channel": "channelname",
  }
}
```

- Items
  - `mio`: Coupon Switch API settings
      - `developerId`: Developer ID used for authentication. Do not change.
      - `accessToken`: Access token acquired by the above-mentioned procedure.
      - `maxDailyAmount`: Preset daily limit (MB).
      - `startingAmount`: Total coupon amount at the beginning of a month (MB). (Determined by your contract.)
  - `switch`: Switch method settings
      - `switchMethod`: Specify method of coupon switch. One of the following can be set:
      1. Directly uses coupon usage data of a user. When the usage is over the preset limit, mioswitch will turn off the user's coupon use.
      1. Use coupon usage and total remaining amount data. When remaining amount becomes less than an amount to remain at the date, mioswitch turnes off the coupon use of the user whose coupon use is largest.
      
      Per-user coupon usage data seem to take time to refresh. So, use of method 2 is advised.
  - `mail`: Email settings
      - `enabled`: Email function enabled if true.
      - `smtpServer`: SMTP server.
      - `smtpPort`: SMTP port.
      - `toAddrs`: Addresses to which emails will be sent.
      - `fromAddr`: From address of emails.
      - `auth`: Set true if the SMTP server requires authentication.
      - `username`: Username for authentication.
      - `password`: Password for authentication.
  - `slack`: Slack settings
      - `enabled`: Slack function enabled if true.
      - `token`: Token for authenticating Slack API access. Refer to [Slack Web API](https://api.slack.com/web) for acquiring tokens.
      - `channel`: Channel name to which notification will be sent.

## Execution
Execute the following commands:

```sh
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
./bin/mioswitch
```

Execute the following command for available options:

```sh
./bin/mioswitch -h
```

## Scheduled execution
Use cron (of equivalents). (In the following example, the program will execute at 15 minute of every hour.)

    #min hour mday mon wday command
    15   *    *    *   *    ${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/mioswitch
