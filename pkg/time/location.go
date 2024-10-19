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

package time

import "time"

var (
	// RoZoneLocation is the Romanian timezone location loaded in the init
	// function. This library does NOT load the time/tzdata package for the
	// embedded timezone database, so the user of this library is responsible
	// to ensure the Europe/Bucharest location is available, otherwise UTC is
	// used and may lead to unexpected results.
	RoZoneLocation *time.Location

	// Allow mocking and testing
	timeNow = time.Now
)

// Now returns time.Now() in Romanian zone location (Europe/Bucharest)
func Now() time.Time {
	return timeNow().In(RoZoneLocation)
}

// TimeInRomania returns the time t in Romanian zone location
// (Europe/Bucharest).
func TimeInRomania(t time.Time) time.Time {
	return t.In(RoZoneLocation)
}

// Date time.Date in RoZoneLocation.
func Date(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, RoZoneLocation)
}

// ParseInRomania time.ParseInLocation(layout, value, RoZoLocation)
func ParseInRomania(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, RoZoneLocation)
}

func init() {
	if loc, err := time.LoadLocation("Europe/Bucharest"); err == nil {
		RoZoneLocation = loc
	} else {
		// If we could not load the Europe/Bucharest location, fallback to
		// time.UTC
		RoZoneLocation = time.UTC
	}
}
