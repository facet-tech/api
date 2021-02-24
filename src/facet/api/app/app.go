package app

import (
	"errors"

	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type App struct {
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
	Name        string                 `json:"name"`
	Id          string                 `json:"id"`
	Environment string                 `json:"environment"`
	WorkspaceId string                 `json:"workspaceId"`
}

const (
	keyApp = "APP"
)

func (app *App) create() error {
	app.Id = db.CreateId(keyApp,app.Name)
	item, error := dynamodbattribute.MarshalMap(app)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.WorkspaceTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (app *App) fetch() error {
	if app.Id == "" {
		app.Id = db.CreateId(keyApp,app.Name)
	}
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"workspaceId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(app.WorkspaceId),
					},
				},
			},
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(app.Id),
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
			error = dynamodbattribute.UnmarshalMap(result.Items[0], app)
		}
	}
	return error
}

func (app *App) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.WorkspaceTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(app.WorkspaceId),
			},
			"id": {
				S: aws.String(app.Id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
