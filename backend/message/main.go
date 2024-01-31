package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AnthonyLaiuppa/ominouspositivity/backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
)

type RepositoryInterface interface {
	GetMessageById(ctx context.Context, id int) (*models.Message, error)
}

var _ RepositoryInterface = (*models.DynamoDBRepository)(nil)

func messageHandler(ctx context.Context, request events.APIGatewayProxyRequest, db RepositoryInterface) (events.APIGatewayProxyResponse, error) {
	//CORS Stuff
	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      allowOrigin,
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Content-Type",
				"Access-Control-Allow-Methods":     "GET,OPTIONS",
			},
		}, nil
	}
	// Retrieve a random message from the DB
	id := rand.Intn(12)
	m, err := db.GetMessageById(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "error marshalling message response for message", "ID", id)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	message := make(map[string]string)
	if m == nil {
		message["message"] = "What is lost may find you before you're ready."
	} else {
		message["message"] = m.Message
	}
	// Pack the message into a JSON response
	messageBytes, err := json.Marshal(message)
	if err != nil {
		slog.ErrorContext(ctx, "error marshalling message response for message", "ID", id)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	// Return the message
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(messageBytes),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": allowOrigin,
		},
	}, nil
}

func main() {
	allowOrigin := os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		panic(fmt.Errorf("must set ALLOW_ORIGIN"))
	}
	tableName := os.Getenv("TABLE_NAME")
	useLocalStr := os.Getenv("USE_LOCAL") // Set to true to use local DynamoDB
	useLocal := false
	if useLocalStr != "" {
		var err error
		useLocal, err = strconv.ParseBool(useLocalStr)
		if err != nil {
			panic(fmt.Errorf("invalid value for USE_LOCAL: %v", err))
		}
	}
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		dynamoClient, err := models.CreateDynamoDBClient(ctx, useLocal)
		db := models.NewDynamoDBRepository(dynamoClient, tableName)
		if err != nil {
			slog.ErrorContext(ctx, "error setting DynamoDB connection", "useLocal", useLocal, "tableName", tableName)
			return events.APIGatewayProxyResponse{}, err
		}
		return messageHandler(ctx, request, db)
	})
}
