package main

import (
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/internal/service"
	"github.com/softwareplace/http-utils/internal/service/api"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/secret"
	"github.com/softwareplace/http-utils/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorHandlerImpl struct {
}

func (p *errorHandlerImpl) Handler(ctx *apicontext.Request[*apicontext.DefaultContext], _ error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}

func init() {
	logger.LogReportCaller = true
	logger.LogSetup()

}

var (
	userPrincipalService                                                     = service.NewPrincipalService()
	errorHandler         errorhandler.ApiHandler[*apicontext.DefaultContext] = &errorHandlerImpl{}
	securityService                                                          = security.New(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
		errorHandler,
	)

	loginService   = service.NewLoginService(securityService)
	secretProvider = service.NewSecretProvider()

	secretHandler = secret.New(
		"./secret/private.key",
		secretProvider,
		securityService,
	)

	apiSecret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcGlLZXkiOiJTL2VTYzVDQ3Jub1laaDAyU2pLdFVsSzFXdmRaaVA1OXpFUU9jNE54K0pjL1c1dkhMa0tndE1ueExHN3dKTUwvIiwiY2xpZW50IjoiU29mdHdhcmUgUGxhY2UiLCJleHAiOjMwMzM5MzczNTcsInNjb3BlIjpbIkhFMSs0cEVwM3YzZFBzWXNLa3FLMGkzdiswSjMvYjFVN01YQkx3ZzhxQ0E9IiwiR2lQWUVNU1IvK1BjNUdaTm9OcUpqZDRkS1FZbjZ6QzBMbmdYTHVxdFc4VzkiLCJjY294TWNaT0tEZ0srTUZuend0YWFEWXgxaEtPSVlKNDl3PT0iLCJxOWRHb3V5bTBxZWxvV1V4bElKZ2Y1U3l6UnIrU3YwWWwvVT0iLCJNOHdrRkN3cmZpeVBKc2hjb3NrQU5GS0RZZ2ZxRnJOWXkwVmljOEdlM3dPSyJdfQ.n5_8kp3nNqXOAZVB73GCIXcv61gNyyihqz6xDIjIA0k"
)

func TestMockServer(t *testing.T) {
	t.Run("expects that can get login response successfully", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			LoginResource(loginService).
			SecurityService(securityService).
			PrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that return 401 when api secret is required for all resources but was not provided", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		secretProvider := service.NewSecretProvider()
		secretHandler := secret.New(
			"./secret/private.key",
			secretProvider,
			securityService,
		)

		server.Default().
			LoginResource(loginService).
			SecretService(secretHandler).
			SecurityService(securityService).
			PrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	t.Run("expects that can get login response successfully when requires api secret and it was provided", func(t *testing.T) {
		// Create a new request
		loginBody := strings.NewReader(`{"username": "my-username","password": "ynT9558iiMga&ayTVGs3Gc6ug1"}`)
		req, err := http.NewRequest("POST", "/login", loginBody)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set(apicontext.XApiKey, apiSecret)

		rr := httptest.NewRecorder()

		server.Default().
			LoginResource(loginService).
			SecretService(secretHandler).
			SecurityService(securityService).
			PrincipalService(userPrincipalService).
			NotFoundHandler().
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("expects that return default not found when a custom was not provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			PrincipalService(userPrincipalService).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "404 page not found") {
			t.Log("Response body contains '404 page not found'")
		} else {
			t.Errorf("Expected response body to contain '404 page not found', but got: %s", rr.Body.String())
		}

	})

	t.Run("expects that return custom not found when a custom was provided", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("POST", "/not-found", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			PrincipalService(userPrincipalService).
			CustomNotFoundHandler(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("Custom 404 Page"))
			}).
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		if strings.Contains(rr.Body.String(), "Custom 404 Page") {
			t.Log("Response body contains 'Custom 404 Page'")
		} else {
			t.Errorf("Expected response body to contain 'Custom 404 Page', but got: %s", rr.Body.String())
		}
	})

	t.Run("expects that return swagger resource when swagger was defined and using the default not found handler", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("GET", "/", nil)

		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()

		server.Default().
			PrincipalService(userPrincipalService).
			EmbeddedServer(api.Handler).
			SwaggerDocHandler("./internal/resource/pet-store.yaml").
			ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMovedPermanently)
		}

		if strings.Contains(rr.Body.String(), "<a href=\"/swagger/index.html\">Moved Permanently</a>.") {
			t.Log("Response body contains '<a href=\"/swagger/index.html\">Moved Permanently</a>.'")
		} else {
			t.Errorf("Expected response body to contain '<a href=\"/swagger/index.html\">Moved Permanently</a>.', but got: %s", rr.Body.String())
		}
	})
}
