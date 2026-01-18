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

package types_test

import (
	"testing"

	"github.com/printesoi/e-factura-go/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestAmountString(t *testing.T) {
	assert := assert.New(t)
	{
		d := types.D(1.2345)
		price := types.NewAmount(d)
		assert.Equal("1.23", price.String())
	}
	{
		d := types.D(1.2354)
		price := types.NewAmount(d)
		assert.Equal("1.24", price.String())
	}
	{
		d := types.D(1.23)
		price := types.NewAmount(d)
		assert.Equal("1.23", price.String())
	}
	{
		d := types.D(1)
		price := types.NewAmount(d)
		assert.Equal("1.00", price.String())
	}
}

func TestPriceAmountString(t *testing.T) {
	assert := assert.New(t)

	{
		d := types.D(1.2345)
		price := types.NewPrice(d)
		assert.Equal("1.2345", price.String())
	}
	{
		d := types.D(1.23)
		price := types.NewPrice(d)
		assert.Equal("1.23", price.String())
	}
	{
		d := types.D(4.5)
		price := types.NewPrice(d)
		assert.Equal("4.50", price.String())
	}
	{
		d := types.D(1)
		price := types.NewPrice(d)
		assert.Equal("1.00", price.String())
	}
}
