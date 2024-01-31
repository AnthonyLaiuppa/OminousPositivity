package models

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

type Message struct {
	Id      int    `json:"id" dynamodbav:"id"`
	Message string `json:"message" dynamodbav:"message"`
}

func (m *Message) GetKey() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberN{Value: strconv.Itoa(m.Id)},
	}
}
