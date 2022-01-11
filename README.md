# go_to_work

.envファイルにマネフォのログイン情報を入れてください。

MONEY_COMPANY_ID="会社のID"

MONEY_EMAIL=" Your Email "

MONEY_PASSWORD=" Your Password "

あとは実行するだけで今月分の出勤、退勤時間を自動入力します。申請はしません。

`go run main.go`


# エラー対処

### “chromedriver” cannot be opened because the developer cannot be verified.

1 $ which chromedriver

2 $ xattr -d com.apple.quarantine 1で出てきたパス



### invalid session id

Chromeのバージョンとドライバーのバージョンは一致してないとダメ。Chromeはデフォルトで勝手に最新になる使用のため注意が必要。

`brew upgrade chromedriver`
