package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockDynamoDBClient is a mock implementation of DynamoDBClientInterface
type MockDynamoDBClient struct {
	mock.Mock
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*dynamodb.CreateTableOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteTable(ctx context.Context, input *dynamodb.DeleteTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*dynamodb.DeleteTableOutput), args.Error(1)
}

func TestGetMessageById(t *testing.T) {
	ctx := context.Background()
	mockDynamoDBClient := new(MockDynamoDBClient)
	tableName := "testTable"

	repo := NewDynamoDBRepository(mockDynamoDBClient, tableName)

	// Prepare mock DynamoDB GetItem response
	mockMessage := &Message{Id: 1, Message: "Test message"}
	item, _ := attributevalue.MarshalMap(mockMessage)
	mockDynamoDBClient.On("GetItem", ctx, mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{Item: item}, nil)

	// Test GetMessageById
	result, err := repo.GetMessageById(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockMessage.Message, result.Message)

	// Additional tests can be added for different scenarios such as error handling or item not found
}
