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
	"strconv"
)

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
