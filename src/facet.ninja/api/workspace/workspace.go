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
	KEY_WORKSPACE = "WORKSPACE"
)

func (workspace *Workspace) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.WorkspaceTableName),
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
			TableName: aws.String(db.WorkspaceTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}


func (workspace *Workspace) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.WorkspaceTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(workspace.Id),
			},
			"id": {
				S: aws.String(workspace.Id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
