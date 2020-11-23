package workspace

import (
	"encoding/json"
	"io/ioutil"

	"facet.ninja/api/middleware"
	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL           = "/workspace"
	ID_QUERY_PARAMATER = "id"
)

func AuthenticatedRoute(router *gin.Engine) {
	router.GET(BASE_URL, middleware.JWTVerify(), Get)
	router.POST(BASE_URL, middleware.JWTVerify(), Post)
	router.DELETE(BASE_URL, middleware.JWTVerify(), Delete)
}

func Get(c *gin.Context) {
	workspace := Workspace{}
	workspace.Id = c.Request.URL.Query().Get(ID_QUERY_PARAMATER)
	error := workspace.fetch()
	util.SetResponseCode(workspace, error, c)
}

func Post(c *gin.Context) {
	workspace := Workspace{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &workspace)
	error = workspace.create()
	util.SetResponseCode(workspace, error, c)
}

func Delete(c *gin.Context) {
	workspace := Workspace{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &workspace)
	error = workspace.delete()
	util.SetResponseCode(nil, error, c)
}
