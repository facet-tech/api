package db

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Define stub
type stubDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *stubDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// Make response
	key := dynamodb.AttributeValue{}
	key.SetS("sample-key")
	val := dynamodb.AttributeValue{}
	val.SetS("sample-val")
	resp := make(map[string]*dynamodb.AttributeValue)
	resp["key"] = &key
	resp["val"] = &val

	// Returned canned response
	output := &dynamodb.GetItemOutput{
		Item: resp,
	}
	return output, nil
}

// Sample Test Case
func TestDynamodb(t *testing.T) {
	svc := &stubDynamoDB{}
	res, err := callDynamodb(svc)
	if err != nil {
		t.Errorf("Error calling Dynamodb %d", err)
	}
	keyRes := *res.Item["key"].S
	valRes := *res.Item["val"].S
	if keyRes != "sample-key" {
		t.Errorf("Wrong key returned. Shoule be sample-key, was %s", keyRes)
	}
	if valRes != "sample-val" {
		t.Errorf("Wrong value returned. Shoule be sample-val, was %s", valRes)
	}
}
