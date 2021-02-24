package configuration

import (
	"errors"
	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Configuration struct {
	Property  string                 `json:"property"`
	Id        string                 `json:"id"`
	Attribute map[string]interface{} `json:"attribute,omitempty"`
}

func (configuration *Configuration) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.ConfigurationTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"property": {
				S: aws.String(configuration.Property),
			},
			"id": {
				S: aws.String(configuration.Id),
			},
		},
	}
	result, error := db.Database.GetItem(input)
	if error == nil && result != nil {
		if len(result.Item) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Item, configuration)
		}
	}
	return error
}

func (configuration *Configuration) create() error {
	item, error := dynamodbattribute.MarshalMap(configuration)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.ConfigurationTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (configuration *Configuration) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.ConfigurationTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"property": {
				S: aws.String(configuration.Property),
			},
			"id": {
				S: aws.String(configuration.Id),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}
