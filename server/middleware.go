package server

import (
	"encoding/json"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"log"
	"net/http"
	"time"
)

// rootAppMiddleware logs each incoming request's method, path, and remote address
func rootAppMiddleware[T api_context.ApiPrincipalContext](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx *api_context.ApiRequestContext[T]

		error_handler.Handler(func() {
			start := time.Now() // Record the start time
			ctx = api_context.Of[T](w, r, "MIDDLEWARE/ROOT_APP")
			queryParam := ""
			if r.URL.RawQuery != "" {
				queryParam = "?" + r.URL.RawQuery
			}

			log.Printf("[%s]:: Incoming request: %s %s from %s", ctx.GetSessionId(), r.Method, r.URL.Path+queryParam, r.RemoteAddr)

			ctx.Next(next)

			duration := time.Since(start)

			log.Printf("[%s]:: => request processed: %s %s in %v",
				ctx.GetSessionId(),
				r.Method,
				r.URL.Path+queryParam,
				duration,
			)

		}, func(err error) {
			onError(err, w)
		})

		defer func() {
			error_handler.Handler(ctx.Flush, func(err error) {
				log.Printf("Error flushing context: %v", err)
			})
			ctx = nil
		}()
	})
}

func onError(err any, w http.ResponseWriter) {
	log.Printf("Error processing request: %+v", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	responseBody := map[string]interface{}{
		"message":    "Failed to process request",
		"statusCode": http.StatusInternalServerError,
		"timestamp":  time.Now().UnixMilli(),
	}

	err = json.NewEncoder(w).Encode(responseBody)

	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
