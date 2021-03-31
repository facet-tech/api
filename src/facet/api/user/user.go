package user

import (
	"errors"
	"facet/api/db"
	"facet/api/util"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	ApiKey      string                 `json:"apiKey,omitempty"`
	User        string                 `json:"user"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
	Email       string                 `json:"email,omitempty"`
}

const (
	KEY_USER     = "USER"
	EMAIL_INDEX  = "email-index"
	APIKEY_INDEX = "apiKey-index"
)

func (user *User) fetch() (WorkspaceUser, error) {
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
	workspaceUser, _ := user.getWorkspaceUserByUserId()
	workspaceUser.Email = user.Email
	workspaceUser.Id = user.Id
	workspaceUser.WorkspaceId = user.WorkspaceId
	workspaceUser.Attribute = user.Attribute

	return workspaceUser, error
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
	//TODO Delete WorkspaceUser
	return err
}

func (user *User) getWorkspaceUserByUserId() (WorkspaceUser, error) {
	// condition represents the boolean condition of whether the item
	// attribute "CodeName" starts with the substring "Ben"
	condition := expression.Contains(expression.Name("id"), user.Id)

	//filt := expression.Name("Artist").Equal(expression.Value("No One You Know"))
	//proj := expression.NamesList(expression.Name("SongTitle"), expression.Name("AlbumTitle"))
	expr, err := expression.NewBuilder().WithFilter(condition).Build()
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(db.WorkspaceTableName),
		IndexName:                 aws.String(APIKEY_INDEX),
	}

	result, err := db.Database.Scan(input)
	workspaceUser := new([]WorkspaceUser)
	if err == nil && result != nil {
		if len(result.Items) == 0 {
			err = errors.New(util.NOT_FOUND)
			return WorkspaceUser{}, err
		} else {
			err = dynamodbattribute.UnmarshalListOfMaps(result.Items, workspaceUser)
		}
	}
	return (*workspaceUser)[0], err

}
