package db

import (
	"facet/api/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Database = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

func CreateRandomId(key string) string {
	return key + "~" + util.GenerateBase64UUID()
}

func CreateId(key string, id string) string {
	return key + "~" + id
}

func CreateKey(name string) string {
	return name + "~"
}

const (
	WorkspaceTableName     = "workspace-prod"
	FacetTableName         = "facet-prod"
	BackendTableName       = "facet-backend-prod"
	ConfigurationTableName = "facet-configuration"
)
