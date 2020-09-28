package grpcgw

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// GWError is used for the error returned from the grpc implementation
// it can handle custom errors
type GWError interface {
	error
	// Status is the http status code
	Status() int
	// Message to outside user
	Message() string
	// Fields return the fields or nil
	Fields() map[string]string
}

type gwError struct {
	Err     error             `json:"-"`
	Msg     string            `json:"message"`
	Status_ int               `json:"status"`
	Fields_ map[string]string `json:"fields,omitempty"`
}

func (gw *gwError) Error() string {
	b, _ := json.Marshal(gw)
	return string(b)
}

func (gw *gwError) Status() int {
	return gw.Status_
}

func (gw *gwError) Message() string {
	return gw.Msg
}

func (gw *gwError) Fields() map[string]string {
	return gw.Fields_
}

// NewNotFound return not found error
func NewNotFound(err error) error {
	return NewBadRequestStatus(err, "Not found", http.StatusNotFound)
}

// NewBadRequest is the bad request
func NewBadRequest(err error, message string) error {
	return NewBadRequestStatus(err, message, http.StatusBadRequest)
}

// NewBadRequestStatus is the bad request
func NewBadRequestStatus(err error, message string, status int) error {
	ret := &gwError{
		Err:     errors.Wrap(err, message),
		Msg:     errors.Wrap(err, message).Error(),
		Status_: status,
	}
	if v, ok := err.(validation.Errors); ok {
		ret.Fields_ = make(map[string]string)
		for k, e := range v {
			ret.Fields_[k] = e.Error()
		}
	}
	return ret
}

type grpcErr interface {
	GRPCStatus() *status.Status
}

func tryGRPCError(err error) GWError {
	g, ok := err.(grpcErr)
	if !ok {
		return &gwError{
			Msg:     "unknown",
			Status_: http.StatusInternalServerError,
		}
	}
	switch g.GRPCStatus().Code() {
	case codes.InvalidArgument:
		return &gwError{
			Status_: http.StatusBadRequest,
			Msg:     "invalid json input",
		}
	default:
		return &gwError{
			Msg:     "internal error with code: " + g.GRPCStatus().Code().String(),
			Status_: http.StatusInternalServerError,
		}
	}
}

func tryJSONError(err error) (GWError, bool) {
	ret := &gwError{}
	txt := err.Error()
	if g, ok := err.(grpcErr); ok {
		txt = g.GRPCStatus().Message()
	}
	if e := json.Unmarshal([]byte(txt), ret); e != nil {
		return nil, false
	}
	return ret, true
}

// defaultHTTPError is my first try to overwrite the default
func defaultHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())
	g, ok := tryJSONError(err)
	if !ok || g.Status() == 0 {
		g = tryGRPCError(err)
	}

	body, ok := g.(*gwError)
	if !ok {
		body = &gwError{
			Err:     err,
			Msg:     g.Message(),
			Status_: g.Status(),
			Fields_: g.Fields(),
		}
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	w.WriteHeader(body.Status())
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func init() {
	runtime.HTTPError = defaultHTTPError
}
