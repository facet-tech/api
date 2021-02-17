package workspace

import (
	"errors"
	"facet/api/db"
	"facet/api/pricing"
	"facet/api/util"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type WorkspaceDto struct {
	DomainId    string `json:"domainId"`
	WorkspaceId string `json:"workspaceId,omitempty"`
	Counter     int    `json:"counter"`
}

type Workspace struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId,omitempty"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

const (
	KEY_WORKSPACE = "WORKSPACE"
)

func (workspace *Workspace) fetchAll() (*[]Workspace, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"workspaceId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(workspace.WorkspaceId),
					},
				},
			},
			"id": {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("DOMAIN~"),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	workspaces := new([]Workspace)
	error = dynamodbattribute.UnmarshalListOfMaps(result.Items, workspaces)

	var resultArr []WorkspaceDto

	for _, workspace := range *workspaces {
		fmt.Println("eeee", workspace)
		element, _ := pricing.Fetch(workspace.WorkspaceId)
		workspaceDto := WorkspaceDto{
			WorkspaceId: workspace.WorkspaceId,
			DomainId:    element.DomainId,
			Counter:     0, //TODO count the elements
		}
		err := append(resultArr, workspaceDto)
		if err == nil && result != nil {
			if len(result.Items) == 0 {
				error = errors.New(util.NOT_FOUND)
			} else {
				error = dynamodbattribute.UnmarshalMap(result.Items[0], workspace)
			}
		}
	}
	return workspaces, error

}

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
