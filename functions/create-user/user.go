package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

// User represents a User Entity
type User struct {
	ID           string
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`
	SQSMessageID string    `json:"sqsMessageId"`
}

func newUser(sqsMessage events.SQSMessage) (*User, error) {
	var user User
	if err := json.Unmarshal([]byte(sqsMessage.Body), &user); err != nil {
		return nil, err
	}
	user.CreatedAt = time.Now()
	user.SQSMessageID = sqsMessage.MessageId
	user.ID = uuid.NewString()
	return &user, nil
}

func (u *User) asDynamoInput(tableName string) (*dynamodb.PutItemInput, error) {
	item, err := attributevalue.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}, nil
}
