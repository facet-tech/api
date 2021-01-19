package main

import (
	"bytes"
	"context"
	"encoding/json"
	"facet.ninja/api/middleware"
	"log"
	"net/http"
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

type FacetMapValue struct {
	Path      string `json:"path"`
	DomRemove bool   `json:"domRemove"`
}

func computeFacetMap(site string) (string, error) {
	globalFacetKey := "GLOBAL-FACET-DECLARATION"
	facetMap := map[string][]FacetMapValue{}
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
					element := FacetMapValue{
						Path:      domElement.Path,
						DomRemove: facetElement.DomRemove,
					}
					if facetElement.Global {
						facetMap[globalFacetKey] = append(facetMap[globalFacetKey], element)
					} else {
						facetMap[facetDto.UrlPath] = append(facetMap[facetDto.UrlPath], element)
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
