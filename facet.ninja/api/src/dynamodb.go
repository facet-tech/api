package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

type Facet struct {
	Name    string   `json:"name"`
	Enabled string   `json:"enabled"`
	Id      []string `json:"id"`
}

type Facets struct {
	Site  string  `json:"site"`
	Facet []Facet `json:"facet"`
}

func getItem(site string) (*Facets, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("facet"),
		Key: map[string]*dynamodb.AttributeValue{
			"site": {
				S: aws.String(site),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}
	facets := new(Facets)
	err = dynamodbattribute.UnmarshalMap(result.Item, facets)

	if err != nil {
		return nil, err
	}

	return facets, nil
}

func putItem(site string, facets Facets) error {
	item, err := dynamodbattribute.MarshalMap(facets)

	if err == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String("facet"),
			Item:      item,
		}
		_, err2 := db.PutItem(input)
		return err2
	} else {
		return err
	}
}

func deleteItem(site string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("facet"),
		Key: map[string]*dynamodb.AttributeValue{
			"site": {
				S: aws.String(site),
			},
		},
	}
	_, err := db.DeleteItem(input)
	return err
}
