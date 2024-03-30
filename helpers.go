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
	"reflect"
	"regexp"
	"strconv"

	"github.com/printesoi/xml-go"
)

const (
	mediaTypeApplicationJSON = "application/json"
	mediaTypeApplicationXML  = "application/xml"
	mediaTypeApplicationPDF  = "application/pdf"
	mediaTypeApplicationZIP  = "application/zip"
	mediaTypeTextXML         = "text/xml"
	mediaTypeTextPlain       = "text/plain"
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
	return responseMediaType(headers) == mediaTypeApplicationJSON
}

func responseBodyIsPlainText(headers http.Header) bool {
	return responseMediaType(headers) == mediaTypeTextPlain
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
	return UnmarshalXML(data, v)
}

// setupUBLXMLEncoder will configure the xml.Encoder to make it suitable for
// marshaling UBL objects to XML.
func setupUBLXMLEncoder(enc *xml.Encoder) *xml.Encoder {
	enc.AddNamespaceBinding(XMLNSUBLcac, "cac")
	enc.AddSkipNamespaceAttrForPrefix(XMLNSUBLcac, "cac")
	enc.AddNamespaceBinding(XMLNSUBLcbc, "cbc")
	enc.AddSkipNamespaceAttrForPrefix(XMLNSUBLcbc, "cbc")
	return enc
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

func buildParseURL(baseUrl, refUrl string, query url.Values) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	return buildURL(base, refUrl, query)
}

func atoi64(s string) (n int64, ok bool) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	return i, true
}

func itoa64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func matchFirstSubmatch(input string, re *regexp.Regexp) (string, bool) {
	ms := re.FindStringSubmatch(input)
	if ms == nil || len(ms) < 2 {
		return "", false
	}
	return ms[1], true
}

func ptrfyString(s string) *string {
	return &s
}

func ptrfyStringNotEmpty(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}

func typeName(v any) string {
	return reflect.TypeOf(v).Name()
}

func typeNameAddrPtr(v any) string {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	return rt.Name()
}

func concatBytes(s ...[]byte) []byte {
	n := 0
	for _, v := range s {
		n += len(v)
	}
	res := make([]byte, 0, n)
	for _, v := range s {
		res = append(res, v...)
	}
	return res
}
