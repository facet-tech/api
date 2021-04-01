package configuration

import (
	"encoding/json"
	"facet/api/middleware"
	"facet/api/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	BaseUrl           = "/facet/configuration"
	PropertyParameter = "property"
	IdParameter       = "id"
)

func Route(router *gin.Engine) {
	router.GET(BaseUrl, middleware.APIKeyVerify(), Get)
	router.POST(BaseUrl, middleware.APIKeyVerify(), Post)
	router.DELETE(BaseUrl, middleware.APIKeyVerify(), Delete)
}

func Get(c *gin.Context) {
	configuration := Configuration{}
	configuration.Property = c.Request.URL.Query().Get(PropertyParameter)
	configuration.Id = c.Request.URL.Query().Get(IdParameter)
	error := configuration.fetch()
	util.SetResponseCode(configuration, error, c)
}

func Post(c *gin.Context) {
	configuration := Configuration{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &configuration)
	error = configuration.create()
	util.SetResponseCode(configuration, error, c)
}

func Delete(c *gin.Context) {
	configuration := Configuration{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &configuration)
	error = configuration.delete()
	util.SetResponseCode(nil, error, c)
}
