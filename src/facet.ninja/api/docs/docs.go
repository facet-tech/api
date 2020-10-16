package docs

import (
	"facet.ninja/api/domain"
)

// swagger:route GET /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   400: notFoundResponseWrapper

// swagger:route POST /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   200:

// swagger:route DELETE /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//   400

// This text will appear as description of your response body.
// swagger:response domainResponseWrapper
type domainResponseWrapper struct {
	// in:body
	Body domain.Domain
}

// Not Found
// swagger:response NotFoundResponseWrapper
type notFoundResponseWrapper struct{}
