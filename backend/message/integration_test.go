//go:build integration
// +build integration

package main

import (
	"context"
	"github.com/AnthonyLaiuppa/ominouspositivity/backend/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageHandlerIntegration(t *testing.T) {
	// Set up stuff, context and DB client
	tableName := "integration_test"
	ctx := context.TODO()
	dynamoClient, err := models.CreateDynamoDBClient(ctx, true)
	assert.NoError(t, err)
	db := models.NewDynamoDBRepository(dynamoClient, tableName)

	// Create a new Table for the test
	attributeDefinitions := []types.AttributeDefinition{
		{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeN, // 'N' for number
		},
	}
	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash, // 'HASH' key type for the primary key
		},
	}
	provisionedThroughput := &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(5), // Specify your read capacity units
		WriteCapacityUnits: aws.Int64(5), // Specify your write capacity units
	}
	_, err = db.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchema,
		TableName:             aws.String(tableName),
		ProvisionedThroughput: provisionedThroughput,
	})
	assert.NoError(t, err)
	// Test a regular GET against an Empty DB
	// Create a request object simulating an API Gateway event
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}
	// Call the messageHandler function
	resp, err := messageHandler(context.TODO(), req, db)
	assert.NoError(t, err)
	// Assert that the response is as expected based on the test data in DynamoDB
	assert.Equal(t, 200, resp.StatusCode)
	// Assertion that response body should contain default message since DB is unpopulated
	assert.Contains(t, resp.Body, "What is lost may find you before you're ready.")
	// Loop to seed data 13 times because I thought RNG was a cute way to pick the messages
	var m models.Message
	for i := 0; i <= 12; i++ {
		m.Id = i
		m.Message = "Test"
		item, err := attributevalue.MarshalMap(m)
		assert.NoError(t, err)
		_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		assert.NoError(t, err)
	}

	//Test to getting a message from a populated table
	resp, err = messageHandler(context.TODO(), req, db)
	assert.NoError(t, err)
	// Assert that the response is as expected based on the test data in DynamoDB
	assert.Equal(t, 200, resp.StatusCode)
	// Assertion that response body should contain default message since DB is unpopulated
	assert.Contains(t, resp.Body, "Test")

	// Tear down the Table
	_, err = db.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	assert.NoError(t, err)
}
