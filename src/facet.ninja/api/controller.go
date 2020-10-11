package main

import (
	"context"
	"encoding/json"
	"facet.ninja/api/domain"
	"facet.ninja/api/facet"
	"facet.ninja/api/user"
	"facet.ninja/api/workspace"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
)

const CREATED = "created"
const DELETED = "deleted"
const SUCCESS = "success"

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {

		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/facet", getFacet)
		r.POST("/facet", createFacet)
		r.DELETE("/facet", deleteFacet)
		r.OPTIONS("/facet", options)

		r.GET("/workspace", getWorkspaceController)
		r.POST("/workspace", createWorkspaceController)
		//r.DELETE("/site", deleteSiteController)

		r.GET("/domain", getDomainController)
		r.POST("/domain", createDomainController)
		//r.DELETE("/site", "/domain")

		r.GET("/user", getUserController)
		r.POST("/user", createUserController)
		r.DELETE("/user", deleteUserController)
		r.OPTIONS("/user", options)

		r.GET("/facet.ninja.js", getJs)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func addCorsHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
}

func getJs(c *gin.Context) {
	addCorsHeaders(c)
	var commaSeperatedIdsString string
	javascript := js()
	site := c.Request.URL.Query().Get("id")
	if &site != nil {
		facets, error := facet.FetchAll(site)
		for _, url := range *facets {
			commaSeperatedIdsString += "\t['" + url.UrlPath + "',new Set(["
			for _, facet := range url.DomElement {
				for _, facetId := range facet.Path {
					commaSeperatedIdsString += "'" + facetId + "',"
				}
				commaSeperatedIdsString = strings.TrimSuffix(commaSeperatedIdsString, ",")
			}
			commaSeperatedIdsString += "])],\n"
		}
		javascript = strings.Replace(javascript, "GO_ARRAY_REPLACE_ME", strings.TrimSuffix(commaSeperatedIdsString, ",\n"), -1)
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
    var res = stack.slice(1).join(' > '); // removes the html element
    return res;
}

var data = new Map([
GO_ARRAY_REPLACE_ME
])

const callback = async function(mutationsList, observer) {
    if ((typeof disableHideFacetNinja === 'undefined' || disableHideFacetNinja === null || disableHideFacetNinja === false) && data.has(window.location.pathname)) {
        let nodesToRemove = data.get(window.location.pathname)       
        for(let mutation of mutationsList) {
	        if(nodesToRemove.has(getDomPath(mutation.target))) {
	        	mutation.target.style.display = "none"
	        	mutation.target.style.setProperty("display", "none", "important");
            }
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
	domainId := c.Request.URL.Query().Get("domainId")
	urlPath := c.Request.URL.Query().Get("urlPath")
	facets, error := facet.Fetch(domainId, urlPath)
	responseCode(facets, error, c)
}

func deleteFacet(c *gin.Context) {
	request := facet.Facet{}
	body, error := ioutil.ReadAll(c.Request.Body)

	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	error = facet.Delete(request)
	responseCode(DELETED, error, c)
}

func options(c *gin.Context) {
	responseCode(SUCCESS, nil, c)
}

func createFacet(c *gin.Context) {
	request := facet.Facet{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	error = facet.Put(request)
	responseCode(CREATED, error, c)
}

func getDomainController(c *gin.Context) {
	domainName := c.Request.URL.Query().Get("domain")
	workspaceId := c.Request.URL.Query().Get("workspaceId")
	site, error := domain.Fetch(domainName, workspaceId)
	responseCode(site, error, c)
}

func createDomainController(c *gin.Context) {
	request := domain.Domain{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	result, error := domain.Create(request)
	responseCode(result, error, c)
}

func getWorkspaceController(c *gin.Context) {
	id := c.Request.URL.Query().Get("id")
	site, error := workspace.Fetch(id)
	responseCode(site, error, c)
}

func createWorkspaceController(c *gin.Context) {
	request := workspace.Workspace{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	result, error := workspace.Create(request)
	responseCode(result, error, c)
}

/*
/*
func deleteSiteController(c *gin.Context) {
	id := c.Param("id")
	error := workspace.Delete(id)
	responseCode(DELETED, error, c)
}*/

func getUserController(c *gin.Context) {
	email := c.Request.URL.Query().Get("email")
	user, error := user.Fetch(email)
	responseCode(user, error, c)
}

func createUserController(c *gin.Context) {
	request := user.User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	id, error := user.Create(request)
	responseCode(id, error, c)
}

func deleteUserController(c *gin.Context) {
	request := user.User{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		c.JSON(500, error)
	}
	json.Unmarshal(body, &request)
	error = user.Delete(request)
	responseCode(DELETED, error, c)
}

func bodyToJson(object interface{}, c *gin.Context) {
	body, error := ioutil.ReadAll(c.Request.Body)
	if error == nil {
		json.Unmarshal(body, &object)
	} else {
		c.JSON(500, error)
	}
}

func responseCode(result interface{}, error interface{}, c *gin.Context) {
	var response interface{}
	var responseCode int

	if error == nil {
		response = result
		responseCode = 500
	} else {
		response = error
	}

	if result == CREATED {
		responseCode = 201
	} else {
		responseCode = 200
	}

	addCorsHeaders(c)
	c.JSON(responseCode, response)
}
