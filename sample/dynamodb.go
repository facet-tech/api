package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func test() {
	error := deleteItem("facet.ninja")

	testFacet := Facets{
		Site: "facet.ninja",
		Facet: []Facet{
			Facet{
				Name:    "sword",
				Enabled: "true",
				Id:      []string{"andy", "cam", "john", "mene"}},
		},
	}
	error = putItem("facet.ninja", testFacet)

	if error != nil {
		fmt.Println(error)
	}

	testFacet2 := Facets{
		Site: "facet.ninja",
		Facet: []Facet{
			Facet{
				Name:    "sword2",
				Enabled: "false",
				Id:      []string{"andy", "cam", "john", "mene"}},
		},
	}
	error = updateItem("facet.ninja", testFacet2)

	if error != nil {
		fmt.Println(error)
	}
	facets, error2 := getItem("mywebsite.facets")
	fmt.Println(facets)
	if error2 != nil {
		fmt.Println(error2)
	}
}

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

func updateItem(site string, facets Facets) error {
	/*item, err := dynamodbattribute.MarshalMap(facets)

	if err == nil {
		input := &dynamodb.UpdateItemInput{
			TableName: aws.String("facet"),
			Key: map[string]*dynamodb.AttributeValue{
				"site": {
					S: aws.String(site),
				},
			},
			UpdateExpression:          aws.String("Set facet = :facet"),
			ExpressionAttributeNames
			ExpressionAttributeValues: item,
		}
		_, err2 := db.UpdateItem(input)
		return err2
	} else {
		return err
	}*/
	deleteItem(site)
	error := putItem(site, facets)
	return error
}
