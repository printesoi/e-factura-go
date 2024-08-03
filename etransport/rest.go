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
	ixml "github.com/printesoi/e-factura-go/xml"
)

const (
	apiBase             = "ETRANSPORT/ws/v1/"
	apiPathUploadV2     = apiBase + "upload/%s/%s/2"
	apiPathMessageList  = apiBase + "lista/%d/%s"
	apiPathMessageState = apiBase + "stareMesaj/%d"
	apiPathInfo         = apiBase + "info"
)

type Message struct {
	Uit          string         `json:"uit"`
	Cod_decl     string         `json:"cod_decl"`
	Ref_decl     string         `json:"ref_decl"`
	Sursa        string         `json:"sursa"`
	Id_incarcare string         `json:"id_incarcare"`
	Data_creare  string         `json:"data_creare"`
	Tip_op       string         `json:"tip_op"`
	Data_transp  string         `json:"data_transp"`
	Pc_tara      string         `json:"pc_tara"`
	Pc_cod       string         `json:"pc_cod"`
	Pc_den       string         `json:"pc_den"`
	Tr_tara      string         `json:"tr_tara"`
	Tr_cod       string         `json:"tr_cod"`
	Tr_den       string         `json:"tr_den"`
	Nr_veh       string         `json:"nr_veh"`
	Nr_rem1      string         `json:"nr_rem1"`
	Nr_rem2      string         `json:"nr_rem2"`
	Mesaje       []MessageError `json:"mesaje,omitempty"`
}

type MessageError struct {
	Type    string `json:"tip"`
	Message string `json:"mesaj"`
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

type (
	GetMessageStateCode string

	GetMessageStateResponse struct {
		Errors []struct {
			ErrorMessage string `json:"errorMessage"`
		} `json:"errors,omitempty"`
		State           GetMessageStateCode `json:"stare"`
		DateResponse    string              `json:"dateResponse"`
		ExecutionStatus int32               `json:"ExecutionStatus"`
		TraceID         string              `json:"trace_id"`
	}
)

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

type (
	uploadStandard string

	UploadV2Response struct {
		DateResponse    string `json:"dateResponse"`
		ExecutionStatus int32  `json:"ExecutionStatus"`
		UploadIndex     int64  `json:"index_incarcare"`
		UIT             string `json:"UIT"`
		TraceID         string `json:"trace_id"`
		DeclarantRef    string `json:"ref_declarant"`
		Attention       string `json:"atentie,omitempty"`
		Errors          []struct {
			ErrorMessage string `json:"errorMessage"`
		} `json:"errors,omitempty"`
	}
)

// IsOk returns true if the response corresponding to fetching messages list
// was successful.
func (r *UploadV2Response) IsOk() bool {
	return r != nil && r.ExecutionStatus == 0
}

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
	ctx context.Context, decl PostingDeclaration, cif string,
) (response *UploadV2Response, err error) {
	xmlReader, err := ixml.MarshalXMLToReader(decl)
	if err != nil {
		return nil, err
	}

	return c.UploadV2XML(ctx, xmlReader, cif)
}
