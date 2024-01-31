package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log/slog"
	"os"
	"time"
)

type DynamoDBRepository struct {
	client    DynamoDBClientInterface
	tableName string
}

type DynamoDBClientInterface interface {
	GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(ctx context.Context, input *dynamodb.DeleteTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
}

func CreateDynamoDBClient(ctx context.Context, useLocal bool) (DynamoDBClientInterface, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	var cfg aws.Config
	var err error

	if useLocal {
		endpoint := os.Getenv("DYNAMO_ENDPOINT")
		if endpoint == "" {
			endpoint = "http://localhost:8000"
		}

		// Setup for local DynamoDB
		cfg, err = config.LoadDefaultConfig(c, config.WithRegion("local"),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint}, nil
				})))
	} else {
		// Default configuration
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	return svc, nil
}

func NewDynamoDBRepository(client DynamoDBClientInterface, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

func (repo *DynamoDBRepository) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return repo.client.GetItem(ctx, input, opts...)
}

func (repo *DynamoDBRepository) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return repo.client.PutItem(ctx, input, opts...)
}

func (repo *DynamoDBRepository) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return repo.client.CreateTable(ctx, input, opts...)
}

func (repo *DynamoDBRepository) DeleteTable(ctx context.Context, input *dynamodb.DeleteTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	return repo.client.DeleteTable(ctx, input, opts...)
}

func (repo *DynamoDBRepository) GetMessageById(ctx context.Context, id int) (*Message, error) {
	message := &Message{Id: id}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       message.GetKey(),
	}

	result, err := repo.client.GetItem(ctx, input)
	if err != nil {
		slog.ErrorContext(ctx, "encountered error attempting to get message", "ID", id, "Table", repo.tableName, "Error", err)
		return nil, err
	}

	if result.Item == nil {
		slog.WarnContext(ctx, "message not found", "ID", id, "Table", repo.tableName)
		return nil, nil // Message not found
	}

	err = attributevalue.UnmarshalMap(result.Item, message)
	if err != nil {
		slog.ErrorContext(ctx, "encountered error attempting to unmarshal message", "ID", id, "Table", repo.tableName, "Error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "retrieved message", "ID", id, "Table", repo.tableName)
	return message, nil
}
