#!/bin/sh
GOOS=linux GOARCH=amd64 go build -o bin/send_used_metric src/functions/send_used_metric/main.go
GOOS=linux GOARCH=amd64 go build -o bin/send_scheduled_metric src/functions/send_scheduled_metric/main.go
rm bin/*.zip
zip bin/send_used_metric bin/send_used_metric
zip bin/send_scheduled_metric bin/send_scheduled_metric
rm bin/send_used_metric
rm bin/send_scheduled_metric
