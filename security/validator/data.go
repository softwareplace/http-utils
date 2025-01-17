package validator

import (
	"github.com/softwareplace/http-utils/api_context"
	"regexp"
	"strings"
)

// HasResourceAccessRight checks if the user has the necessary roles to access the requested resource.
// It compares the roles assigned to the user with those required for the resource's path.
// If the path does not require any roles, the function returns true.
//
// Parameters:
//
//	ctx - The API request context containing user roles and request metadata.
//
// Returns:
//
//	bool - True if the user has the required roles or if the path does not require roles, false otherwise.
func HasResourceAccessRight[T api_context.ApiContextData](ctx api_context.ApiRequestContext[T]) bool {
	userRoles := (*ctx.AccessRolesLoader)(ctx)
	requiredRoles, isRoleRequired := GetRolesForPath(ctx)

	for _, requiredRole := range requiredRoles {
		for _, userRole := range userRoles {
			if requiredRole == userRole {
				return true
			}
		}
	}

	return !isRoleRequired
}

// GetRolesForPath retrieves the roles associated with a request path.
//
// This function takes the API request context and determines the roles required
// for accessing the specified path. The roles are matched based on predefined
// patterns or exact path matches.
//
// Parameters:
//
//	ctx - The API request context containing request metadata.
//
// Returns:
//
//	[]string - A slice of required roles for the path or nil if no roles are defined.
//	bool - True if roles are required for the path, false otherwise.
func GetRolesForPath[T api_context.ApiContextData](ctx api_context.ApiRequestContext[T]) ([]string, bool) {
	path := ctx.Request.Method + "::" + ctx.Request.URL.Path

	for pattern, requiredRoles := range roles {
		regexPattern := convertPathToRegex(pattern)
		regex := regexp.MustCompile(regexPattern)

		if regex.MatchString(path) || path == pattern {
			return requiredRoles, true
		}
	}

	return nil, false
}

// IsPublicPath checks if the provided path is registered as a public route.
//
// This function takes the API request context and verifies whether the current
// request path matches any registered public routes. Public routes are those
// that do not require any roles to access.
//
// Parameters:
//
//	ctx - The API request context containing request metadata.
//
// Returns:
//
//	bool - True if the path is a public route, false otherwise.
func IsPublicPath[T api_context.ApiContextData](ctx api_context.ApiRequestContext[T]) bool {
	path := ctx.Request.Method + "::" + ctx.Request.URL.Path
	for _, openPath := range openPaths {
		regexPattern := convertPathToRegex(openPath)
		regex := regexp.MustCompile(regexPattern)
		if regex.MatchString(path) || path == openPath {
			return true
		}
	}
	return false
}

// convertPathToRegex converts a path with dynamic segments (e.g., ":param") into a regex pattern.
func convertPathToRegex(path string) string {
	// Escape slashes and replace dynamic segments with regex groups.
	escapedPath := strings.ReplaceAll(path, "/", `\/`)
	return "^" + re.ReplaceAllString(escapedPath, `[^\/]+`) + "$"
}
