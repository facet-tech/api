package domain

import (
	"errors"

	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Domain struct {
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
	Domain      string                 `json:"domain"`
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId"`
}

const (
	KEY_DOMAIN      = "DOMAIN"
	DOMAIN_ID_INDEX = "domain-workspaceId-index"
)

func (domain *Domain) create() error {
	domain.Id = db.CreateRandomId(KEY_DOMAIN)
	item, error := dynamodbattribute.MarshalMap(domain)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.WorkspaceTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func FetchAll(workspaceId string) (*[]Domain, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"workspaceId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(workspaceId),
					},
				},
			},
			"id": {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(db.CreateKey(KEY_DOMAIN)),
					},
				},
			},

		},
	}
	result, error := db.Database.Query(input)
	array := new([]Domain)
	if error == nil && result.Items != nil {
		error = dynamodbattribute.UnmarshalListOfMaps(result.Items, array)
	}
	return array, error
}

func (domain *Domain) fetch() error {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		IndexName: aws.String(DOMAIN_ID_INDEX),
		KeyConditions: map[string]*dynamodb.Condition{
			"workspaceId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(domain.WorkspaceId),
					},
				},
			},
			"domain": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(domain.Domain),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	if error == nil && result.Items != nil {
		if len(result.Items) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Items[0], domain)
		}
	}
	return error
}

func (domain *Domain) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.WorkspaceTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(domain.WorkspaceId),
			},
			"id": {
				S: aws.String(domain.Id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
