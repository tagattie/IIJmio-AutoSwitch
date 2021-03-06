# IIJmio-AutoSwitch
[![Build Status](https://travis-ci.org/tagattie/IIJmio-AutoSwitch.svg?branch=master)](https://travis-ci.org/tagattie/IIJmio-AutoSwitch)
[![Go Report Card](https://goreportcard.com/badge/github.com/tagattie/IIJmio-AutoSwitch)](https://goreportcard.com/report/github.com/tagattie/IIJmio-AutoSwitch)

[English](README.md)

IIJmio向け: クーポン使用量が所定の値を超えたとき、クーポン使用を自動的にOFFします。

## 背景
このプログラムは[IIJmioモバイル](https://www.iijmio.jp/)の利用者向けです。

動画や音楽のストリーミングサービスでパケットを使いすぎて月末にクーポン残量がなくなってしまったり、ファミリーシェアプランでクーポンをシェアしているメンバーがひとりでパケットを使いすぎたりして困ることがあります。

Androidには[Mio Mix](https://play.google.com/store/apps/details?id=com.itworks.miomix)というアプリがあり、一日あたり所定の通信量を超えた場合にクーポン使用を自動オフする機能を提供しています。

しかし、iOSにはこういったアプリがないようです。また、AndroidでもOSによるバックグラウンドタスクキル機能などの影響でアプリがうまく動作しないケースもあります。このプログラムは、Unix-like OSあるいはWindowsベースの端末、サーバ上で動作し、設定した一日あたりのクーポン使用量を超えたときに、クーポン使用を自動でOFFにします。

## 機能
以下の機能を提供します。

- クーポン使用の自動ON/OFF
  - クーポン使用ONのとき:
    - 事前に設定した一日あたりのクーポン使用量を超えたら、クーポン使用をOFFにします
  - クーポン使用OFFのとき:
    - 一日あたりのクーポン使用量を下回り、かつクーポンが残っていればクーポン使用をONにします
    - (日付が変わって通信量がクリアされるときを想定しています)
- メール送信機能有効化時:
  - クーポンON/OFF状態変化時に設定した宛先にメール送信
  - 認証エラー発生時に設定した宛先にメール送信
- 認証エラー発生時に設定したチャネルにSlackメッセージを送信(有効化時のみ)

## 動作環境
- Unix-like OS / Windows
- Go
- GNU Make

以下の環境で動作確認をしています。

- FreeBSD 11.1-RELEASE 64bit (Go 1.9, GNU Make 4.2.1)
- Windows 10 64bit (Go 1.9, GNU Make 4.2.1)

## ビルド
以下のコマンドを実行します:

```sh
go get github.com/tagattie/IIJmio-AutoSwitch
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
make
```

`${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/mioswitch`に実行ファイルが生成されます。(Windowsの場合は、実行ファイル名が`mioswitch.exe`となります。)

## 設定
このプログラムは、IIJmioが提供する[クーポンスイッチAPI](https://www.iijmio.jp/hdd/coupon/mioponapi.jsp)を使用します。以下で、プログラムの動作に必要な設定を行ないます。

### アクセストークンの取得
まず、認証用のアクセストークンを取得します。クーポンスイッチAPIの注意事項に、「認可は標準ブラウザなどを利用して行ってください。アプリケーション側にmioIDとパスワードを入力させて認可する実装はしないでください。」とありますので、Webブラウザで以下の手順を実行してください。

1. 以下のURLにアクセスします:

    <https://api.iijmio.jp/mobile/d/v1/authorization/?response_type=token&client_id=nWmKQvVQbEfM11PzENM&state=auth-request&redirect_uri=jp.or.iij4u.rr.tagattie.autoswitch>

1. mioIDとmioパスワードを入力してログインします。
1. 「使いすぎ防止オートスイッチが以下の機能の許可を求めています。...」という確認画面になりますので、連携を許可する場合は「許可する」のボタンを押下します。
1. 許可すると、以下のURLにリダイレクトされます:

    <https://api.iijmio.jp/mobile/d/v1/authorization/jp.or.iij4u.rr.tagattie.autoswitch#access_token=YOUR_ACCESS_TOKEN&state=auth-request&token_type=Bearer&expires_in=7776000>

    この際、待ち受けているアプリがいないので、ブラウザには

        {"returnCode": "Requested resource is not found"}

    と表示されますが、URLには必要なアクセストークンが含まれていますので、これを記録しておきます。(URLの`access_token=YOUR_ACCESS_TOKEN`の部分) アクセストークンは90日間有効です。トークンの有効期限が切れた場合は、本手順を改めて行ない、アクセストークンを再取得する必要があります。

### 設定
このプログラムは、デフォルトで

- `/usr/local/etc` あるいは
- カレントワーキングディレクトリ

に`mioswitch.json`という名前の設定ファイルがあることを期待しています。プログラムを実行する前に、以下の設定を行ないます:

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

- 設定項目
  - `mio`: クーポンスイッチAPI関連の設定
    - `developerId`: クーポンスイッチAPIの認証に使用する開発者IDです。変更しないでください。
    - `accessToken`: 認証に使用するアクセストークンです。上記の手順で取得した値を設定してください。
    - `maxDailyAmount`: ここで設定した値(MB)を超えるとクーポン使用をOFFにします。
    - `startingAmount`: 月初におけるクーポン残量(MB)を指定します。
  - `switch`:
    - `switchMethod`: クーポンの使用量を求める方式を指定します。以下に示す1あるいは2の方式が指定できます。
      1. ユーザのクーポン使用量データを直接使用します。クーポン使用量が上記の設定値を上回った場合に、クーポン使用をOFFにします。
      1. ユーザのクーポン使用量とクーポンのトータル残量データを使用します。上記の設定値から、当日までに残っているべきクーポン残量を算出し、クーポン残量がこの値を下回った場合に、使用量の最も多いユーザのクーポン使用をOFFにします。
      
      ユーザのクーポン使用量データはタイムリーに更新されないようなので、2の方式を推奨します。
  - `mail`: メール関連の設定
    - `enabled`: メール送信機能の有効化(trueで有効)。
    - `smtpServer`: メール送信に使用するサーバーを指定します。
    - `smtpPort`: メール送信に使用するサーバーのポート番号を指定します。
    - `toAddrs`: メール送信先のアドレスを指定します(複数可)。
    - `fromAddr`: メールの送信元となるアドレスを指定します。
    - `auth`: メールサーバーが認証を必要とする場合、trueを指定します。
    - `username`: 認証に使用するユーザー名を指定します。
    - `password`: 認証に使用するパスワードを指定します。
  - `slack`: Slack関連の設定
    - `enabled`: Slackメッセージ送信機能の有効化(trueで有効)。
    - `token`: Slack APIの認証に使用するトークンを指定します。トークンの取得については[Slack Web API](https://api.slack.com/web)を参照してください。
    - `channel`: メッセージを送信するSlackのチャネル名を指定します。

## 実行
以下のコマンドを実行します:

```sh
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
./bin/mioswitch
```

コマンドラインオプションの一覧は以下で確認できます:

```sh
./bin/mioswitch -h
```

## 定期的な実行
Cron(あるいは同等のプログラム)を使用します。(以下の例では、毎時15分にプログラムを実行します。)

    #min hour mday mon wday command
    15   *    *    *   *    ${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/mioswitch
