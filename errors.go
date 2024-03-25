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
	"regexp"
	"strconv"
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

// newErrorResponseParse creates a new *ErrorResponse from the given
// http.Response and error. If parse is true, we try to unmarshal the body as
// JSON to extract some additional info about the error.
func newErrorResponseParse(resp *http.Response, err error, parse bool) *ErrorResponse {
	errResp := &ErrorResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Method:     resp.Request.Method,
		Url:        resp.Request.URL,
		Err:        err,
	}

	data, err := peekResponseBody(resp)
	errResp.ResponseBody = data
	if !parse {
		return errResp
	}

	responseSeemsJSON := func() bool {
		if responseBodyIsJSON(resp.Header) {
			return true
		}
		if responseBodyIsPlainText(resp.Header) {
			if len(data) > 0 && data[0] == '{' {
				return true
			}
		}
		return false
	}

	if err == nil && len(data) > 0 && responseSeemsJSON() {
		var b struct {
			TraceID *string `json:"trace_id,omitempty"`
			Message *string `json:"message,omitempty"`
			Error   *string `json:"eroare,omitempty"`
		}
		// Don't care if we get an error here.
		_ = json.Unmarshal(data, &b)
		errResp.TraceID = b.TraceID
		errResp.Message = b.Message
		if b.Message == nil && b.Error != nil {
			errResp.Message = b.Error
		}
	}

	return errResp
}

// newErrorResponse is synonym to newErrorResponseParse(resp, err, true)
func newErrorResponse(resp *http.Response, err error) *ErrorResponse {
	return newErrorResponseParse(resp, err, true)
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
	Builder string
	Term    *string
}

// newBuilderNameErrorf creates a new Builder error for the given builder name,
// term and format string and args.
func newBuilderNameErrorf(builder, term string, format string, a ...any) *BuilderError {
	return &BuilderError{
		error:   fmt.Errorf(format, a...),
		Builder: builder,
		Term:    ptrfyStringNotEmpty(term),
	}
}

// newBuilderErrorf same as newBuilderNameErrorf but the builder name is taken
// from the type name of builder.
func newBuilderErrorf(builder any, term string, format string, a ...any) *BuilderError {
	return newBuilderNameErrorf(typeNameAddrPtr(builder), term, format, a...)
}

func (e *BuilderError) Error() string {
	errMsg := fmt.Sprintf("%s: ", e.Builder)
	if e.Term != nil {
		errMsg += fmt.Sprintf("term=%s: ", *e.Term)
	}
	errMsg += e.error.Error()
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
	*ErrorResponse
	Limit int64
}

func newLimitExceededError(r *http.Response, limit int64, err error) *LimitExceededError {
	return &LimitExceededError{
		ErrorResponse: newErrorResponseParse(r, err, false),
		Limit:         limit,
	}
}

func (e *LimitExceededError) Error() string {
	return e.ErrorResponse.Error()
}

var (
	regexLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) .* in cursul zilei")
	regexDownloadID       = regexp.MustCompile("id_descarcare=(\\d+)")
	regexCUI              = regexp.MustCompile("CUI=\\s*(\\d+)")

	// regexDownloadLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) descarcari la mesajul cu id_descarcare=(\\d+)\\s* in cursul zilei")
	// regexGetMessagesLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) de interogari de tip lista mesaje .* de catre CUI=\\s*(\\d+)\\s* in cursul zilei")
)

func errorMessageMatchLimitExceeded(err string) (limit int64, match bool) {
	if m, ok := matchFirstSubmatch(err, regexLimitExceededMsg); ok {
		limit, _ := strconv.ParseInt(m, 10, 64)
		return limit, true
	}
	return
}

func newErrorResponseDetectType(resp *http.Response) error {
	data, err := peekResponseBody(resp)
	if err == nil && len(data) > 0 {
		var b struct {
			Title string `json:"titlu"`
			Error string `json:"eroare"`
		}
		// Don't care if we get an error here.
		_ = json.Unmarshal(data, &b)
		if b.Error != "" {
			rerr := fmt.Errorf("%s: %s", b.Title, b.Error)
			if limit, ok := errorMessageMatchLimitExceeded(b.Error); ok {
				return newLimitExceededError(resp, limit, rerr)
			}

			return newErrorResponseParse(resp, rerr, false)
		}
	}

	return newErrorResponse(resp, nil)
}
