package workspace

import (
	"errors"
	"facet/api/db"
	"facet/api/domain"
	"facet/api/pricing"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type WorkspaceDto struct {
	DomainId string `json:"domainId"`
	Counter  *int64 `json:"counter"`
	Domain   string `json:"Domain"`
}

type WorkspaceEntityDto struct {
	WorkspaceId string         `json:"workspaceId,omitempty"`
	Workspaces  []WorkspaceDto `json:"workspaces,omitempty"`
}

type Workspace struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId,omitempty"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

const (
	KEY_WORKSPACE = "WORKSPACE"
)

func (workspace *Workspace) fetchAll() (WorkspaceEntityDto, error) {
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

	result, err := db.Database.Query(input)
	workspaces := new([]Workspace)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, workspaces)
	resultArr := make([]WorkspaceDto, 0)

	// TODO This needs to be decoupled/paginated. Danger of OOT exception
	for _, workspace := range *workspaces {
		domainElement, _ := domain.Fetch(workspace.WorkspaceId, workspace.Id)
		counter, _ := pricing.Count(workspace.Id)
		workspaceDto := WorkspaceDto{
			DomainId: workspace.Id,
			Counter:  counter,
			Domain:   domainElement.Domain,
		}
		resultArr = append(resultArr, workspaceDto)
	}
	resultObj := WorkspaceEntityDto{WorkspaceId: workspace.WorkspaceId, Workspaces: resultArr}
	return resultObj, err

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
