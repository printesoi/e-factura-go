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

package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	api_helpers "github.com/printesoi/e-factura-go/internal/helpers/api"
)

func (c *baseClient) debugRequest(req *http.Request, resp *http.Response) {
	var reqDump []byte
	var respDump []byte

	if resp != nil {
		dumpResponseBody := false
		if api_helpers.ResponseBodyIsJSON(resp.Header) || api_helpers.ResponseBodyIsXML(resp.Header) {
			dumpResponseBody = true
		}
		respDump, _ = httputil.DumpResponse(resp, dumpResponseBody)
		if resp.Request != nil {
			req = resp.Request
		}
	}
	if req != nil {
		reqDump, _ = httputil.DumpRequest(req, true)
	}

	log.Default().Writer().Write([]byte(fmt.Sprintf("Request: %s\n\nResponse: %s\n", string(reqDump), string(respDump))))
}
