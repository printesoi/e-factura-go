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
	"github.com/printesoi/xml-go"
	"github.com/shopspring/decimal"
)

// Decimal is a wrapper of the github.com/shopspring/decimal.Decimal type in
// order to ensure type safety and lossless computation.
type Decimal struct {
	decimal.Decimal
}

// Zero constant, to make computations faster.
// Zero should never be compared with == or != directly, please use
// decimal.Equal or decimal.Cmp instead.
var Zero Decimal = DD(decimal.Zero)

// NewFromDecimal converts a decimal.Decimal to Decimal.
func NewFromDecimal(d decimal.Decimal) Decimal {
	return Decimal{Decimal: d}
}

// DD is a synonym for NewFromDecimal.
func DD(d decimal.Decimal) Decimal {
	return NewFromDecimal(d)
}

// NewFromFloat converts a float64 to Decimal.
func NewFromFloat(f float64) Decimal {
	return NewFromDecimal(decimal.NewFromFloat(f))
}

// D is a synonym for NewFromFloat.
func D(f float64) Decimal {
	return NewFromFloat(f)
}

// NewFromString returns a new Decimal from a string representation.
// Trailing zeroes are not trimmed.
func NewFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}
	return NewFromDecimal(d), nil
}

// Ptr returns a pointer to d. Useful ins contexts where a pointer is needed.
func (d Decimal) Ptr() *Decimal {
	return &d
}

// IsInitialized if the decimal is initialized (ie is created explicitly with a
// constructor, not implicitly via var declaration).
func (d *Decimal) IsInitialized() bool {
	if d == nil {
		return false
	}
	return d.Decimal != decimal.Decimal{}
}

// Value returns the value of the pointer receiver. If the receiver is nil,
// Zero is returned.
func (d *Decimal) Value() Decimal {
	if d == nil {
		return Zero
	}
	return *d
}

// MarshalXML implements the xml.Marshaler interface.
func (d *Decimal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if d == nil {
		return nil
	}
	return e.EncodeElement(d.String(), start)
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	return DD(d.Decimal.Add(d2.Decimal))
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	return DD(d.Decimal.Sub(d2.Decimal))
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return DD(d.Decimal.Neg())
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	return DD(d.Decimal.Mul(d2.Decimal))
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	return DD(d.Decimal.Div(d2.Decimal))
}

// DivRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//
//	for a positive quotient digit 5 is rounded up, away from 0
//	if the quotient is negative then digit 5 is rounded down, away from 0
//
// Note that precision<0 is allowed as input.
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	return DD(d.Decimal.DivRound(d2.Decimal, precision))
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	return DD(d.Decimal.Mod(d2.Decimal))
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 Decimal) Decimal {
	return DD(d.Decimal.Pow(d2.Decimal))
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
// Example:
//
//	DD(decimal.NewFromString("123.456")).Truncate(2).String() // "123.45"
func (d Decimal) Truncate(precision int32) Decimal {
	return DD(d.Decimal.Truncate(precision))
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
//	NewFromFloat(5.45).Round(1).String() // output: "5.5"
//	NewFromFloat(545).Round(-1).String() // output: "550"
func (d Decimal) Round(places int32) Decimal {
	return DD(d.Decimal.Round(places))
}

// Returns the Decimal suitable to use as an amount, ie. rounds is to two
// decimal places.
func (d Decimal) AsAmount() Decimal {
	return DD(d.Decimal.Round(2))
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
func (d Decimal) Cmp(d2 Decimal) int {
	return d.Decimal.Cmp(d2.Decimal)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Cmp(d2) == 0
}
