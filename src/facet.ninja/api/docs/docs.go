package docs

import (
	"facet.ninja/api/domain"
)

// swagger:route GET /domain Domain idOfFoobarEndpoint
// Get a domain
// responses:
//	200: okResponseWrapper
//  400: notFoundResponseWrapper

// swagger:route POST /domain Domain domainResponseWrapper
// Create a Domain
// responses:
//	200: domainResponseWrapper
//	400: notFoundResponseWrapper

// swagger:route DELETE /domain Domain domainResponseWrapper
// Foobar does some amazing stuff.
// responses:
//	400: notFoundResponseWrapper

// This text will appear as description of your response body.
// swagger:response domainResponseWrapper
type domainResponseWrapper struct {
	// in:body
	Body domain.Domain
}

// OK
// swagger:response okResponseWrapper
type okResponseWrapper struct{}

// Not Found
// swagger:response notFoundResponseWrapper
type notFoundResponseWrapper struct{}

// swagger:parameters getDomainId
type domainParamWrapper struct {
	id string
}

// swagger:parameters idOfFoobarEndpoint
type foobarParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:query
	DomainID    string `json:"domainId"`
	WorkspaceId string `json:"workspaceId"`
}
