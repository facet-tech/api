package user

import (
	"encoding/json"
	"io/ioutil"

	"facet.ninja/api/util"
	"github.com/gin-gonic/gin"
)

const (
	BASE_URL              = "/user"
	VERIFY_SETUP_URL      = "/user/signup/verify"
	SIGNUP_URL            = "/user/signup"
	LOGIN_URL             = "/user/login"
	EMAIL_QUERY_PARAMATER = "email"
)

func Route(router *gin.Engine) {
	router.GET(BASE_URL, Get)
	router.POST(BASE_URL, Post)
	router.POST(SIGNUP_URL, Signup)
	router.POST(LOGIN_URL, Login)
	router.POST(VERIFY_SETUP_URL, VerifySetupURL)
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

func Signup(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	error = user.addUserToUserPool()
	if error != nil {
		// need to cancel all entries that were previously added
	}
	util.SetResponseCode(user, error, c)
}

func VerifySetupURL(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	error = user.verify()
	util.SetResponseCode(nil, error, c)
}

func Login(c *gin.Context) {
	user := User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &user)
	authOutput := user.login()
	util.SetResponseCode(authOutput, error, c)
}
