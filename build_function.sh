#!/bin/sh
GOOS=linux GOARCH=amd64 go build -o bin/send_webhook src/functions/send_webhook/main.go
rm bin/*.zip
zip bin/send_webhook bin/send_webhook
rm bin/send_webhook
