package db

import (
	"facet.ninja/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Database = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

func CreateId(key string) string {
	return key + "~" + util.GenerateBase64UUID()
}
