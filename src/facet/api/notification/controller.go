package notification

import (
	"encoding/json"
	"facet/api/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	BASE_URL              = "/notification"
)

func Route(router *gin.Engine) {
	router.POST(BASE_URL, Post)
}


func Post(c *gin.Context) {
	notification := Notification{}
	body, error := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(body, &notification)
	error = notification.SendBatch()
	util.SetResponseCode(notification, error, c)
}