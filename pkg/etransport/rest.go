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

package etransport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	ierrors "github.com/printesoi/e-factura-go/internal/errors"
	"github.com/printesoi/e-factura-go/pkg/types"
	ixml "github.com/printesoi/e-factura-go/pkg/xml"
)

const (
	apiBase             = "ETRANSPORT/ws/v1/"
	apiPathUploadV2     = apiBase + "upload/%s/%s/2"
	apiPathMessageList  = apiBase + "lista/%d/%s"
	apiPathMessageState = apiBase + "stareMesaj/%d"
	apiPathInfo         = apiBase + "info"
)

type MessageState string

const (
	MessageStateOK  MessageState = "OK"
	MessageStateERR MessageState = "ERR"
)

type Message struct {
	UIT                          UITType         `json:"uit"`
	DeclarantCode                int64           `json:"cod_decl"`
	DeclarantRef                 string          `json:"ref_decl"`
	Source                       string          `json:"sursa"`
	UploadID                     int64           `json:"id_incarcare"`
	CreatedDate                  string          `json:"data_creare"`
	State                        MessageState    `json:"stare"`
	Op                           string          `json:"tip"`
	OpType                       int             `json:"tip_op,omitempty"`
	TransportDate                string          `json:"data_transp"`
	CommercialPartnerCountryCode CountryCodeType `json:"pc_tara,omitempty"`
	CommercialPartnerCode        string          `json:"pc_cod,omitempty"`
	CommercialPartnerName        string          `json:"pc_den,omitempty"`
	TransportOrgCountryCode      CountryCodeType `json:"tr_tara,omitempty"`
	TransportOrgCode             string          `json:"tr_cod,omitempty"`
	TransportOrgName             string          `json:"tr_den,omitempty"`
	LicensePlace                 string          `json:"nr_veh,omitempty"`
	Trailer1LicensePlate         string          `json:"nr_rem1,omitempty"`
	Trailer2LicensePlate         string          `json:"nr_rem2,omitempty"`
	Messages                     []MessageError  `json:"mesaje,omitempty"`
	GrossTotalWeight             types.Decimal   `json:"gr_tot_bruta,omitempty"`
	NetTotalWeight               types.Decimal   `json:"gr_tot_neta,omitempty"`
	TotalValue                   types.Decimal   `json:"val_tot,omitempty"`
	LineCount                    int64           `json:"nr_linii,omitempty"`
	PostIncident                 string          `json:"post_avarie,omitempty"`
}

type MessageErrorType string

const (
	MessageErrorTypeErr  MessageErrorType = "ERR"
	MessageErrorTypeWarn MessageErrorType = "WARN"
	MessageErrorTypeInfo MessageErrorType = "INFO"
)

type MessageError struct {
	Type    MessageErrorType `json:"tip"`
	Message string           `json:"mesaj"`
}

func (me MessageError) IsErr() bool {
	return me.Type == MessageErrorTypeErr
}

func (me MessageError) IsWarn() bool {
	return me.Type == MessageErrorTypeWarn
}

func (me MessageError) IsInfo() bool {
	return me.Type == MessageErrorTypeInfo
}

// MessagesListResponse is the parsed response from the list messages endpoint.
type MessagesListResponse struct {
	Errors []struct {
		ErrorMessage string `json:"errorMessage"`
	} `json:"errors,omitempty"`
	Messages        []Message `json:"mesaje,omitempty"`
	Serial          string    `json:"serial"`
	CUI             string    `json:"cui"`
	Title           string    `json:"titlu"`
	DateResponse    string    `json:"dateResponse"`
	ExecutionStatus int32     `json:"ExecutionStatus"`
	TraceID         string    `json:"trace_id"`
}

// IsOk returns true if the response corresponding to fetching messages list
// was successful.
func (r *MessagesListResponse) IsOk() bool {
	return r != nil && (len(r.Errors) == 0 || len(r.Errors) == 1 && strings.HasPrefix(r.Errors[0].ErrorMessage, "Nu exista mesaje in "))
}

// GetFirstErrorMessage returns the first error message. If no error messages
// are set for the response, empty string is returned.
func (r *MessagesListResponse) GetFirstErrorMessage() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].ErrorMessage
}

// GetMessagesList fetches the list of messages for a provided cif and number
// of days.
// NOTE: If there are no messages for the given interval, ANAF APIs
// return an error. For this case, the response.IsOk() returns true and the
// Messages slice is empty, since I don't think that no messages should result
// in an error.
func (c *Client) GetMessagesList(
	ctx context.Context, cif string, numDays int,
) (response *MessagesListResponse, err error) {
	path := fmt.Sprintf(apiPathMessageList, numDays, cif)
	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, path, nil, nil)
	if err = er; err != nil {
		return
	}

	res := new(MessagesListResponse)
	if err = c.apiClient.DoUnmarshalJSON(req, res, func(r *http.Response, _ any) error {
		for _, em := range res.Errors {
			if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(em.ErrorMessage); ok {
				return ierrors.NewLimitExceededError(r, limit, errors.New(em.ErrorMessage))
			}
		}
		return nil
	}); err == nil {
		response = res
	}
	return
}

type GetMessageStateResponse struct {
	Errors []struct {
		ErrorMessage string `json:"errorMessage"`
	} `json:"errors,omitempty"`
	State           GetMessageStateCode `json:"stare"`
	DateResponse    string              `json:"dateResponse"`
	ExecutionStatus int32               `json:"ExecutionStatus"`
	TraceID         string              `json:"trace_id"`
}

type GetMessageStateCode string

const (
	GetMessageStateCodeOk         GetMessageStateCode = "ok"
	GetMessageStateCodeNok        GetMessageStateCode = "nok"
	GetMessageStateCodeProcessing GetMessageStateCode = "in prelucrare"
	GetMessageStateCodeInvalidXML GetMessageStateCode = "XML cu erori nepreluat de sistem"
)

// IsOk returns true if the message state if ok (processed, and can be
// downloaded).
func (r *GetMessageStateResponse) IsOk() bool {
	return r != nil && r.State == GetMessageStateCodeOk
}

// IsNok returns true if the message state is nok (there was an error
// processing the invoice).
func (r *GetMessageStateResponse) IsNok() bool {
	return r != nil && r.State == GetMessageStateCodeNok
}

// IsProcessing returns true if the message state is processing.
func (r *GetMessageStateResponse) IsProcessing() bool {
	return r != nil && r.State == GetMessageStateCodeProcessing
}

// IsInvalidXML returns true if the message state is processing.
func (r *GetMessageStateResponse) IsInvalidXML() bool {
	return r != nil && r.State == GetMessageStateCodeInvalidXML
}

// GetFirstErrorMessage returns the first error message. If no error messages
// are set for the response, empty string is returned.
func (r *GetMessageStateResponse) GetFirstErrorMessage() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].ErrorMessage
}

// GetMessageState fetch the state of a message. The uploadIndex must a result
// from an upload operation.
func (c *Client) GetMessageState(
	ctx context.Context, uploadIndex int64,
) (response *GetMessageStateResponse, err error) {
	path := fmt.Sprintf(apiPathMessageState, uploadIndex)
	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, path, nil, nil)
	if err = er; err != nil {
		return
	}

	res := new(GetMessageStateResponse)
	if err = c.apiClient.DoUnmarshalJSON(req, res, func(r *http.Response, _ any) error {
		for _, em := range res.Errors {
			if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(em.ErrorMessage); ok {
				return ierrors.NewLimitExceededError(r, limit, errors.New(em.ErrorMessage))
			}
		}
		return nil
	}); err == nil {
		response = res
	}
	return
}

type UploadV2Response struct {
	DateResponse    string  `json:"dateResponse"`
	ExecutionStatus int32   `json:"ExecutionStatus"`
	UploadIndex     int64   `json:"index_incarcare"`
	UIT             UITType `json:"UIT"`
	TraceID         string  `json:"trace_id"`
	DeclarantRef    string  `json:"ref_declarant"`
	Attention       string  `json:"atentie,omitempty"`
	Errors          []struct {
		ErrorMessage string `json:"errorMessage"`
	} `json:"errors,omitempty"`
}

// IsOk returns true if the response corresponding to fetching messages list
// was successful.
func (r *UploadV2Response) IsOk() bool {
	return r != nil && r.ExecutionStatus == 0
}

// GetUploadIndex returns the upload index (should only be called when
// IsOk() == true).
func (r *UploadV2Response) GetUploadIndex() int64 {
	if r == nil {
		return 0
	}
	return r.UploadIndex
}

// GetUIT returns the UIT (should only be called when IsOk() == true).
func (r *UploadV2Response) GetUIT() UITType {
	if r == nil {
		return ""
	}
	return r.UIT
}

// GetFirstErrorMessage returns the first error message. If no error messages
// are set for the response, empty string is returned.
func (r *UploadV2Response) GetFirstErrorMessage() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].ErrorMessage
}

type uploadStandard string

const (
	uploadStandardETransp uploadStandard = "ETRANSP"
)

func (c *Client) UploadV2XML(
	ctx context.Context, xml io.Reader, cif string,
) (response *UploadV2Response, err error) {
	path := fmt.Sprintf(apiPathUploadV2, uploadStandardETransp, cif)
	req, er := c.apiClient.NewRequest(ctx, http.MethodPost, path, nil, xml)
	if err = er; err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/xml")

	res := new(UploadV2Response)
	if err = c.apiClient.DoUnmarshalJSON(req, res, func(r *http.Response, _ any) error {
		for _, em := range res.Errors {
			if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(em.ErrorMessage); ok {
				return ierrors.NewLimitExceededError(r, limit, errors.New(em.ErrorMessage))
			}
		}
		return nil
	}); err == nil {
		response = res
	}
	return
}

func (c *Client) UploadPostingDeclarationV2(
	ctx context.Context, decl PostingDeclarationV2, cif string,
) (response *UploadV2Response, err error) {
	xmlReader, err := ixml.MarshalXMLToReader(decl)
	if err != nil {
		return nil, err
	}

	return c.UploadV2XML(ctx, xmlReader, cif)
}
