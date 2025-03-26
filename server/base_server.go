package server

import (
	"github.com/gorilla/mux"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/login"
	"github.com/softwareplace/goserve/security/principal"
	"github.com/softwareplace/goserve/security/secret"
	"net/http"
	"sync"
)

type baseServer[T apicontext.Principal] struct {
	router                              *mux.Router
	principalService                    principal.Service[T]
	errorHandler                        apicontext.ApiHandler[T]
	loginService                        login.Service[T]
	securityService                     security.Service[T]
	secretService                       secret.Service[T]
	apiKeyGeneratorService              ApiKeyGeneratorService[T]
	server                              *http.Server // Add a server instance
	mu                                  sync.Mutex   // Add a mutex for thread safety
	swaggerIsEnabled                    bool
	loginResourceEnable                 bool
	apiSecretKeyGeneratorResourceEnable bool
	contextPath                         string
	port                                string
}

func Default(
	topMiddlewares ...ApiMiddleware[*apicontext.DefaultContext],
) Api[*apicontext.DefaultContext] {
	return CreateApiRouter[*apicontext.DefaultContext](topMiddlewares...)
}

func CreateApiRouter[T apicontext.Principal](topMiddlewares ...ApiMiddleware[T]) Api[T] {
	router := mux.NewRouter()
	router.Use(rootAppMiddleware[T])

	api := &baseServer[T]{
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

func CreateApiRouterWith[T apicontext.Principal](router mux.Router) Api[T] {
	router.Use(rootAppMiddleware[T])
	api := &baseServer[T]{
		router:                              &router,
		apiSecretKeyGeneratorResourceEnable: true,
		loginResourceEnable:                 true,
		contextPath:                         apiContextPath(),
		port:                                apiPort(),
	}

	return api.NotFoundHandler()
}
