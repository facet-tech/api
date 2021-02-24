package backend

import (
	"errors"
	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Framework struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Language struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Signature struct {
	Enabled    bool        `json:"enabled"`
	Name       string      `json:"name"`
	Signature  string      `json:"signature"`
	Parameter  []Parameter `json:"parameter"`
	ReturnType string      `json:"returnType"`
}

type DTO struct {
	AppId              string                 `json:"appId"`
	FullyQualifiedName string                 `json:"fullyQualifiedName"`
	Signature          []Signature            `json:"signature"`
	Version            string                 `json:"version"`
	Attribute          map[string]interface{} `json:"attribute,omitempty"`
	Language           Language               `json:"language"`
	Framework          []Framework            `json:"framework,omitempty"`
	InterfaceSignature []string               `json:"interfaceSignature,omitempty"`
	Type               string                 `json:"type,omitempty"`
	ParentSignature    string                 `json:"parentSignature,omitempty"`
}

func FetchAll(AppId string) (*[]DTO, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.BackendTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"appId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(AppId),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	facets := new([]DTO)
	if error == nil && result.Items != nil {
		error = dynamodbattribute.UnmarshalListOfMaps(result.Items, facets)
	}
	return facets, error
}

func (facet *DTO) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
			"fullyQualifiedName": {
				S: aws.String(facet.FullyQualifiedName),
			},
		},
	}
	result, error := db.Database.GetItem(input)
	if error == nil && result != nil {
		if len(result.Item) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Item, facet)
		}
	}
	return error
}

func (facet *DTO) create() error {
	item, error := dynamodbattribute.MarshalMap(facet)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.BackendTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (facet *DTO) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
			"fullyQualifiedName": {
				S: aws.String(facet.FullyQualifiedName),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}

func (facet *DTO) deleteAll() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.BackendTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appId": {
				S: aws.String(facet.AppId),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}
