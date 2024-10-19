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

// Package contants provides some useful constants that are not bound to any
// specific API.
package constants

const (
	// ApiBaseURL is the base URL for the ANAF protected services.
	ApiBaseURL = "https://api.anaf.ro/"

	ApiBasePathSandbox = "test/"
	// ApiBaseSandbox points to the sandbox (testing) version of the APIs.
	ApiBaseSandbox = ApiBaseURL + ApiBasePathSandbox

	ApiBasePathProd = "prod/"
	// ApiBaseProd points to the production version of the APIs.
	ApiBaseProd = ApiBaseURL + ApiBasePathProd

	// PublicApiBaseURL is the base URL for the ANAF public (unprotected) APIs.
	PublicApiBaseURL = "https://webservicesp.anaf.ro/"

	PublicApiBasePathProd = "prod/"
	// PublicApiBaseProd points to the production version of the public
	// (unprotected) APIs.
	PublicApiBaseProd = PublicApiBaseURL + PublicApiBasePathProd
)
