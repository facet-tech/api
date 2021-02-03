package main

import (
	"context"
	"facet/api/domain"
	"facet/api/facet"
	"facet/api/middleware"
	"facet/api/notification"
	"facet/api/user"
	"facet/api/util"
	"facet/api/workspace"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		log.Printf("Gin cold start")
		router := gin.Default()
		defaultRoutes(router)
		router.Group("/")
		{
			router.GET("/js", getJs)
			router.GET("/js/facetmap", getFacetMap)
			notification.Route(router)
		}

		// authenticated routes
		router.Group("/")
		{
			router.Use(middleware.JWTVerify())
			facet.Route(router)
			workspace.Route(router)
			domain.Route(router)
			user.Route(router)
		}
		ginLambda = ginadapter.New(router)
	}
	return ginLambda.ProxyWithContext(ctx, req)
}

func defaultRoutes(route *gin.Engine) {
	route.OPTIONS("/*anyPath", util.Options)
}

var mutationObserverTemplate *template.Template
var moFile []byte

func getFacetMap(c *gin.Context) {
	util.SetCorsHeaders(c)
	domainId := c.Request.URL.Query().Get("id")
	facetMapJsonString, err := facet.ComputeMutationObserverFacetMap(domainId)
	util.SetResponseCode(facetMapJsonString, err, c)
}

func getJs(c *gin.Context) {
	util.SetCorsHeaders(c)

	site := c.Request.URL.Query().Get("id")
	if &site == nil {
		c.JSON(http.StatusNotFound, "id is required")
		return
	}

	var err error
	if mutationObserverTemplate == nil {
		moFile, err = ioutil.ReadFile("./resources/templates/mutationObserver.js")
		if err != nil {
			log.Print(err)
		}
	}

	c.Data(http.StatusOK, "text/javascript", moFile)

}
