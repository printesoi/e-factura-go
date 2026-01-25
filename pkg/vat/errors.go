// Copyright 2026 Victor Dodon
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

package vat

import "fmt"

// CheckVatValidationError is a simple wrapper over a standard error. It's
// useful to distinguish validation errors before the actual call to ANAF APIs
// are made.
type CheckVatValidationError struct {
	error
}

// NewCheckVatValidationErrorf wraps a fmt.Errorf into a CheckVatValidationError
func NewCheckVatValidationErrorf(format string, a ...any) *CheckVatValidationError {
	return &CheckVatValidationError{error: fmt.Errorf(format, a...)}
}

// NewCheckVatValidationError wraps a standard error into a CheckVatValidationError
func NewCheckVatValidationError(err error) *CheckVatValidationError {
	return &CheckVatValidationError{error: err}
}

func (e *CheckVatValidationError) Error() string {
	return e.error.Error()
}

func (e *CheckVatValidationError) Unwrap() error {
	return e.error
}
