package main

import (
	"bytes"
	"context"
	"encoding/json"
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

func computeFacetMap(site string) (string, error) {
	globalFacetKey := "GLOBAL-FACET-DECLARATION"
	facetMap := map[string][]string{}
	var err error
	if &site != nil {
		facets, errFetch := facet.FetchAll(site)
		if errFetch != nil {
			return "{}", errFetch
		}
		for _, facetDto := range *facets {
			for _, facetElement := range facetDto.Facet {
				if facetElement.Enabled == false {
					continue
				}
				for _, domElement := range facetElement.DomElement {
					if facetElement.Global {
						facetMap[globalFacetKey] = append(facetMap[globalFacetKey], domElement.Path)
					} else {
						facetMap[facetDto.UrlPath] = append(facetMap[facetDto.UrlPath], domElement.Path)
					}
				}
			}
		}
	}
	facetMapJSON, _ := json.Marshal(facetMap)
	facetMapJSONString := string(facetMapJSON)
	return facetMapJSONString, err
}

var mutationObserverTemplate *template.Template

func getJs(c *gin.Context) {

	util.SetCorsHeaders(c)
	if mutationObserverTemplate == nil {
		t, err := template.ParseFiles("./resources/templates/mutationObserver.js")
		if err != nil {
			log.Print(err)
		}
		mutationObserverTemplate = t
	}

	site := c.Request.URL.Query().Get("id")
	if &site == nil {
		c.JSON(http.StatusNotFound, "id is required")
		return
	}

	facetMap, error := computeFacetMap(site)
	config := map[string]string{
		"GO_ARRAY_REPLACE_ME": facetMap,
	}

	var tpl bytes.Buffer
	if err := mutationObserverTemplate.Execute(&tpl, config); err != nil {
		log.Print(err)
	}

	result := tpl.String()

	if error == nil {
		c.Data(http.StatusOK, "text/javascript", []byte(result))
	} else {
		c.JSON(http.StatusInternalServerError, error)
	}
}
