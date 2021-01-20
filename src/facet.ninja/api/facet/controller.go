package facet

import (
	"encoding/json"
	"fmt"

	"io/ioutil"

	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL                  = "/facet"
	DOMAIN_ID_QUERY_PARAMATER = "domainId"
	URL_PATH_QUERY_PARAMATER  = "urlPath"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	router.DELETE(BASE_URL, Delete)
}

func Get(c *gin.Context) {
	facet := FacetDTO{}
	facet.DomainId = c.Request.URL.Query().Get(DOMAIN_ID_QUERY_PARAMATER)
	facet.UrlPath = c.Request.URL.Query().Get(URL_PATH_QUERY_PARAMATER)
	var facetArray *[]FacetDTO
	var error error
	if facet.UrlPath == "" {
		facetArray, error = FetchAll(facet.DomainId)
	} else {
		error = facet.fetch()
		facetArray = &[]FacetDTO{facet}
	}
	util.SetResponseCode(facetArray, error, c)
}

func Post(c *gin.Context) {
	facet := FacetDTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	error = facet.create()
	fmt.Println("ERR",error)
	util.SetResponseCode(facet, error, c)
}

func Delete(c *gin.Context) {
	facet := FacetDTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	error = facet.delete()
	util.SetResponseCode(nil, error, c)
}
