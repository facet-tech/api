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
