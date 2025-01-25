// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package gen

import (
	"fmt"

	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

const (
	OAuth2Scopes = "OAuth2.Scopes"
)

// BaseResponse defines model for BaseResponse.
type BaseResponse struct {
	Code      *int    `json:"code,omitempty"`
	Message   *string `json:"message,omitempty"`
	Success   *bool   `json:"success,omitempty"`
	Timestamp *int    `json:"timestamp,omitempty"`
}

func GetBaseResponseBody[T api_context.ApiPrincipalContext](
	ctx *api_context.ApiRequestContext[T],
	onSuccess server.OnSuccess[BaseResponse, T],
	onError server.OnError[T],
) {
	server.GetRequestBody(ctx, BaseResponse{}, onSuccess, onError)
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func GetLoginRequestBody[T api_context.ApiPrincipalContext](
	ctx *api_context.ApiRequestContext[T],
	onSuccess server.OnSuccess[LoginRequest, T],
	onError server.OnError[T],
) {
	server.GetRequestBody(ctx, LoginRequest{}, onSuccess, onError)
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Expires *int    `json:"expires,omitempty"`
	Token   *string `json:"token,omitempty"`
}

func GetLoginResponseBody[T api_context.ApiPrincipalContext](
	ctx *api_context.ApiRequestContext[T],
	onSuccess server.OnSuccess[LoginResponse, T],
	onError server.OnError[T],
) {
	server.GetRequestBody(ctx, LoginResponse{}, onSuccess, onError)
}

// GetTestVersionParams defines parameters test of for GetTestVersion.
type GetTestVersionParams struct {
	// Ref Any data
	Ref string `form:"ref" json:"ref"`

	// Authorization jwt
	Authorization string `json:"Authorization"`
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody = LoginRequest

// ServerInterface represents all server handlers represents all server handlers..
type ServerInterface[T api_context.ApiPrincipalContext] interface {
	// Authentication endpoint
	// (POST /login)
	PostLogin(ctx *api_context.ApiRequestContext[T])
	// Public endpoint
	// (GET /test)
	GetTest(ctx *api_context.ApiRequestContext[T])
	// Secured endpoint
	// (GET /test/{version})
	GetTestVersion(ctx *api_context.ApiRequestContext[T])
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

//PostLogin operation middleware test execution
//GetTest operation middleware test execution
//GetTestVersion operation middleware test execution

// ResourcesHandler registers API endpoints and sets up the Swagger documentation.
//
// This function takes an instance of `ApiRouterHandler` and `ServerInterface`,
// and configures the following:
// - Sets up Swagger documentation using the provided `GetSwagger` function.
// - PostLogin
// - GetTest
// - GetTestVersion
// Parameters:
//   - apiServer: The API router handler used for setting up routes and middleware.
//   - server: The server interface implementation containing the endpoint handlers.
//
// Generics:
//   - T: A type that satisfies the `ApiPrincipalContext` interface, representing the principal/context
//     involved in the API operations.
func ResourcesHandler[T api_context.ApiPrincipalContext](apiServer server.ApiRouterHandler[T], server ServerInterface[T]) {

	apiServer.Add(server.PostLogin, "/login", "POST", []string{}...)

	apiServer.Add(server.GetTest, "/test", "GET", []string{}...)

	apiServer.Add(server.GetTestVersion, "/test/{version}", "GET", []string{"api:example:admin"}...)

}
