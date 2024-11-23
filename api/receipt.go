// Package receipt provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package receipt

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Item defines model for Item.
type Item struct {
	// The total price payed for this item.
	Price string `json:"price"`

	// The Short Product Description for the item.
	ShortDescription string `json:"shortDescription"`
}

// Receipt defines model for Receipt.
type Receipt struct {
	Items []Item `json:"items"`

	// The date of the purchase printed on the receipt.
	PurchaseDate openapi_types.Date `json:"purchaseDate"`

	// The time of the purchase printed on the receipt. 24-hour time expected.
	PurchaseTime string `json:"purchaseTime"`

	// The name of the retailer or store the receipt is from.
	Retailer string `json:"retailer"`

	// The total amount paid on the receipt.
	Total string `json:"total"`
}

// PostReceiptsProcessJSONBody defines parameters for PostReceiptsProcess.
type PostReceiptsProcessJSONBody Receipt

// PostReceiptsProcessJSONRequestBody defines body for PostReceiptsProcess for application/json ContentType.
type PostReceiptsProcessJSONRequestBody PostReceiptsProcessJSONBody

// Bind implements render.Binder.
func (PostReceiptsProcessJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// PostReceiptsProcessJSON200Response is a constructor method for a PostReceiptsProcess response.
// A *Response is returned with the configured status code and content type from the spec.
func PostReceiptsProcessJSON200Response(body struct {
	ID string `json:"id"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetReceiptsIDPointsJSON200Response is a constructor method for a GetReceiptsIDPoints response.
// A *Response is returned with the configured status code and content type from the spec.
func GetReceiptsIDPointsJSON200Response(body struct {
	Points *int64 `json:"points,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Submits a receipt for processing
	// (POST /receipts/process)
	PostReceiptsProcess(w http.ResponseWriter, r *http.Request) *Response
	// Returns the points awarded for the receipt
	// (GET /receipts/{id}/points)
	GetReceiptsIDPoints(w http.ResponseWriter, r *http.Request, id string) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// PostReceiptsProcess operation middleware
func (siw *ServerInterfaceWrapper) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostReceiptsProcess(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetReceiptsIDPoints operation middleware
func (siw *ServerInterfaceWrapper) GetReceiptsIDPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetReceiptsIDPoints(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Post("/receipts/process", wrapper.PostReceiptsProcess)
		r.Get("/receipts/{id}/points", wrapper.GetReceiptsIDPoints)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6yVUW+jRhDHv8pqem+HbUxc68JbW0uVVaWKLvcWUmnNDvFeze52dricFfHdqwVsY0wu",
	"iXQvCYZh5r/z/83wDLktnTVo2EP6DD7fYimbyzVjGf47sg6JNfr2l84xXCj0OWnH2hpI4csWBVuWO9EE",
	"CCf3qERhSfBWe6EZyylEgN9l6XYIKSyni2uIwElmpJDhnyxTH7NsmmXqOak/QAS8dyHSM2nzCHUEfmuJ",
	"V/26YzLuQpS4JauqnEUvvJODI2pubGVYaiNW+CTmye1f59Lus+wpy3yWTR4+jiirIyD8r9KECtL7S5lR",
	"17WH45t28xVzDmf6jDlqx5eNDiLPLz4QFpDCL7OTZbPOr1ljVh1Bqc26jZ8fi0kiuQ8PXUX5VnpcSX7B",
	"QiUZhS2aLh2ig6OGUQlrmvvUKj5vYBInySSeT+I5RFBYKiVDCiHdmJGH1F90+RJLunyzEJEsJltbUfsS",
	"fneYM6pzffOr9FxaiB2TRshS75DGZRl5knWIFJaEZ0vYFyW0FwXZIWZZFcfJ8kb8YckgiRtJ/yK/xFob",
	"PEpcBM2w/WgOZRmYFk7qHzv3/kEc4H7s2ACwgc1RB/JB+uUwhMTaFPbyVL8Jr4PcY3cd2Ry9t6Eoa24O",
	"0k1SmPzjs29Ivk0xn8bTODTOOjTSaUjhahpPr9qjb5sBm3Xp/azL30yl9Xyp6K7alJq9kEdJYbl0r4Uu",
	"NYVIhvC1ghRuredOoe8UQttH9Py7VftQI7eG0TTlpHM7nTfvz776dte1w/7aKjhslPrcKKYKmxveWePb",
	"HZPE8bvKDjaUCn9PJEm1WW5+XcaTGLGYLJJNPrlW8+VEFYtPxVWMn643yZC0uzfsU61egOXcks/IFRnf",
	"kL5eCem9fjSoBNs+/AGBRXvsy8npDa823+ROq0aMr8pS0v4ttofwE0bPWtUzZ3X3cX3EEZL6sttQIZ8k",
	"qePn8yR9yNSfeERqvbpty4T+kiyRkTyk92OnXK9OK+yQWIeHYRAggrDkIA1tH/IT9Vl43caHn4rbqY9H",
	"5OZx3Fvp2vBycZIRvhWPSI0lr8LTLPeq3IR1Xgx8aJFZXFr3t+1hUJmDYZLFBTjvcLmu6/r/AAAA//95",
	"EUThlwkAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
