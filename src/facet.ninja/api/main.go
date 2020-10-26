package main

import (
	"context"
	"facet.ninja/api/domain"
	"facet.ninja/api/facet"
	"facet.ninja/api/user"
	"facet.ninja/api/util"
	"facet.ninja/api/workspace"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"strings"
)

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		router := gin.Default()
		defaultRoutes(router)
		facet.Route(router)
		workspace.Route(router)
		domain.Route(router)
		user.Route(router)
		router.GET("/facet.ninja.js", getJs)
		ginLambda = ginadapter.New(router)
	}
	return ginLambda.ProxyWithContext(ctx, req)
}

func defaultRoutes(route *gin.Engine) {
	route.OPTIONS("/*anyPath", util.Options)
}

func getJs(c *gin.Context) {
	util.SetCorsHeaders(c)
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
    return res.replace(/\s+/g, '');
}

var data = new Map([
GO_ARRAY_REPLACE_ME
])

var facetedNodes = new Set()

const callback = async function(mutationsList, observer) {
    if ((typeof disableHideFacetNinja === 'undefined' || disableHideFacetNinja === null || disableHideFacetNinja === false) && data.has(window.location.pathname)) {
        let nodesToRemove = data.get(window.location.pathname)       
        for(let mutation of mutationsList) {
            let domPath = getDomPath(mutation.target)
	        if(nodesToRemove.has(domPath) && !facetedNodes.has(domPath)) {
	        	facetedNodes.add(domPath)
                mutation.target.style.display = "none"
	        	mutation.target.style.setProperty("display", "none", "important");
            }
            //console.log(getDomPath(mutation.target))
        }
    }
};

const targetNode = document
const config = { subtree: true, childList: true, attributes: true};
const observer = new MutationObserver(callback);
observer.observe(targetNode, config);`
	return script
}
