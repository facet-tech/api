package middleware

import (
	"facet/api/db"
	"facet/api/user"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func APIKeyVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: cache valid JWTs https://github.com/facets-io/api/issues/14
		skipAuthentication, _ := strconv.ParseBool(os.Getenv("SKIP_AUTHENTICATION"))
		if skipAuthentication {
			return
		}
		apiKey := c.GetHeader("ApiKey")
		containsApiKey, _ := fetchApiKey(apiKey)
		if !containsApiKey {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// TODO: verify claims against the Cognito Pool https://github.com/facets-io/api/issues/13
	}
}

func fetchApiKey(apiKey string) (bool, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(db.WorkspaceTableName),
		IndexName: aws.String(user.APIKEY_INDEX),
		KeyConditions: map[string]*dynamodb.Condition{
			"apiKey": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(apiKey),
					},
				},
			},
		},
	}
	result, error := db.Database.Query(input)

	if result == nil || result.Count == nil || *result.Count == 0 {
		return false, error
	}
	return true, error
}
