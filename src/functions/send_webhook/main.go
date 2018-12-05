package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyAttribute struct {
	Room string `json:"room"`
}

type MyPlacementInfo struct {
	Attributes MyAttribute `json:"attributes"`
}

type MyEvent struct {
	PlacementInfo MyPlacementInfo `json:"placementInfo"`
}

type Payload struct {
	Message string  `json:"message"`
	Meta    MyEvent `json:"meta"`
}

// HandleRequest puts lastmodified to s3
func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	log.Print(event)
	var webhookURL = os.Getenv("WEBHOOK_URL")
	var message = os.Getenv("MESSAGE")
	log.Print(message)
	log.Print(fmt.Sprintf(message, event.PlacementInfo.Attributes.Room))
	var payload = Payload{
		Message: fmt.Sprintf(message, event.PlacementInfo.Attributes.Room),
		Meta:    event}
	params, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		webhookURL,
		bytes.NewBuffer(params),
	)
	if err != nil {
		return "", err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "normal end", nil
}

func main() {
	lambda.Start(HandleRequest)
}
