package domain

import (
	"encoding/json"
	"io/ioutil"

	"facet/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL                          = "/domain"
	DOMAIN_QUERY_PARAMATER            = "domain"
	WORKSAPCE_ID_PATH_QUERY_PARAMATER = "workspaceId"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	router.DELETE(BASE_URL, Delete)
}

func Get(c *gin.Context) {
	domain := Domain{}
	domain.Domain = c.Request.URL.Query().Get(DOMAIN_QUERY_PARAMATER)
	domain.WorkspaceId = c.Request.URL.Query().Get(WORKSAPCE_ID_PATH_QUERY_PARAMATER)
	error := domain.fetch()
	util.SetResponseCode(domain, error, c)
}

func Post(c *gin.Context) {
	domain := Domain{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &domain)
	error = domain.create()
	util.SetResponseCode(domain, error, c)
}

func Delete(c *gin.Context) {
	domain := Domain{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &domain)
	error = domain.delete()
	util.SetResponseCode(nil, error, c)
}
