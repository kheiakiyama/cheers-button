package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kheiakiyama/iot-button-metrics/src/lib"
	"github.com/nlopes/slack"
)

type SlashCommand struct {
	Token          string `json:"token"`
	TeamID         string `json:"team_id"`
	TeamDomain     string `json:"team_domain"`
	EnterpriseID   string `json:"enterprise_id,omitempty"`
	EnterpriseName string `json:"enterprise_name,omitempty"`
	ChannelID      string `json:"channel_id"`
	ChannelName    string `json:"channel_name"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Command        string `json:"command"`
	Text           string `json:"text"`
	ResponseURL    string `json:"response_url"`
	TriggerID      string `json:"trigger_id"`
}

// HandleRequest puts metrics based lastmodified
func HandleRequest(ctx context.Context, param SlashCommand) (string, error) {
	log.Print(param)
	var slackVerifiedToken = os.Getenv("SLACK_VERIFIED_TOKEN")
	if slackVerifiedToken != param.Token {
		return "", errors.New("Token Invalid")
	}
	switch param.Command {
	default:
		return DefaultCommand(ctx, param, slackVerifiedToken)
	}
}

// DefaultCommand write recent summary uses
func DefaultCommand(ctx context.Context, param SlashCommand, slackVerifiedToken string) (string, error) {
	var BUCKET = os.Getenv("BUCKET")
	var lastModifiedKeyPrifix = os.Getenv("LASTMODIFIED_KEY_PRIFIX")
	var buttonPrifix = os.Getenv("BUTTON_PREFIX")
	var buttonCount, _ = strconv.Atoi(os.Getenv("BUTTON_COUNT"))
	var TIMEOUT, _ = strconv.ParseInt(os.Getenv("TIMEOUT"), 10, 64)

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(endpoints.ApNortheast1RegionID),
	})

	// 最終
	lastClicked, message, errlo := lib.GetLastButtonClicked(svc)
	if errlo != nil {
		return message, errlo
	}

	var usedCount = 0
	for index := 1; index <= buttonCount; index++ {
		var key = fmt.Sprintf("%s/%s%d", lastModifiedKeyPrifix, buttonPrifix, index)
		log.Print(key)

		// Object取得
		goo, errgo := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(BUCKET),
			Key:    &key,
		})
		var getSuccess = true
		if errgo != nil {
			if aerr, ok := errgo.(awserr.Error); ok {

				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					log.Print("bucket does not exist at GetObject")
					return "bucket does not exist at GetObject", aerr

				case s3.ErrCodeNoSuchKey:
					// 新規作成
					getSuccess = false
				default:
					log.Printf("aws error %v at GetObject", aerr.Error())
					return "aws error at GetObject", aerr
				}
			}
		}
		if getSuccess {
			defer goo.Body.Close()
		}
		t := time.Now()
		inTime := (t.Unix() - lastClicked[fmt.Sprintf("%s%d", buttonPrifix, index)]) < TIMEOUT
		if inTime {
			usedCount++
		}
	}
	defer log.Print("normal end")
	msg := fmt.Sprintf("%d used / %d all", usedCount, buttonCount)
	api := slack.New(slackVerifiedToken)
	api.PostMessage(param.ChannelID, msg, slack.PostMessageParameters{})
	return "normal end", nil
}

func main() {
	lambda.Start(HandleRequest)
}
