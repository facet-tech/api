package workspace

import (
	"facet.ninja/api/db"
	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Workspace struct {
	Id        string                 `json:"id"`
	Key       string                 `json:"key"`
	Domain    string                 `json:"domain,omitempty"`
	Email     string                 `json:"email,omitempty"`
	Attribute map[string]interface{} `json:"attribute,omitempty"`
}

const (
	KEY_WORKSPACE        = "WORKSPACE"
	WORKSPACE_TABLE_NAME = "workspace"
)

func createKey(key string, value string) string {
	return key + "#" + value
}

func Fetch(id string) (*Workspace, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(WORKSPACE_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"key": {
				S: aws.String(KEY_WORKSPACE),
			},
		},
	}
	item, err := db.Database.GetItem(input)
	result := new(Workspace)

	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalMap(item.Item, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Create(workspace Workspace) (result Workspace, error error) {
	workspace.Id = util.GenerateBase64UUID()
	workspace.Key = KEY_WORKSPACE
	item, err := dynamodbattribute.MarshalMap(workspace)
	if err == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(WORKSPACE_TABLE_NAME),
			Item:      item,
		}
		_, err2 := db.Database.PutItem(input)
		log.Print(err2)
		return workspace, err2
	} else {
		log.Print(err)
		return workspace, err
	}
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
}

func AssosiateUser(userSiteMappng UserSiteMapping) error {
	item, err := dynamodbattribute.MarshalMap(userSiteMappng)

	if err == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(USER_SITE_MAPPING),
			Item:      item,
		}
		_, err2 := db.Database.PutItem(input)
		return err2
	} else {
		return err
	}
}

func FetchUserMapping(userId string) (*UserSiteMapping, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(USER_SITE_MAPPING),
		KeyConditions: map[string]*dynamodb.Condition{
			"userId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList:     []*dynamodb.AttributeValue {
					{
						S: aws.String(userId),
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
	userSiteMapping := new(UserSiteMapping)
	err = dynamodbattribute.UnmarshalMap(result.Items, userSiteMapping)

	if err != nil {
		return nil, err
	}

	return userSiteMapping, nil
}
*/
