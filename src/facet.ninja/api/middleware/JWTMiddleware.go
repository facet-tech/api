package middleware

import (
	"errors"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

func AuthorizeJWT() gin.HandlerFunc {
	fmt.Println("mpika man mou!")
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		//fmt.Println("jwtToken!", tokenString)
		//keySet, err := jwk.Fetch("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_oM4ne6cSf/.well-known/jwks.json")
		//fmt.Println("keySet!", keySet, "err", err)

		//tokenString := "eyJraWQiOiIrKzZGaFhvcXpwMkhaakV6Z254Q0JobFR6cWZmdFJNQlFDYnZ0dzRUMHh3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIyOTkxYTUxNC04Y2UyLTQ4MjktYTljNi0zZmIzMWVmZmM0ODQiLCJldmVudF9pZCI6ImFiMjczYmZmLTcyM2UtNDYyYi04MWZkLWFkNjEyYTViM2U0MCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE2MDYwMzYzOTgsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy13ZXN0LTIuYW1hem9uYXdzLmNvbVwvdXMtd2VzdC0yX29NNG5lNmNTZiIsImV4cCI6MTYwNjAzOTk5OCwiaWF0IjoxNjA2MDM2Mzk4LCJqdGkiOiIzYTFlMjdmNS02MzBhLTQyZTYtYjg1YS00ODg5NThiZDE2ZmEiLCJjbGllbnRfaWQiOiI2ZmE0ZmhjdG5vanVmM2htbHZvMG12c3AwMiIsInVzZXJuYW1lIjoiZXJtaW4uZWlyaWtAcHJpbWFyeWFsZS5jb20ifQ.JqmCjpnF0Eb_6wL-h8a9_e2kMBwFVBSwDOokjZn7L-5ox7lqCRYwEM2LEJ4kGXgd45Fa3WG32hlQ8ZOwZg49ViMaiqvEapdGq3OyX6M37wPql76Bo-wAnoKQFPZ4LEE0ku95lb66HEpijpSsvfx2hfHNkWtU17-cDWyOqxIc-6vhuq7ZSRnCt8D1p4pPNic2FNgOxeQ6NKWOjcQGzv25dT-aTtJMAh10swf68kkeuf8ugEczWnerqE-nhUa3shXH2t0x0vv2j2WedyI3TTg3BVhRv7v0Y5_twe6Jlm-8WtJ4tJzFkS704oOJ0L4hzI8AlqdL5iEhuTVVGWQnen_9dw"
		fmt.Println(tokenString)

		token, err := jwt.Parse(tokenString, getKey)
		if err != nil {
			panic(err)
		}
		claims := token.Claims.(jwt.MapClaims)
		for key, value := range claims {
			fmt.Printf("%s\t%v\n", key, value)
		}
	}
}

const jwksURL = `https://cognito-idp.us-west-2.amazonaws.com/us-west-2_oM4ne6cSf/.well-known/jwks.json`

func getKey(token *jwt.Token) (interface{}, error) {

	// TODO: cache response so we don't have to make a request every time
	// we want to verify a JWT
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
