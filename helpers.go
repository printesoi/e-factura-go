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
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/m29h/xml"
)

// This is a copy of the drainBody from src/net/http/httputil/dump.go
func drainBody(b io.ReadCloser) (body []byte, r2 io.ReadCloser, err error) {
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

func peekResponseBody(r *http.Response) (body []byte, err error) {
	body, r.Body, err = drainBody(r.Body)
	return
}

func peekRequestBody(r *http.Request) (body []byte, err error) {
	body, r.Body, err = drainBody(r.Body)
	return
}

func responseMediaType(headers http.Header) (mediaType string) {
	mediaType, _, _ = mime.ParseMediaType(headers.Get("Content-Type"))
	return
}

func responseBodyIsJSON(headers http.Header) bool {
	return responseMediaType(headers) == "application/json"
}

func responseBodyIsXML(headers http.Header) bool {
	switch responseMediaType(headers) {
	case "application/xml", "text/xml":
		return true
	}
	return false
}

func responseIsSuccess(status int) bool {
	return status >= 200 && status < 300
}

// jsonUnmarshalReader reads all the content from the given reader r and
// unmarshals the data as JSON into the value v.
func jsonUnmarshalReader(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// xmlUnmarshalReader reads all the content from the given reader r and
// unmarshals the data as XML into the value v.
func xmlUnmarshalReader(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// xmlMarshalReader returns the XML encoding of v as a io.Reader.
func xmlMarshalReader(v any) (r io.Reader, err error) {
	var b bytes.Buffer
	if _, err := b.WriteString(xml.Header); err != nil {
		return nil, err
	}
	enc := xml.NewEncoder(&b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return &b, nil
}

func buildURL(base *url.URL, refUrl string, query url.Values) (string, error) {
	u, err := base.Parse(refUrl)
	if err != nil {
		return "", err
	}

	// TODO: support adding query params to and URL already containing query
	// params.
	u.RawQuery = query.Encode()
	return u.String(), nil
}
