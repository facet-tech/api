package main

import (
	"bytes"
	"context"
	"facet.ninja/api/middleware"
	"log"
	"strings"
	"text/template"

	"facet.ninja/api/domain"
	"facet.ninja/api/facet"
	"facet.ninja/api/user"
	"facet.ninja/api/util"
	"facet.ninja/api/workspace"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
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
		router.GET("/facet.ninja.js", getJs)
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

func getJs(c *gin.Context) {
	util.SetCorsHeaders(c)
	if mutationObserverTemplate == nil {
		t, err := template.ParseFiles("./resources/templates/mutationObserver.js") // Parse template file.
		if err != nil {
			log.Print(err)
		}
		mutationObserverTemplate = t
	}

	var commaSeperatedIdsString string
	site := c.Request.URL.Query().Get("id")
	if &site != nil {
		facets, error := facet.FetchAll(site)
		for _, facetDto := range *facets {
			commaSeperatedIdsString += "\t['" + facetDto.UrlPath + "',new Set(["
			for _, facet := range facetDto.Facet {
				for _, domElement := range facet.DomElement {
					commaSeperatedIdsString += "'" + domElement.Path + "',"
				}
			}
			commaSeperatedIdsString = strings.TrimSuffix(commaSeperatedIdsString, ",")
			commaSeperatedIdsString += "])],\n"
		}

		wantedString := strings.TrimSuffix(commaSeperatedIdsString, ",\n")

		config := map[string]string{
			"GO_ARRAY_REPLACE_ME": wantedString,
		}

		var tpl bytes.Buffer
		if err := mutationObserverTemplate.Execute(&tpl, config); err != nil {
			log.Print(err)
		}

		result := tpl.String()

		if error == nil {
			c.Data(200, "text/javascript", []byte(result))
		} else {
			c.JSON(500, error)
		}
	} else {
		c.JSON(400, "id is required")
	}
}
