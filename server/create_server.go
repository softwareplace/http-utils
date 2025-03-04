package server

import (
	"github.com/gorilla/mux"
	"github.com/softwareplace/http-utils/api_context"
)

func CreateApiRouter[T api_context.ApiPrincipalContext](topMiddlewares ...ApiMiddleware[T]) ApiRouterHandler[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &apiRouterHandlerImpl[T]{
		router:                              router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	router.Use(api.errorHandlerWrapper)

	for _, middleware := range topMiddlewares {
		api.RegisterMiddleware(middleware, "")
	}
	return api.NotFoundHandler()
}

func CreateApiRouterWith[T api_context.ApiPrincipalContext](router mux.Router) ApiRouterHandler[T] {
	router.Use(rootAppMiddleware[T])
	api := &apiRouterHandlerImpl[T]{
		router:                              &router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	return api.NotFoundHandler()
}
