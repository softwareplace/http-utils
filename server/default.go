package server

import (
	apicontext "github.com/softwareplace/http-utils/context"
)

func Default(topMiddlewares ...ApiMiddleware[*apicontext.DefaultContext]) Api[*apicontext.DefaultContext] {
	return CreateApiRouter[*apicontext.DefaultContext](topMiddlewares...)
}
