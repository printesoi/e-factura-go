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

package etransport_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	itest "github.com/printesoi/e-factura-go/internal/test/etransport"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAuthCodeURL(t *testing.T) {
	assert := assert.New(t)

	cfg, err := itest.SetupTestEnvOAuth2Config(true)
	if !assert.NoError(err) {
		return
	}
	if cfg == nil {
		t.Skipf("Skipping test, no credentials to create oauth2 config")
	}

	uuid, err := uuid.NewRandom()
	if !assert.NoError(err) {
		return
	}

	url := cfg.AuthCodeURL(uuid.String())
	assert.True(url != "")
	fmt.Printf("%s\n", url)
}

func TestExchangeCode(t *testing.T) {
	assert := assert.New(t)

	cfg, err := itest.SetupTestEnvOAuth2Config(true)
	if !assert.NoError(err) {
		return
	}
	if cfg == nil {
		t.Skipf("Skipping test, no credentials to create oauth2 config")
	}

	code := os.Getenv("ETRANSPORT_TEST_EXCHANGE_CODE")
	if code == "" {
		t.Skipf("Skipping test, no exchange code")
	}

	token, err := cfg.Exchange(context.Background(), code)
	if !assert.NoError(err) {
		return
	}

	tj, _ := json.Marshal(token)
	fmt.Printf("%s\n", string(tj))
}
