package user

import (
	"facet.ninja/api/db"
	"facet.ninja/api/util"
	"facet.ninja/api/workspace"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strings"
)

type User struct {
	Id          string                 `json:"id"`
	WorkspaceId string                 `json:"workspaceId"`
	Email       string                 `json:"email,omitempty"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

type Id struct {
	Id          string `json:"id"`
	WorkspaceId string `json:"workspaceId"`
}

const USER_TABLE_NAME = "workspace"

const (
	KEY_WORKSPACE        = "WORKSPACE"
	KEY_USER             = "USER"
	KEY_DOMAIN           = "DOMAIN"
	WORKSPACE_TABLE_NAME = "workspace"
)

func createKey(key string, value string) string {
	return key + "#" + value
}

func getValueFromKey(key string, prefix string) string {
	return strings.TrimPrefix(key, prefix+"#")
}

func Create(user User) (result User, error error) {
	workspace := workspace.Workspace{}
	workspace.Id = user.WorkspaceId
	user.Id = util.GenerateBase64UUID()
	workspace.Key = createKey(KEY_USER, user.Id)
	workspace.Email = user.Email
	workspace.Attribute = user.Attribute
	item, err := dynamodbattribute.MarshalMap(workspace)
	if err == nil {

		input := &dynamodb.PutItemInput{
			TableName: aws.String(USER_TABLE_NAME),
			Item:      item,
		}
		_, err2 := db.Database.PutItem(input)

		return user, err2
	} else {
		return user, err
	}
}

func Fetch(email string) (*Id, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(USER_TABLE_NAME),
		IndexName: aws.String("email-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
		},
	}
	result, err := db.Database.Query(input)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, nil
	}
	workspace := new(workspace.Workspace)
	err = dynamodbattribute.UnmarshalMap(result.Items[0], workspace)
	id := new(Id)
	id.Id = getValueFromKey(workspace.Key, KEY_USER)
	id.WorkspaceId = workspace.Id
	if err != nil {
		return nil, err
	}

	return id, nil
}

func Delete(user User) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(user.WorkspaceId),
			},
			"key": {
				S: aws.String(createKey(KEY_USER, user.Id)),
			},
		},
	}
	log.Println(user)
	log.Println(input)
	_, err := db.Database.DeleteItem(input)
	return err
}
