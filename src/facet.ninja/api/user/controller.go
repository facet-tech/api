package user

import (
	"encoding/json"
	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	BASE_URL              = "/user"
	EMAIL_QUERY_PARAMATER = "email"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	router.DELETE(BASE_URL, Delete)
}

func Get(c *gin.Context) {
	user := User{}
	user.Email = c.Request.URL.Query().Get(EMAIL_QUERY_PARAMATER)
	error := user.fetch()
	util.SetResponseCode(user, error, c)
}

func Post(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	error = user.Update()
	util.SetResponseCode(user, error, c)
}

func Delete(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	error = user.delete()
	util.SetResponseCode(nil, error, c)
}
