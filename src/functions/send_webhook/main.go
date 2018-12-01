package main

import (
	"context"
	"log"

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
	return "normal end", nil
}

func main() {
	lambda.Start(HandleRequest)
}
