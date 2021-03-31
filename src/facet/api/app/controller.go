package app

import (
	"encoding/json"
	"facet/api/middleware"
	"io/ioutil"

	"facet/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BaseUrl                       = "/app"
	IdQueryParameter              = "id"
	NameQueryParameter            = "name"
	EnvironmentIdQueryParameter   = "environment"
	WorkspaceIdPathQueryParameter = "workspaceId"
)

func Route(router *gin.Engine) {
	router.GET(BaseUrl, middleware.APIKeyVerify(), Get)
	router.POST(BaseUrl, middleware.APIKeyVerify(), Post)
	router.DELETE(BaseUrl, middleware.APIKeyVerify(), Delete)
}

func Get(c *gin.Context) {
	app := App{}
	app.Id = c.Request.URL.Query().Get(IdQueryParameter)
	app.WorkspaceId = c.Request.URL.Query().Get(WorkspaceIdPathQueryParameter)
	app.Name = c.Request.URL.Query().Get(NameQueryParameter)
	app.Environment = c.Request.URL.Query().Get(EnvironmentIdQueryParameter)
	var appArray *[]App
	var error error
	if app.Name == "" && app.Id == "" {
		appArray, error = FetchAll(app.WorkspaceId)
	} else {
		error = app.fetch()
		appArray = &[]App{app}
	}
	util.SetResponseCode(appArray, error, c)
}

func Post(c *gin.Context) {
	app := App{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &app)
	error = app.create()
	util.SetResponseCode(app, error, c)
}

func Delete(c *gin.Context) {
	app := App{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &app)
	error = app.delete()
	util.SetResponseCode(nil, error, c)
}
