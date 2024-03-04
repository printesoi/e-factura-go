package efactura

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ErrorResponse is an error returns if the HTTP requests was finished (we got
// a *http.Response from the HTTP client, but it was not a successful response,
// or it was an error parsing the response.
type ErrorResponse struct {
	StatusCode   int
	Status       string
	Method       string
	Url          *url.URL
	ResponseBody []byte
	Err          error
	TraceID      *string
	Message      *string
}

func newErrorResponse(resp *http.Response, err error) *ErrorResponse {
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
			Message *string `json:"message,omitempty"`
		}
		// Don't care if we get an error here.
		json.Unmarshal(data, &b)
		errResp.TraceID = b.TraceID
		errResp.Message = b.Message
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
	if r.Message != nil {
		m += fmt.Sprintf("; message=%s", *r.Message)
	}
	return m
}
