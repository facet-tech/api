package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/facet/:id", getFacet)
		r.POST("/facet/:id", createFacet)
		r.DELETE("/facet/:id", deleteFacet)
		r.GET("/js/:id/facet.ninja.js", getJs)
		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

type myId struct {
	Id string `form:"id"`
}

func getJs(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	var commaSeperatedIdsString string
	javascript := js()
	site := c.Param("id")
	if &site != nil {
		log.Println("id2:=" + site)
		facets, error := getItem(site)
		for _, facet := range facets.Facet {
			for _, facetId := range facet.Id {
				commaSeperatedIdsString += commaSeperatedIdsString + ",\"" + facetId + "\""
			}
		}
		javascript = strings.Replace(javascript, "GO_ARRAY_REPLACE_ME", strings.TrimPrefix(commaSeperatedIdsString, ","), -1)
		if error == nil {
			c.Data(200, "text/javascript", []byte(javascript))
		} else {
			c.JSON(500, error)
		}
	} else {
		c.JSON(400, "id is required")
	}
}

func js() string {
	script := `function getDomPath(el) {
	var stack = [];
	while (el.parentNode != null) {
		var sibCount = 0;
		var sibIndex = 0;
		for (var i = 0; i < el.parentNode.childNodes.length; i++) {
			var sib = el.parentNode.childNodes[i];
			if (sib.nodeName == el.nodeName) {
				if (sib === el) {
					sibIndex = sibCount;
				}
				sibCount++;
			}
		}
		if (el.hasAttribute('id') && el.id != '') {
			stack.unshift(el.nodeName.toLowerCase() + '#' + el.id);
		} else if (sibCount > 1) {
			stack.unshift(el.nodeName.toLowerCase() + ':eq(' + sibIndex + ')');
		} else {
			stack.unshift(el.nodeName.toLowerCase());
		}
		el = el.parentNode;
	}
	var aa = stack.slice(1);
	return aa.join(' > ');
}
//aHR0cHM6Ly9teXdlYnNpdGUuZmFjZXQubmluamEv
var all = document.getElementsByTagName("*");
var nodesToRemove = new Set([
    GO_ARRAY_REPLACE_ME
])

const callback = async function(mutationsList, observer) {
for(let mutation of mutationsList) {
if(nodesToRemove.has(getDomPath(mutation.target))) {
mutation.target.remove();
}
}
};

const targetNode = document
const config = { subtree: true, childList: true };
const observer = new MutationObserver(callback);
observer.observe(targetNode, config);`
	return script
}

func getFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	facets, error := getItem(site)

	if error == nil {
		c.JSON(200, facets)
	} else {
		c.JSON(500, error)
	}
}

func deleteFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	error := deleteItem(site)

	if error == nil {
		c.JSON(200, "Site "+site+"deleted.")
	} else {
		c.JSON(500, error)
	}

}

func createFacet(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Header", "*")
	site := c.Param("id")
	newFacet := Facets{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &newFacet)
	error = putItem(site, newFacet)

	if error == nil {
		c.JSON(201, "Created facet: "+site)
	} else {
		c.JSON(500, error)
	}

}
