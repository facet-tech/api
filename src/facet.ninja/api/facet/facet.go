package facet

import (
	"facet.ninja/api/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const FACET_TABLE_NAME = "facet"

type DomElement struct {
	Enabled string   `json:"enabled"`
	Path    []string `json:"path"`
}

type Facet struct {
	DomainId   string       `json:"domainId"`
	UrlPath    string       `json:"urlPath"`
	DomElement []DomElement `json:"domElement"`
}

func FetchAll(siteId string) (*[]Facet, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(FACET_TABLE_NAME),

		KeyConditions: map[string]*dynamodb.Condition{
			"domainId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(siteId),
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
	facets := new([]Facet)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, facets)

	if err != nil {
		return nil, err
	}

	return facets, nil
}

func Fetch(domainId string, urlPath string) (*Facet, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(FACET_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"domainId": {
				S: aws.String(domainId),
			},
			"urlPath": {
				S: aws.String(urlPath),
			},
		},
	}
	result, err := db.Database.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}
	facets := new(Facet)
	err = dynamodbattribute.UnmarshalMap(result.Item, facets)

	if err != nil {
		return nil, err
	}

	return facets, nil
}

func Put(facet Facet) error {
	log.Print(facet)
	item, err := dynamodbattribute.MarshalMap(facet)
	log.Print(item)
	log.Print(err)
	if err == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(FACET_TABLE_NAME),
			Item:      item,
		}
		_, err2 := db.Database.PutItem(input)
		log.Print(err2)
		return err2
	} else {
		return err
	}
}

func Delete(facet Facet) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(FACET_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"domainId": {
				S: aws.String(facet.DomainId),
			},
			"urlPath": {
				S: aws.String(facet.UrlPath),
			},
		},
	}
	_, err := db.Database.DeleteItem(input)
	return err
}
