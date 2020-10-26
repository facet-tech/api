package workspace

import (
	"encoding/json"
	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	BASE_URL           = "/workspace"
	ID_QUERY_PARAMATER = "id"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	//router.DELETE(BASE_URL, Delete)
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
