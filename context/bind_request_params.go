package context

import (
	"fmt"
	"github.com/gorilla/mux"
	goservereflect "github.com/softwareplace/goserve/reflect"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// FieldSource indicates where a field originated from
type FieldSource string

// RequestError represents a validation error with contextual information
type RequestError struct {
	Field   string      `json:"field"`      // The original field name from request
	Source  FieldSource `json:"source"`     // Where the field came from
	Message string      `json:"message"`    // Human-readable error message
	Code    int         `json:"statusCode"` // HTTP status code
}

// Error implements the error interface
func (e *RequestError) Error() string {
	return fmt.Sprintf("%s %s", e.Source, e.Message)
}

func (ctx *Request[T]) FormValues() url.Values {
	r := ctx.Request
	if r.Form == nil {
		err := r.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			return nil
		}
	}
	return r.Form
}

// BindRequestParams extracts and binds request parameters such as query, form data, headers, or route vars into a target struct.
// It validates the target struct and returns a RequestError with details on validation failure, or nil on success.
func (ctx *Request[T]) BindRequestParams(target interface{}) *RequestError {
	r := ctx.Request

	contentType := ctx.Request.Header.Get(ContentType)

	if strings.Contains(contentType, MultipartFormData) {
		targetType := reflect.TypeOf(target)
		if body, ok := goservereflect.FindField(targetType, "body"); ok {
			_ = goservereflect.ParamsExtract(body,
				goservereflect.ParamsExtractorSource{
					Tree: ctx.FormValues(),
				},
			)
		}
	}
	_ = goservereflect.ParamsExtract(target,
		goservereflect.ParamsExtractorSource{
			Tree: r.URL.Query(),
		}, goservereflect.ParamsExtractorSource{
			Tree: r.Header,
		}, goservereflect.ParamsExtractorSource{
			Source: mux.Vars(r),
		},
	)

	err := ctx.StructValidation(target)

	if err != nil {
		return &RequestError{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}
	return nil
}
