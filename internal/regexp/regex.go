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

package regexp

import (
	"regexp"
)

// MatchString reports whether the string input
// contains any match of the regular expression re.
func MatchString(re *regexp.Regexp, input string) bool {
	return re.MatchString(input)
}

// MatchFirstSubmatch will return the first group matched. If the string
// input does not match the regex re or the number of matched groups is 0,
// the second value returned is false.
func MatchFirstSubmatch(re *regexp.Regexp, input string) (string, bool) {
	ms := re.FindStringSubmatch(input)
	if ms == nil || len(ms) < 2 {
		return "", false
	}
	return ms[1], true
}
