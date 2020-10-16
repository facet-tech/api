package docs

import (
	"facet.ninja/api/domain"
	"github.com/pdrum/swagger-automation/api"
)

// swagger:route GET /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   200: domainResponseWrapper

// swagger:route POST /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   200:

// swagger:route DELETE /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   200

// This text will appear as description of your response body.
// swagger:response foobarResponse
type foobarResponseWrapper struct {
	// in:body
	Body api.FooBarResponse
}

// This text will appear as description of your response body.
// swagger:response domainResponseWrapper
type domainResponseWrapper struct {
	// in:body
	Body domain.Domain
}

// swagger:parameters idOfFoobarEndpoint
type foobarParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body api.FooBarRequest
}
