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

package constants

const (
	ApiBaseURL       = "https://api.anaf.ro/"
	PublicApiBaseURL = "https://webservicesp.anaf.ro/"

	// ApiBaseSandbox points to the sandbox (testing) version of the API
	ApiBasePathSandbox = "test/"
	ApiBaseSandbox     = ApiBaseURL + ApiBasePathSandbox

	// ApiBaseProd points to the production version of the API
	ApiBasePathProd = "prod/"
	ApiBaseProd     = ApiBaseURL + ApiBasePathProd

	// ApiPublicBaseProd points to the production version of the public
	// (unprotected) API.
	PublicApiBasePathProd = "prod/"
	PublicApiBaseProd     = PublicApiBaseURL + PublicApiBasePathProd
)
