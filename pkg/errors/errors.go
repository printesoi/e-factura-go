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

package errors

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrInvalidTokenSource       = errors.New("invalid token source")
	ErrInvalidOAuth2Credentials = errors.New("invalid OAuth2 credentials")
	ErrInvalidOAuth2Endpoint    = errors.New("invalid OAuth2 endpoint")
	ErrInvalidOAuth2RedirectURL = errors.New("invalid OAuth2 redirect URL")
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
	Err     error
	Builder string
	Term    *string
}

func (e *BuilderError) Error() string {
	errMsg := fmt.Sprintf("%s: ", e.Builder)
	if e.Term != nil {
		errMsg += fmt.Sprintf("term=%s: ", *e.Term)
	}
	errMsg += e.Err.Error()
	return errMsg
}

// ValidateSignatureError is an error returned if the signature cannot be
// succesfully validated.
type ValidateSignatureError struct {
	error
}

func newValidateSignatureError(err error) *ValidateSignatureError {
	return &ValidateSignatureError{error: err}
}

// LimitExceededError is an error returned if we hit an API limit.
type LimitExceededError struct {
	// ErrorResponse has information about the HTTP response.
	*ErrorResponse
	// Limit stores the API limit that was hit for the day.
	Limit int64
}

func (e *LimitExceededError) Error() string {
	return e.ErrorResponse.Error()
}
