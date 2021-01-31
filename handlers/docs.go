// Package classification of Cluster API
//
// Documentation for Cluster API
//
//	Schemes: https
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta

package handlers
import "github.com/bybrisk/cluster-api/data"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// Response structure after cluster formation
// swagger:response clusterCreateResponse
type clusterCreateResponseWrapper struct {
	// Response structre after cluster creation
	// in: body
	Body data.CreateClusterResponse
}

// swagger:parameters createCluster
type createClusterParamsWrapper struct {
	// Cluster data structure to create cluster with BybID.
	// Note: the number of clusters should be provided in every iterations
	// in: body
	// required: true
	Body data.CreateClusterRequest
}