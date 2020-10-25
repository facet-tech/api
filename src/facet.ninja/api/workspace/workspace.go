package workspace

import (
	"errors"
	"facet.ninja/api/db"
	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Workspace struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId,omitempty"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

const (
	KEY_WORKSPACE        = "WORKSPACE"
	WORKSPACE_TABLE_NAME = "workspace-temp"
)

func (workspace *Workspace) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(WORKSPACE_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(workspace.Id),
			},
			"id": {
				S: aws.String(workspace.Id),
			},
		},
	}
	result, error := db.Database.GetItem(input)
	if error == nil && result != nil {
		if len(result.Item) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Item, workspace)
		}
	}
	return error
}

func (workspace *Workspace) create() error {
	workspace.Id = db.CreateId(KEY_WORKSPACE)
	workspace.WorkspaceId = workspace.Id
	item, error := dynamodbattribute.MarshalMap(workspace)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(WORKSPACE_TABLE_NAME),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

/*
func Delete(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(SITE_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}*/
