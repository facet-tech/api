package facet

import (
	"encoding/json"
	"facet.ninja/api/middleware"
	"io/ioutil"

	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL                  = "/facet"
	DOMAIN_ID_QUERY_PARAMATER = "domainId"
	URL_PATH_QUERY_PARAMATER  = "urlPath"
)

func AuthenticatedRoute(router *gin.Engine) {
	router.GET(BASE_URL, middleware.JWTVerify(), Get)
	router.POST(BASE_URL, middleware.JWTVerify(), Post)
	router.DELETE(BASE_URL, middleware.JWTVerify(), Delete)
}

func Get(c *gin.Context) {
	facet := FacetDTO{}
	facet.DomainId = c.Request.URL.Query().Get(DOMAIN_ID_QUERY_PARAMATER)
	facet.UrlPath = c.Request.URL.Query().Get(URL_PATH_QUERY_PARAMATER)
	error := facet.fetch()
	util.SetResponseCode(facet, error, c)
}

func Post(c *gin.Context) {
	facet := FacetDTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	error = facet.create()
	util.SetResponseCode(facet, error, c)
}

func Delete(c *gin.Context) {
	facet := FacetDTO{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &facet)
	error = facet.delete()
	util.SetResponseCode(nil, error, c)
}
