package user

import (
	"encoding/json"
	"facet.ninja/api/middleware"
	"io/ioutil"
	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL              = "/user"
	EMAIL_QUERY_PARAMATER = "email"
)

func AuthenticatedRoute(router *gin.Engine) {
	router.GET(BASE_URL, middleware.JWTVerify(), Get)
	router.POST(BASE_URL, middleware.JWTVerify(), Post)
	router.DELETE(BASE_URL, middleware.JWTVerify(), Delete)
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
	error = user.create()
	if error != nil {
		// need to cancel all entries that were added
	}
	util.SetResponseCode(user, error, c)
}

func Delete(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	error = user.delete()
	util.SetResponseCode(nil, error, c)
}
