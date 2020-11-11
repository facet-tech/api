package db

import (
	"fmt"

	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var Database = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

func CreateId(key string) string {
	return key + "~" + util.GenerateBase64UUID()
}

const (
	WorkspaceTableName = "workspace-prod"
	FacetTableName     = "facet-prod"
)

func main() {
	// Create Dynamodb AWS session
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	res, err := callDynamodb(svc)
	if err != nil {
		fmt.Printf("Error returned %d", err)
	}
	fmt.Printf("Result %s", res.GoString())
}

func callDynamodb(svc dynamodbiface.DynamoDBAPI) (*dynamodb.GetItemOutput, error) {
	// Call GetItem
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("db-name"),
		Key: map[string]*dynamodb.AttributeValue{
			"Attribute": {
				S: aws.String("key"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
