package facet

import (
	"errors"
	"fmt"

	"facet.ninja/api/db"
	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DomElement struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Facet struct {
	Enabled    bool         `json:"enabled"`
	Name       string       `json:"name"`
	DomElement []DomElement `json:"domElement"`
	// facet should be visible throughout the domain navigation (regardless pathname)
	Global bool `json:"global"`
	// should remove facet from the DOM (as opposed to simple display: 'none')
	DomRemove bool `json:"domRemove"`
}

type FacetDTO struct {
	DomainId string  `json:"domainId"`
	UrlPath  string  `json:"urlPath"`
	Facet    []Facet `json:"facet"`
	Version  string  `json:"version"`
}

func FetchAll(siteId string) (*[]FacetDTO, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.FacetTableName),
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
	result, error := db.Database.Query(input)
	facets := new([]FacetDTO)
	if error == nil && result.Items != nil {
		error = dynamodbattribute.UnmarshalListOfMaps(result.Items, facets)
	}
	return facets, error
}

func (facet *FacetDTO) fetch() error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.FacetTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domainId": {
				S: aws.String(facet.DomainId),
			},
			"urlPath": {
				S: aws.String(facet.UrlPath),
			},
		},
	}
	result, error := db.Database.GetItem(input)
	if error == nil && result != nil {
		if len(result.Item) == 0 {
			error = errors.New(util.NOT_FOUND)
		} else {
			error = dynamodbattribute.UnmarshalMap(result.Item, facet)
		}
	}
	return error
}

func (facet *FacetDTO) create() error {
	fmt.Println("CHECKME",facet)
	item, error := dynamodbattribute.MarshalMap(facet)
	if error == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.FacetTableName),
			Item:      item,
		}
		_, error = db.Database.PutItem(input)
	}
	return error
}

func (facet *FacetDTO) delete() error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.FacetTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domainId": {
				S: aws.String(facet.DomainId),
			},
			"urlPath": {
				S: aws.String(facet.UrlPath),
			},
		},
	}
	_, error := db.Database.DeleteItem(input)
	return error
}
