package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		/*r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "*"
			},
		}))*/

		r.GET("/facet/:id", getFacet)
		r.POST("/facet/:id", createFacet)
		r.DELETE("/facet/:id", deleteFacet)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

func getFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	facets, error := getItem(site)

	if error == nil {
		c.JSON(200, facets)
	} else {
		c.JSON(500, error)
	}

}

func deleteFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	error := deleteItem(site)

	if error == nil {
		c.JSON(200, "Site "+site+"deleted.")
	} else {
		c.JSON(500, error)
	}

}

func createFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	newFacet := Facets{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &newFacet)
	error = putItem(site, newFacet)

	if error == nil {
		c.JSON(201, "Created facet: "+site)
	} else {
		c.JSON(500, error)
	}

}
