package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyAttribute struct {
	DEVICE string `json:"DEVICE"`
}

type MyPlacementInfo struct {
	Attributes MyAttribute `json:"attributes"`
}

type MyEvent struct {
	PlacementInfo MyPlacementInfo `json:"placementInfo"`
}

// HandleRequest puts lastmodified to s3
func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	log.Print(event)
	var webhookURL = os.Getenv("WEBHOOK_URL")
	params, err := json.Marshal(event)
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
