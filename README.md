# IIJmio-AutoSwitch
[![Build Status](https://travis-ci.org/tagattie/IIJmio-AutoSwitch.svg?branch=master)](https://travis-ci.org/tagattie/IIJmio-AutoSwitch)

Automatically disable IIJmio coupon use while packet usage is over preset amount.

IIJmio向け: クーポン使用量が所定の値を超えたとき、クーポン使用を自動的にOFFします。

## 背景
このプログラムは[IIJmioモバイル](https://www.iijmio.jp/)の利用者向けです。

動画や音楽のストリーミングサービスでパケットを使いすぎて月末にクーポン残量がなくなってしまったり、ファミリーシェアプランでクーポンをシェアしているメンバーがひとりでパケットを使いすぎたりして困ることがあります。

Androidには[Mio Mix](https://play.google.com/store/apps/details?id=com.itworks.miomix)というアプリがあり、一日あたり所定の通信量を超えた場合にクーポン使用を自動オフする機能を提供しています。

しかし、iOSにはこういったアプリがないようです。また、AndroidでもOSなどによるバックグラウンドタスクキル機能の影響でアプリがうまく動作しないケースもあります。このプログラムは、Unix-like OSベースの端末、サーバ上で動作し、設定した一日あたりのクーポン使用量を超えたときに、クーポンを自動でOFFにします。

## 機能
以下の機能を提供します。

- クーポン使用ONのとき:
  - 事前に設定した一日あたりのクーポン使用量を超えたら、クーポン使用をOFFにします
- クーポン使用OFFのとき:
  - 一日あたりのクーポン使用量を下回り、かつクーポンが残っていればクーポン使用をONにします
  - (日付が変わって通信量がクリアされるときを想定しています)

## 動作環境
- Unix-like OS
- Go
- GNU Make

以下の環境で動作確認をしています。

- Ubuntu 14.04.5 (Go 1.2.1, GNU Make 3.81)
- FreeBSD 11.1-RELEASE (Go 1.9, GNU Make 4.2.1)

## ビルド
以下のコマンドを実行します:

```sh
go get github.com/tagattie/IIJmio-AutoSwitch
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
make
```

`${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/autoswitch`に実行ファイルが生成されます。

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

に`autoSwitch.json`という名前の設定ファイルがあることを期待しています。プログラムを実行する前に、以下の設定を行ないます:

```json:autoSwitch.json
{
    "accessToken": "YOUR_ACCESS_TOKEN"
}
```

取得したアクセストークンを設定してください。

## 実行
以下のコマンドを実行します:

```sh
cd ${GOPATH-$HOME/go}/src/github.com/tagattie/IIJmio-AutoSwitch
./bin/autoswitch
```

コマンドラインオプションの一覧は以下で確認できます:

```sh
./bin/autoswitch -h
```

## 定期的な実行
Cron(あるいは同等のプログラム)を使用します。(以下の例では、毎時15分にプログラムを実行します。)

    #minute hour    mday    month   wday    command
    15      *       *       *       *       ${GOPATH}/src/github.com/tagattie/IIJmio-AutoSwitch/bin/autoswitch
