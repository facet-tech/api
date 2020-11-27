package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"os"
)

func JWTVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: cache valid JWTs https://github.com/facets-io/api/issues/14
		tokenString := c.GetHeader("AccessToken")
		_, err := jwt.Parse(tokenString, getKey)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// TODO verify claims against the Cognito Pool https://github.com/facets-io/api/issues/13
	}
}
func getKey(token *jwt.Token) (interface{}, error) {
	jwksURL := os.Getenv("COGNITO_JWKS_URL")
	set, err := jwk.FetchHTTP(jwksURL)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	keys := set.LookupKeyID(keyID)
	if len(keys) == 0 {
		return nil, fmt.Errorf("key %v not found", keyID)
	}

	if key := set.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, fmt.Errorf("unable to find key %q", keyID)
}
