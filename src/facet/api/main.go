package main

import (
	"context"
	"encoding/json"
	"facet/api/domain"
	"facet/api/facet"
	"facet/api/middleware"
	"facet/api/notification"
	"facet/api/user"
	"facet/api/util"
	"facet/api/workspace"
	"fmt"
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
			router.GET("/js/computefacetmap", computeFacetMap)
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

// TODO refactor and call this method through facet entity
func computeFacetMap(c *gin.Context) {
	util.SetCorsHeaders(c)
	fmt.Println("MPIADWADWWA")
	site := c.Request.URL.Query().Get("id")
	fmt.Println("SITE",site)
	globalFacetKey := "GLOBAL-FACET-DECLARATION"
	facetMap := map[string][]string{}
	if &site != nil {
		facets, errFetch := facet.FetchAll(site)
		fmt.Println("ela man",facets)
		if errFetch != nil {
			return
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
	fmt.Println(facetMap)
	facetMapJSON, _ := json.Marshal(facetMap)
	facetMapJSONString := string(facetMapJSON)
	c.Data(http.StatusOK, "text/javascript", []byte(facetMapJSONString))
}

var mutationObserverTemplate *template.Template
var moFile []byte

func getJs(c *gin.Context) {
	util.SetCorsHeaders(c)

	fmt.Println("mMPIKAAA")
	site := c.Request.URL.Query().Get("id")
	fmt.Println("mMPIKAAA", site)
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
