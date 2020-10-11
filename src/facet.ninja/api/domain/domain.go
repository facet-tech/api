package domain

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

type Domain struct {
	Id          string                 `json:"id"`
	Domain      string                 `json:"domain"`
	WorkspaceId string                 `json:"workspaceId"`
	Attribute   map[string]interface{} `json:"attribute,omitempty"`
}

type Id struct {
	Id string `json:"id"`
}

const (
	KEY_DOMAIN           = "DOMAIN"
	WORKSPACE_TABLE_NAME = "workspace"
)

func createKey(key string, value string) string {
	return key + "#" + value
}

func getValueFromKey(key string, prefix string) string {
	return strings.TrimPrefix(key, prefix+"#")
}

func Create(domain Domain) (result Domain, error error) {
	workspace := workspace.Workspace{}
	workspace.Id = domain.WorkspaceId
	domain.Id = util.GenerateBase64UUID()
	workspace.Key = createKey(KEY_DOMAIN, domain.Id)
	workspace.Domain = domain.Domain
	workspace.Attribute = domain.Attribute
	log.Print(domain)
	item, err := dynamodbattribute.MarshalMap(workspace)
	log.Print(item)
	log.Println(err)
	if err == nil {

		input := &dynamodb.PutItemInput{
			TableName: aws.String(WORKSPACE_TABLE_NAME),
			Item:      item,
		}
		_, err2 := db.Database.PutItem(input)
		log.Println(err2)
		return domain, err2
	} else {
		return domain, err
	}
}

func Fetch(domain string, workspaceId string) (*Id, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(WORKSPACE_TABLE_NAME),
		IndexName: aws.String("domain-id-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(workspaceId),
					},
				},
			},
			"domain": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(domain),
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
	id.Id = getValueFromKey(workspace.Key, KEY_DOMAIN)
	if err != nil {
		return nil, err
	}

	return id, nil
}
