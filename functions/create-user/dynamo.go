package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

type DynamoPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func newDynamoClient(cfg aws.Config) *dynamodb.Client {
	awsv2.AWSV2Instrumentor(&cfg.APIOptions)
	return dynamodb.NewFromConfig(cfg)
}

func saveUserToDynamo(ctx context.Context, api DynamoPutItemAPI, input *dynamodb.PutItemInput) error {
	_, err := api.PutItem(ctx, input)
	return err
}
