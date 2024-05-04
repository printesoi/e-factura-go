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
	"encoding/json"
	"testing"
	_ "time/tzdata"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMessage(t *testing.T) {
	assert := assert.New(t)

	type messageTest struct {
		JSON []byte

		ExpectedMessage             Message
		ExpectedID                  int64
		ExpectedUploadIndex         int64
		ExpectedIsError             bool
		ExpectedIsSentInvoice       bool
		ExpectedIsReceivedInvoice   bool
		ExpectedIsBuyerMessage      bool
		ExpectedIsSelfBilledInvoice bool
		ExpectedSellerCIF           string
		ExpectedBuyerCIF            string
	}
	tests := []messageTest{
		{
			JSON: []byte(`{}`),
		},
		{
			JSON: []byte(`{
				"cif": "123456789",
				"data_creare": "202401020304",
				"detalii": "Erori de validare identificate la factura transmisa cu id_incarcare=42",
				"id": "128",
				"id_solicitare": "42",
				"tip": "ERORI FACTURA"
			}`),
			ExpectedMessage: Message{
				ID:           "128",
				Type:         MessageTypeError,
				UploadIndex:  "42",
				CIF:          "123456789",
				Details:      "Erori de validare identificate la factura transmisa cu id_incarcare=42",
				CreationDate: "202401020304",
			},
			ExpectedID:          int64(128),
			ExpectedUploadIndex: int64(42),
			ExpectedIsError:     true,
		},
		{
			JSON: []byte(`{
				"cif": "123456789",
				"data_creare": "202401020304",
				"detalii": "Erori de validare identificate la factura de tip declarat=AUTOFACTURA, transmisa cu id_incarcare=42",
				"id": "128",
				"id_solicitare": "42",
				"tip": "ERORI FACTURA"
			}`),
			ExpectedMessage: Message{
				ID:           "128",
				Type:         MessageTypeError,
				UploadIndex:  "42",
				CIF:          "123456789",
				Details:      "Erori de validare identificate la factura de tip declarat=AUTOFACTURA, transmisa cu id_incarcare=42",
				CreationDate: "202401020304",
			},
			ExpectedID:                  int64(128),
			ExpectedUploadIndex:         int64(42),
			ExpectedIsError:             true,
			ExpectedIsSelfBilledInvoice: true,
		},
		{
			JSON: []byte(`{
				"cif": "123456789",
				"data_creare": "202401020304",
				"detalii": "Factura cu id_incarcare=42 emisa de cif_emitent=123456789 pentru cif_beneficiar=987654321",
				"id": "128",
				"id_solicitare": "42",
				"tip": "FACTURA TRIMISA"
			}`),
			ExpectedMessage: Message{
				ID:           "128",
				Type:         MessageTypeSentInvoice,
				UploadIndex:  "42",
				CIF:          "123456789",
				Details:      "Factura cu id_incarcare=42 emisa de cif_emitent=123456789 pentru cif_beneficiar=987654321",
				CreationDate: "202401020304",
			},
			ExpectedID:            int64(128),
			ExpectedUploadIndex:   int64(42),
			ExpectedIsSentInvoice: true,
			ExpectedSellerCIF:     "123456789",
			ExpectedBuyerCIF:      "987654321",
		},
		{
			JSON: []byte(`{
				"cif": "123456789",
				"data_creare": "202401020304",
				"detalii": "Factura cu id_incarcare=42 emisa de cif_emitent=987654321 pentru cif_beneficiar=123456789",
				"id": "128",
				"id_solicitare": "42",
				"tip": "FACTURA PRIMITA"
			}`),
			ExpectedMessage: Message{
				ID:           "128",
				Type:         MessageTypeReceivedInvoice,
				UploadIndex:  "42",
				CIF:          "123456789",
				Details:      "Factura cu id_incarcare=42 emisa de cif_emitent=987654321 pentru cif_beneficiar=123456789",
				CreationDate: "202401020304",
			},
			ExpectedID:                int64(128),
			ExpectedUploadIndex:       int64(42),
			ExpectedIsReceivedInvoice: true,
			ExpectedSellerCIF:         "987654321",
			ExpectedBuyerCIF:          "123456789",
		},
		{
			JSON: []byte(`{
				"cif": "123456789",
				"data_creare": "202401020304",
				"detalii": "Factura cu id_incarcare=42 transmisa de cif=123456789  ca autofactutra in numele cif=987654321",
				"id": "128",
				"id_solicitare": "42",
				"tip": "FACTURA PRIMITA"
			}`),
			ExpectedMessage: Message{
				ID:           "128",
				Type:         MessageTypeReceivedInvoice,
				UploadIndex:  "42",
				CIF:          "123456789",
				Details:      "Factura cu id_incarcare=42 transmisa de cif=123456789  ca autofactutra in numele cif=987654321",
				CreationDate: "202401020304",
			},
			ExpectedID:                  int64(128),
			ExpectedUploadIndex:         int64(42),
			ExpectedIsReceivedInvoice:   true,
			ExpectedIsSelfBilledInvoice: true,
			ExpectedSellerCIF:           "987654321",
			ExpectedBuyerCIF:            "123456789",
		},
		// TODO: we also need to check self billed invoices that were issued on
		// our behalf.
	}
	for _, mt := range tests {
		var um Message
		if !assert.NoError(json.Unmarshal(mt.JSON, &um)) {
			continue
		}

		// Test raw fields
		assert.Equal(mt.ExpectedMessage.ID, um.ID)
		assert.Equal(mt.ExpectedMessage.Type, um.Type)
		assert.Equal(mt.ExpectedMessage.UploadIndex, um.UploadIndex)
		assert.Equal(mt.ExpectedMessage.CIF, um.CIF)
		assert.Equal(mt.ExpectedMessage.Details, um.Details)
		assert.Equal(mt.ExpectedMessage.CreationDate, um.CreationDate)

		// Test getters/helpers
		assert.Equal(mt.ExpectedID, um.GetID())
		assert.Equal(mt.ExpectedUploadIndex, um.GetUploadIndex())
		assert.Equal(mt.ExpectedIsError, um.IsError())
		assert.Equal(mt.ExpectedIsSentInvoice, um.IsSentInvoice())
		assert.Equal(mt.ExpectedIsReceivedInvoice, um.IsReceivedInvoice())
		assert.Equal(mt.ExpectedIsBuyerMessage, um.IsBuyerMessage())
		assert.Equal(mt.ExpectedIsSelfBilledInvoice, um.IsSelfBilledInvoice())
		assert.Equal(mt.ExpectedSellerCIF, um.GetSellerCIF())
		assert.Equal(mt.ExpectedBuyerCIF, um.GetBuyerCIF())
		if mt.ExpectedMessage.CreationDate != "" {
			mdate, ok := um.GetCreationDate()
			assert.True(ok)
			assert.Equal(mt.ExpectedMessage.CreationDate, mdate.Format(messageTimeLayout))
		}
	}
}
