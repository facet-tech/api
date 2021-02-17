package pricing

import (
	"facet/api/db"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"time"
)

type Pricing struct {
	RequestId string `json:"requestId"`
	DomainId  string `json:"domainId"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"userAgent,omitempty"`
}

func (pricing *Pricing) Create() error {
	time.Sleep(15000)
	pricing.RequestId = db.CreateId(db.PricingTableName)
	pricing.Timestamp = time.Now().UTC().Format(time.RFC3339)
	item, e := dynamodbattribute.MarshalMap(pricing)
	fmt.Println("ela", item)
	if e == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.PricingTableName),
			Item:      item,
		}
		_, e = db.Database.PutItem(input)
	}
	return e
}

func FetchAll(domainId string) (*[]Pricing, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.FacetTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"domainId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(domainId),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)
	pricingRecords := new([]Pricing)
	if error == nil && result.Items != nil {
		error = dynamodbattribute.UnmarshalListOfMaps(result.Items, pricingRecords)
	}
	return pricingRecords, error
}

func Fetch(workspaceId string) (*Pricing, error) {
	pricingElement := new(Pricing)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.PricingTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"domainId": {
				S: aws.String(workspaceId),
			},
		},
	}
	result, err := db.Database.GetItem(input)
	if err == nil && result != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, pricingElement)
	}

	return pricingElement, err
}
