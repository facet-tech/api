package main

import (
	"bytes"
	"context"
	"facet/api/domain"
	"facet/api/facet"
	"facet/api/middleware"
	"facet/api/notification"
	"facet/api/pricing"
	"facet/api/user"
	"facet/api/util"
	"facet/api/workspace"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
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
	if mutationObserverTemplate == nil {
		t, err := template.ParseFiles("./resources/templates/mutationObserver.js")
		if err != nil {
			log.Print(err)
		}
		mutationObserverTemplate = t
	}

	domainId := c.Request.URL.Query().Get("id")
	if &domainId == nil {
		c.JSON(http.StatusNotFound, "id is required")
		return
	}

	facetMap, error := facet.ComputeMutationObserverFacetMap(domainId)
	config := map[string]string{
		"GO_ARRAY_REPLACE_ME": facetMap,
	}

	var tpl bytes.Buffer
	if err := mutationObserverTemplate.Execute(&tpl, config); err != nil {
		log.Print(err)
	}

	result := tpl.String()

	if error == nil {
		ua := c.Request.Header.Get("User-Agent")
		pricingElement := pricing.Pricing{
			DomainId:  domainId,
			UserAgent: ua,
		}
		go pricingElement.Create()
		c.Data(http.StatusOK, "text/javascript", []byte(result))

	} else {
		c.JSON(http.StatusInternalServerError, error)
	}
}
