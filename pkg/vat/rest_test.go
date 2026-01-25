// Copyright 2026 Victor Dodon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License
package vat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/printesoi/e-factura-go/internal/errors"
	"github.com/printesoi/e-factura-go/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestCIF(t *testing.T) {
	assert := assert.New(t)

	type test struct {
		input       string
		expectError bool
		expected    CIF
	}
	for _, test := range []test{
		{"", true, CIF(0)},
		{"ANAF", true, CIF(0)},
		{"RO", true, CIF(0)},
		{"0", false, CIF(0)},
		{"1234", false, CIF(1234)},
		{"RO1234", false, CIF(1234)},
	} {
		cif, err := MakeCIFFromString(test.input)
		if test.expectError {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(test.expected, cif)
		}
	}
}

func TestCheckVatValidation(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Use dummy handler to ensure no actual calls to ANAF APIs are executed.
	client, teardown, err := setupTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	if teardown != nil {
		defer teardown()
	}
	assert.NoError(err)

	{
		// It should return a CheckVatValidationError for an empty request.
		_, err := client.CheckVatV9(ctx, CheckVatRequest{})
		assert.Error(err)

		var validationError *CheckVatValidationError
		assert.True(errors.As(err, &validationError))
	}
	{
		// It should return a CheckVatValidationError if the number of CIFs is
		// > VatRequestCIFLimit
		r := CheckVatRequest{}
		for i := int64(0); i < VatRequestCIFLimit+1; i++ {
			r.AddCIF(MakeCIF(i))
		}
		_, err := client.CheckVatV9(ctx, r)
		assert.Error(err)

		var validationError *CheckVatValidationError
		assert.True(errors.As(err, &validationError))
	}
}

func TestUnmarshalResponse(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	writeJSON := func(w http.ResponseWriter, data any, status int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		b, _ := json.Marshal(data)
		w.Write(b)
	}

	getCifs := func(items []CheckVatRequestItem) []CIF {
		cifs := make([]CIF, 0, len(items))
		for _, item := range items {
			cifs = append(cifs, item.CIF)
		}
		return cifs
	}

	mockCompanyData := func(item CheckVatRequestItem) CompanyGeneralData {
		companyData := CompanyGeneralData{}
		companyData.CIF = item.CIF
		companyData.Date = item.Date
		companyData.Name = fmt.Sprintf("Company %d SRL", item.CIF)
		companyData.OwnershipForm = "PROPR.PRIVATA-CAPITAL PRIVAT AUTOHTON"
		companyData.OrganizationalForm = "PERSOANA JURIDICA"
		companyData.LegalForm = "SOCIETATE COMERCIALĂ CU RĂSPUNDERE LIMITATĂ"
		companyData.ROeFacturaStatus = true
		companyData.CAEN = "6201"
		companyData.Address = "București"
		companyData.CompetentFiscalAuthority = "Administraţia Sector 1 a Finanţelor Publice"
		registrationDate := types.MakeDate(2000, time.January, 1)
		companyData.RegistrationDate = makeNullDate(registrationDate)
		companyData.RegistrationState = fmt.Sprintf("INREGISTRAT din data %s", registrationDate.Format("2006.01.02"))
		return companyData
	}

	setupTest := func(items []CheckVatRequestItem, handler func(
		w http.ResponseWriter,
		r *http.Request,
		checkVatRequest CheckVatRequest,
	)) (client *Client, teardown func(), err error) {
		// Use dummy handler to ensure no actual calls to ANAF APIs are executed.
		client, teardown, err = setupTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate path & method
			assert.Equal(apiPathCheckVatV9, r.RequestURI)
			assert.Equal(http.MethodPost, r.Method)

			// Validate valid JSON body
			var payload CheckVatRequest
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Validate payload is expected
			if !assert.Equal(len(items), len(payload)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			for i, item := range payload {
				if !assert.Equal(items[i].CIF, item.CIF) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if !assert.True(items[i].Date.Equal(item.Date)) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			// Additional actions
			handler(w, r, payload)
		}))
		return
	}

	{
		// All CIFs not found

		items := []CheckVatRequestItem{
			{
				CIF:  CIF(1234),
				Date: types.MakeDate(2026, time.January, 12),
			},
			{
				CIF:  CIF(5678),
				Date: types.MakeDate(2026, time.January, 23),
			},
		}

		// Test HTTP status 404 when none of the CIFs are found.
		client, teardown, err := setupTest(items, func(w http.ResponseWriter, r *http.Request, payload CheckVatRequest) {
			response := CheckVatResponse{}
			for _, item := range payload {
				response.NotFound = append(response.NotFound, item.CIF)
			}

			writeJSON(w, response, http.StatusNotFound)
		})
		if teardown != nil {
			defer teardown()
		}
		assert.NoError(err)

		response, err := client.CheckVatV9(ctx, MakeCheckVatRequestFromItems(items...))
		if assert.NoError(err) {
			assert.Equal(getCifs(items), response.NotFound)
			assert.Empty(response.Found)
		}
	}
	{
		// All CIFs found

		items := []CheckVatRequestItem{
			{
				CIF:  CIF(1234),
				Date: types.MakeDate(2026, time.January, 12),
			},
			{
				CIF:  CIF(5678),
				Date: types.MakeDate(2026, time.January, 23),
			},
		}

		// Test HTTP status 404 when none of the CIFs are found.
		client, teardown, err := setupTest(items, func(w http.ResponseWriter, r *http.Request, payload CheckVatRequest) {
			response := CheckVatResponse{}
			for _, item := range items {
				data := CompanyData{
					GeneralData: mockCompanyData(item),
				}
				response.Found = append(response.Found, data)
			}

			writeJSON(w, response, http.StatusOK)
		})
		if teardown != nil {
			defer teardown()
		}
		assert.NoError(err)

		response, err := client.CheckVatV9(ctx, MakeCheckVatRequestFromItems(items...))
		if assert.NoError(err) {
			assert.Empty(response.NotFound)

			assert.Equal(len(items), len(response.Found))
			for i, data := range response.Found {
				companyData := CompanyData{
					GeneralData: mockCompanyData(items[i]),
				}
				assert.Equal(companyData, data)
			}
		}

	}
	{
		// Mix of found and not found

		foundItems := []CheckVatRequestItem{
			{
				CIF:  CIF(1234),
				Date: types.MakeDate(2026, time.January, 12),
			},
			{
				CIF:  CIF(5678),
				Date: types.MakeDate(2026, time.January, 23),
			},
		}
		notFoundItems := []CheckVatRequestItem{
			{
				CIF:  CIF(9876),
				Date: types.MakeDate(2026, time.January, 1),
			},
		}
		allItems := append(foundItems, notFoundItems...)

		client, teardown, err := setupTest(allItems, func(w http.ResponseWriter, r *http.Request, payload CheckVatRequest) {
			response := CheckVatResponse{}
			for _, item := range foundItems {
				data := CompanyData{
					GeneralData: mockCompanyData(item),
				}
				response.Found = append(response.Found, data)
			}
			for _, item := range notFoundItems {
				response.NotFound = append(response.NotFound, item.CIF)
			}

			writeJSON(w, response, http.StatusOK)
		})
		if teardown != nil {
			defer teardown()
		}
		assert.NoError(err)

		response, err := client.CheckVatV9(ctx, MakeCheckVatRequestFromItems(allItems...))
		if assert.NoError(err) {
			assert.Equal(getCifs(notFoundItems), response.NotFound)

			assert.Equal(len(foundItems), len(response.Found))
			for i, data := range response.Found {
				assert.Equal(foundItems[i].CIF, data.GeneralData.CIF)
				assert.True(foundItems[i].Date.Equal(data.GeneralData.Date))
			}
		}
	}
}
