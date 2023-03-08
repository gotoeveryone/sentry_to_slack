# Sentry to Slack

Sentry のデータを Slack に連携するための関数群

## Requirements

- Go 1.19+

## Environment Variables

- SLACK_API_TOKEN: Slack の API トークン
- SLACK_CHANNEL: 通知先のチャンネル
- TAGS: 通知に利用したいタグをカンマ区切りで定義 (例: `server_name,environment`)
