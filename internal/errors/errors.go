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
	"encoding/json"
	std_errors "errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	api_helpers "github.com/printesoi/e-factura-go/internal/helpers/api"
	"github.com/printesoi/e-factura-go/internal/ptr"
	iregexp "github.com/printesoi/e-factura-go/internal/regexp"
	errors "github.com/printesoi/e-factura-go/pkg/errors"
)

var (
	New    = std_errors.New
	As     = std_errors.As
	Is     = std_errors.Is
	Unwrap = std_errors.Unwrap
)

// NewErrorResponseParse creates a new *errors.ErrorResponse from the given
// http.Response and error. If parse is true, we try to unmarshal the body as
// JSON to extract some additional info about the error.
func NewErrorResponseParse(resp *http.Response, err error, parse bool) *errors.ErrorResponse {
	errResp := &errors.ErrorResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Method:     resp.Request.Method,
		Url:        resp.Request.URL,
		Err:        err,
	}

	data, err := api_helpers.PeekResponseBody(resp)
	errResp.ResponseBody = data
	if !parse {
		return errResp
	}

	responseSeemsJSON := func() bool {
		if api_helpers.ResponseBodyIsJSON(resp.Header) {
			return true
		}
		if api_helpers.ResponseBodyIsPlainText(resp.Header) {
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

// NewErrorResponse is synonym to NewErrorResponseParse(resp, err, true)
func NewErrorResponse(resp *http.Response, err error) *errors.ErrorResponse {
	return NewErrorResponseParse(resp, err, true)
}

// NewBuilderNameErrorf creates a new Builder error for the given builder name,
// term and format string and args.
func NewBuilderNameErrorf(builder, term string, format string, a ...any) *errors.BuilderError {
	return &errors.BuilderError{
		Err:     fmt.Errorf(format, a...),
		Builder: builder,
		Term:    ptr.StringNotEmpty(term),
	}
}

// NewBuilderErrorf same as NewBuilderNameErrorf but the builder name is taken
// from the type name of builder.
func NewBuilderErrorf(builder any, term string, format string, a ...any) *errors.BuilderError {
	return NewBuilderNameErrorf(typeNameAddrPtr(builder), term, format, a...)
}

func NewLimitExceededError(r *http.Response, limit int64, err error) *errors.LimitExceededError {
	return &errors.LimitExceededError{
		ErrorResponse: NewErrorResponseParse(r, err, false),
		Limit:         limit,
	}
}

var (
	regexLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) .* in cursul zilei")
	regexDownloadID       = regexp.MustCompile("id_descarcare=(\\d+)")
	regexCUI              = regexp.MustCompile("CUI=\\s*(\\d+)")

	// regexDownloadLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) descarcari la mesajul cu id_descarcare=(\\d+)\\s* in cursul zilei")
	// regexGetMessagesLimitExceededMsg = regexp.MustCompile("S-au facut deja (\\d+) de interogari de tip lista mesaje .* de catre CUI=\\s*(\\d+)\\s* in cursul zilei")
)

func ErrorMessageMatchLimitExceeded(err string) (limit int64, match bool) {
	if m, ok := iregexp.MatchFirstSubmatch(regexLimitExceededMsg, err); ok {
		limit, _ := strconv.ParseInt(m, 10, 64)
		return limit, true
	}
	return
}

func NewErrorResponseDetectType(resp *http.Response) error {
	data, err := api_helpers.PeekResponseBody(resp)
	if err == nil && len(data) > 0 {
		var b struct {
			Title string `json:"titlu"`
			Error string `json:"eroare"`
		}
		// Don't care if we get an error here.
		_ = json.Unmarshal(data, &b)
		if b.Error != "" {
			rerr := fmt.Errorf("%s: %s", b.Title, b.Error)
			if limit, ok := ErrorMessageMatchLimitExceeded(b.Error); ok {
				return NewLimitExceededError(resp, limit, rerr)
			}

			return NewErrorResponseParse(resp, rerr, false)
		}
	}

	return NewErrorResponse(resp, nil)
}

func typeNameAddrPtr(v any) string {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	return rt.Name()
}
