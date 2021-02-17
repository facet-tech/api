package workspace

import (
	"encoding/json"
	"io/ioutil"

	"facet/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL           = "/workspace"
	ID_QUERY_PARAMATER = "id"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	router.DELETE(BASE_URL, Delete)
}

func Get(c *gin.Context) {
	workspace := Workspace{}
	workspace.WorkspaceId = c.Request.URL.Query().Get(ID_QUERY_PARAMATER)
	workspace.Id = c.Request.URL.Query().Get(ID_QUERY_PARAMATER)
	workspaces, error := workspace.fetchAll()
	util.SetResponseCode(workspaces, error, c)
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
