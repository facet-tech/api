package pricing

import (
	"facet/api/db"
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
	pricing.RequestId = db.CreateId(db.PricingTableName)
	pricing.Timestamp = time.Now().UTC().Format(time.RFC3339)
	item, e := dynamodbattribute.MarshalMap(pricing)
	if e == nil {
		input := &dynamodb.PutItemInput{
			TableName: aws.String(db.PricingTableName),
			Item:      item,
		}
		_, e = db.Database.PutItem(input)
	}
	return e
}

func Count(domainId string) (*int64, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.PricingTableName),
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
	result, err := db.Database.Query(input)
	return result.ScannedCount, err
}
