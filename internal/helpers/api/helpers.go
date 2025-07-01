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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
)

const (
	MediaTypeApplicationJSON = "application/json"
	MediaTypeApplicationXML  = "application/xml"
	MediaTypeApplicationPDF  = "application/pdf"
	MediaTypeApplicationZIP  = "application/zip"
	MediaTypeTextXML         = "text/xml"
	MediaTypeTextPlain       = "text/plain"
)

// This is a copy of the drainBody from src/net/http/httputil/dump.go
func DrainBody(b io.ReadCloser) (body []byte, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		return nil, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return buf.Bytes(), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func PeekResponseBody(r *http.Response) (body []byte, err error) {
	body, r.Body, err = DrainBody(r.Body)
	return
}

func PeekRequestBody(r *http.Request) (body []byte, err error) {
	body, r.Body, err = DrainBody(r.Body)
	return
}

func ResponseMediaType(headers http.Header) (mediaType string) {
	mediaType, _, _ = mime.ParseMediaType(headers.Get("Content-Type"))
	return
}

func ResponseBodyIsJSON(headers http.Header) bool {
	return ResponseMediaType(headers) == MediaTypeApplicationJSON
}

func ResponseBodyIsPlainText(headers http.Header) bool {
	return ResponseMediaType(headers) == MediaTypeTextPlain
}

func ResponseBodyIsXML(headers http.Header) bool {
	switch ResponseMediaType(headers) {
	case MediaTypeApplicationXML, MediaTypeTextXML:
		return true
	}
	return false
}

func ResponseIsSuccess(status int) bool {
	return status >= 200 && status < 300
}

// UnmarshalReaderJSON reads all the content from the given reader r and
// unmarshals the data as JSON into the value v.
func UnmarshalReaderJSON(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func BuildURL(base *url.URL, refUrl string, query url.Values) (string, error) {
	if base == nil {
		return "", fmt.Errorf("nil base URL")
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("empty base URL")
	}
	u, err := base.Parse(refUrl)
	if err != nil {
		return "", err
	}

	if u.RawQuery == "" {
		u.RawQuery = query.Encode()
	} else {
		qs := u.Query()
		for k, v := range query {
			qs[k] = v
		}
		u.RawQuery = qs.Encode()
	}
	return u.String(), nil
}

func BuildParseURL(baseUrl, refUrl string, query url.Values) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	return BuildURL(base, refUrl, query)
}
