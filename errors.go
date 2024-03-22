// Copyright 2024 Victor Dodon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

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

// BuilderError is an error returned by the builders.
type BuilderError struct {
	error
	Term *string
}

func NewBuilderErrorf(format string, a ...any) *BuilderError {
	return &BuilderError{
		error: fmt.Errorf(format, a...),
	}
}

func NewBuilderTermErrorf(term string, format string, a ...any) *BuilderError {
	return &BuilderError{
		error: fmt.Errorf(format, a...),
	}
}

func (e *BuilderError) Error() string {
	if e.Term == nil {
		return e.error.Error()
	}

	return fmt.Sprintf("term: %s, error: %s", *e.Term, e.error.Error())
}

// ValidateSignatureError is an error returned if the signature cannot be
// succesfully validated.
type ValidateSignatureError struct {
	error
}

func NewValidateSignatureError(err error) *ValidateSignatureError {
	return &ValidateSignatureError{error: err}
}
