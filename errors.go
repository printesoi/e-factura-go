package efactura

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ErrorResponse struct {
	StatusCode   int
	Status       string
	Method       string
	Url          *url.URL
	ResponseBody []byte
	TraceID      *string
	Err          error
}

func NewErrorResponse(resp *http.Response, err error) *ErrorResponse {
	errResp := &ErrorResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Method:     resp.Request.Method,
		Url:        resp.Request.URL,
		Err:        err,
	}

	data, err := peekResponseBody(resp)
	errResp.ResponseBody = data
	if err == nil && len(data) > 0 && responseBodyIsJSON(resp.Header) {
		var b struct {
			TraceID *string `json:"trace_id,omitempty"`
		}
		_ = json.Unmarshal(data, &b)
		errResp.TraceID = b.TraceID
	}

	return errResp
}

func (r *ErrorResponse) Error() string {
	m := fmt.Sprintf("ANAF API call error: %v %v: %s", r.Method, r.Url, r.Status)
	if r.TraceID != nil {
		m += fmt.Sprintf("; trace_id=%s", *r.TraceID)
	}
	if r.Err != nil {
		m += fmt.Sprintf("; error=%s", r.Err.Error())
	}
	return m
}
