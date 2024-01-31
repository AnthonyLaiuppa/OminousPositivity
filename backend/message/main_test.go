package main

import (
	"context"
	"encoding/json"
	"github.com/AnthonyLaiuppa/ominouspositivity/backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock for DynamoDBClientInterface
type MockDynamoDBClient struct {
	mock.Mock
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) GetMessageById(ctx context.Context, id int) (*models.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Message), args.Error(1)
}

func TestMessageHandler(t *testing.T) {
	ctx := context.TODO()
	mockDB := new(MockDynamoDBClient)

	// Set up the mock expectation for a successful response
	mockMsg := &models.Message{Id: 1, Message: "Test message"}
	mockDB.On("GetMessageById", ctx, mock.AnythingOfType("int")).Return(mockMsg, nil)

	// Test for a successful response
	req := events.APIGatewayProxyRequest{HTTPMethod: "GET"}
	res, err := messageHandler(ctx, req, mockDB)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var msg map[string]string
	err = json.Unmarshal([]byte(res.Body), &msg)
	assert.NoError(t, err)
	assert.Equal(t, "Test message", msg["message"])

	// Test for OPTIONS request
	req = events.APIGatewayProxyRequest{HTTPMethod: "OPTIONS"}
	res, err = messageHandler(ctx, req, mockDB)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// Reset mock expectations for the item not found scenario
	mockDB.ExpectedCalls = nil
	mockDB.Calls = nil
	mockDB.On("GetMessageById", ctx, mock.AnythingOfType("int")).Return((*models.Message)(nil), nil)

	// Test for Item Not Found
	req = events.APIGatewayProxyRequest{HTTPMethod: "GET"}
	response, _ := messageHandler(ctx, req, mockDB)
	assert.Contains(t, response.Body, "What is lost may find you before you're ready.")
}
