package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var logger *log.Logger
var dynamoClient *dynamodb.Client
var tableName string

// handler receives a batch of SQS messages and persists to DynamoDB.  All messages must be successful or the entire batch will be reprocessed.
func handler(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		if err := handleMessage(ctx, message); err != nil {
			logger.Printf("error: %+v\n", err)
			return err
		}
	}
	logger.Println("processed all SQS messages")
	return nil
}

func main() {
	logger = log.Default()
	logger.SetPrefix("user_creation ")
	logger.SetFlags(log.Lshortfile | log.Lmsgprefix)
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	dynamoClient = newDynamoClient(cfg)
	if env, ok := os.LookupEnv("DYNAMO_TABLE_NAME"); !ok {
		panic(errors.New("missing required env DYNAMO_TABLE_NAME"))
	} else {
		tableName = env
	}
	lambda.Start(handler)
}

func handleMessage(ctx context.Context, message events.SQSMessage) error {
	logger.Printf("handling SQS message %s\n", message.MessageId)
	user, err := newUser(message)
	if err != nil {
		return err
	}
	dynamoInput, err := user.asDynamoInput(tableName)
	if err != nil {
		return err
	}
	if err = saveUserToDynamo(ctx, dynamoClient, dynamoInput); err != nil {
		return err
	}
	return nil
}
