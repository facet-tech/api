package backend

import (
	"encoding/json"
	"facet/api/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	BaseUrl                          = "/facet/backend"
	AppIdQueryParameter              = "appId"
	FullyQualifiedNameQueryParamater = "fullyQualifiedName"
)

func Route(router *gin.Engine) {
	router.GET(BaseUrl, Get)
	router.POST(BaseUrl, Post)
	router.DELETE(BaseUrl, Delete)
}

func Get(c *gin.Context) {
	facet := DTO{}
	facet.AppId = c.Request.URL.Query().Get(AppIdQueryParameter)
	facet.FullyQualifiedName = c.Request.URL.Query().Get(FullyQualifiedNameQueryParamater)
	var facetArray *[]DTO
	var error error
	if facet.FullyQualifiedName == "" {
		facetArray, error = FetchAll(facet.AppId)
	} else {
		error = facet.fetch()
		facetArray = &[]DTO{facet}
	}
	util.SetResponseCode(facetArray, error, c)
}

func Post(c *gin.Context) {
	facet := DTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	error = facet.create()
	util.SetResponseCode(facet, error, c)
}

func Delete(c *gin.Context) {
	facet := DTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	if facet.FullyQualifiedName == "" {
		facetArray, _ := FetchAll(facet.AppId)
		for _, facetDto := range *facetArray {
			error = facetDto.delete()
		}
	} else {
		error = facet.delete()
	}

	util.SetResponseCode(nil, error, c)
}
