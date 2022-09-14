package main

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestNewUser(t *testing.T) {
	cases := map[string]struct {
		given     events.SQSMessage
		want      *User
		shouldErr bool
	}{
		"valid message returns user": {
			given: events.SQSMessage{
				MessageId: "messageID",
				Body:      "{\"name\":\"person\",\"email\":\"email@test.com\"}",
			},
			want: &User{
				SQSMessageID: "messageID",
				Name:         "person",
				Email:        "email@test.com",
			},
		},
		"empty message returns error": {
			given: events.SQSMessage{
				MessageId: "bad-message-ID",
				Body:      "",
			},
			shouldErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := newUser(tt.given)
			if tt.shouldErr {
				if err == nil {
					t.Error("expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s\n", err)
				}
				if got.Name != tt.want.Name {
					t.Errorf("incorrect name: %s\n", got.Name)
				}
				if got.Email != tt.want.Email {
					t.Errorf("incorrect email: %s\n", got.Email)
				}
				if got.SQSMessageID != tt.want.SQSMessageID {
					t.Errorf("incorrect SQS Message ID: %s\n", got.SQSMessageID)
				}
			}
		})
	}
}

func TestAsDynamoInput(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		given     *User
		tableName string
		want      *dynamodb.PutItemInput
		shouldErr bool
	}{
		"valid user maps to input": {
			given: &User{
				ID:           "id",
				SQSMessageID: "sqsID",
				Email:        "email",
				Name:         "name",
				CreatedAt:    now,
			},
			tableName: "table",
			want: &dynamodb.PutItemInput{
				TableName: aws.String("table"),
				Item: map[string]types.AttributeValue{
					"ID":           &types.AttributeValueMemberS{Value: "id"},
					"Name":         &types.AttributeValueMemberS{Value: "name"},
					"Email":        &types.AttributeValueMemberS{Value: "email"},
					"SQSMessageID": &types.AttributeValueMemberS{Value: "sqsID"},
					"CreatedAt":    &types.AttributeValueMemberN{Value: strconv.Itoa(int(now.Unix()))},
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := tt.given.asDynamoInput(tt.tableName)
			if tt.shouldErr {
				if err == nil {
					t.Error("expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s\n", err)
				}
				if !isEqualPutItemInput(got, tt.want) {
					t.Logf("got: %#v\n", got)
					t.Logf("wanted: %#v\n", tt.want)
					t.Error("incorrect response")
				}
			}
		})
	}
}

func isEqualPutItemInput(a, b *dynamodb.PutItemInput) bool {
	isEqual := true
	isEqual = aws.ToString(a.TableName) == aws.ToString(b.TableName)
	for k, v := range a.Item {
		vv, ok := b.Item[k]
		if !ok {
			isEqual = false
		}
		isEqual = reflect.DeepEqual(v, vv)
	}
	return isEqual
}
