package user

import (
	"errors"
	"facet/api/db"
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId"`
	Email       string                 `json:"email"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

type WorkspaceUser struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId,omitempty"`
	ApiKey      string                 `json:"apiKey"`
	User        string                 `json:"user"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

const (
	KEY_USER     = "USER"
	EMAIL_INDEX  = "email-index"
	APIKEY_INDEX = "apiKey-index"
)

func (user *User) fetch() error {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		IndexName: aws.String(EMAIL_INDEX),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(user.Email),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	if error == nil && result != nil {
		if len(result.Items) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Items[0], user)
		}
	}
	return error
}

func (user *User) Update() error {
	if len(user.Id) == 0 {
		user.Id = db.CreateRandomId(KEY_USER)
	}
	// TODO need to check if userId already exists ..
	apiKey := util.GenerateBase64UUID()
	userWorkspace := WorkspaceUser{
		Id:          user.Id + ":apiKey~" + apiKey,
		WorkspaceId: user.WorkspaceId,
		ApiKey:      apiKey,
		User:        user.Id,
		Attribute:   nil,
	}
	err := createUserWorkspace(&userWorkspace)
	if err != nil {
		return err
	}
	item, error := dynamodbattribute.MarshalMap(user)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.WorkspaceTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func createUserWorkspace(workspace *WorkspaceUser) error {
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

func (user *User) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.WorkspaceTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"workspaceId": {
				S: aws.String(user.WorkspaceId),
			},
			"id": {
				S: aws.String(user.Id),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
