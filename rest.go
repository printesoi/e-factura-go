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
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/m29h/xml"
)

type (
	Code             string
	ValidateStandard string
	UploadStandard   string

	// ValidateResponse is the parsed response from the validate XML endpoint
	ValidateResponse struct {
		State    Code   `json:"stare"`
		TraceID  string `json:"trace_id"`
		Messages []struct {
			Message string `json:"message"`
		} `json:"Messages,omitempty"`
	}

	// GeneratePDFResponseError is the error response of the XML-To-PDF
	// endpoint
	GeneratePDFResponseError struct {
		State    Code   `json:"stare"`
		TraceID  string `json:"trace_id"`
		Messages []struct {
			Message string `json:"message"`
		} `json:"Messages,omitempty"`
	}

	// GeneratePDFResponse is the parsed response from the XML-To-PDF endpoint
	GeneratePDFResponse struct {
		Error *GeneratePDFResponseError
		PDF   []byte
	}

	MessageRASP struct {
		UploadIndex int    `xml:"index_incarcare,attr"`
		Message     string `xml:"message,attr"`
	}

	// UploadResponse is the parsed response from the upload endpoint
	UploadResponse struct {
		ResponseDate    string `xml:"dateResponse,attr,omitempty"`
		ExecutionStatus *int   `xml:"ExecutionStatus,attr,omitempty"`
		UploadIndex     *int64 `xml:"index_incarcare,attr,omitempty"`
		Errors          []struct {
			ErrorMessage string `xml:"errorMessage,attr"`
		} `xml:"Errors,omitempty"`

		XMLName xml.Name `xml:"header"`
	}

	GetMessageStateCode string

	// GetMessageStateResponse is the parsed response from the get message
	// state endoint
	GetMessageStateResponse struct {
		State      GetMessageStateCode `xml:"stare,attr"`
		DownloadID int64               `xml:"id_descarcare,attr,omitempty"`
		Errors     []struct {
			ErrorMessage string `xml:"errorMessage,attr"`
		} `xml:"Errors,omitempty"`

		XMLName xml.Name `xml:"header"`
	}

	MessageFilterType int

	Message struct {
		CreationDate string `json:"data_creare"`
		CIF          string `json:"cif"`
		UploadIndex  string `json:"id_solicitare"`
		Details      string `json:"detalii"`
		Type         string `json:"tip"`
		ID           string `json:"id"`
		CIFSeller    string `json:"cif_emitent"`
		CIFCustomer  string `json:"cif_beneficiar"`
	}

	// MessagesListResponse is the parsed response from the list messages
	// endpoint.
	MessagesListResponse struct {
		Error    string    `json:"eroare"`
		Title    string    `json:"titlu"`
		Serial   string    `json:"serial"`
		CUI      string    `json:"cui"`
		Messages []Message `json:"mesaje"`
	}

	// MessagesListPaginationResponse is the parsed response from the list
	// messages with pagination endpoint.
	MessagesListPaginationResponse struct {
		MessagesListResponse

		RecordsInPage       int64 `json:"numar_inregistrari_in_pagina"`
		TotalRecordsPerPage int64 `json:"numar_total_inregistrari_per_pagina"`
		TotalRecords        int64 `json:"numar_total_inregistrari"`
		TotalPages          int64 `json:"numar_total_pagini"`
		CurrentPageIndex    int64 `json:"index_pagina_curenta"`
	}

	// DownloadInvoiceResponseError is the error response from the download
	// invoice endpoint.
	DownloadInvoiceResponseError struct {
		Error string `json:"eroare"`
		Title string `json:"titlu,omitempty"`
	}

	// DownloadInvoiceResponse is the parsed response from the download invoice
	// endpoint.
	DownloadInvoiceResponse struct {
		Error *DownloadInvoiceResponseError
		Zip   []byte
	}
)

const (
	CodeOk  Code = "ok"
	CodeNok Code = "nok"

	ValidateStandardFACT1 ValidateStandard = "FACT1"
	ValidateStandardFCN   ValidateStandard = "FCN"

	GetMessageStateCodeOk         GetMessageStateCode = "ok"
	GetMessageStateCodeNok        GetMessageStateCode = "nok"
	GetMessageStateCodeInvalidXML GetMessageStateCode = "XML cu erori nepreluat de sistem"
	GetMessageStateCodeProcessing GetMessageStateCode = "in prelucrare"

	UploadStandardUBL  UploadStandard = "UBL"
	UploadStandardCN   UploadStandard = "CN"
	UploadStandardCII  UploadStandard = "CII"
	UploadStandardRASP UploadStandard = "RASP"

	MessageFilterAll MessageFilterType = iota
	MessageFilterErrors
	MessageFilterSent
	MessageFilterReceived
	MessageFilterCustomerMessage
)

func (s ValidateStandard) String() string {
	return string(s)
}

func (s UploadStandard) String() string {
	return string(s)
}

// IsOk returns true if the validate response was successful.
func (r *ValidateResponse) IsOk() bool {
	return r != nil && r.State == CodeOk
}

// IsOk returns true if the XML-To-PDF response was successful.
func (r *GeneratePDFResponse) IsOk() bool {
	return r != nil && r.Error == nil
}

// IsOk returns true if the response corresponding to an upload was successful.
func (r *UploadResponse) IsOk() bool {
	return r != nil && r.ExecutionStatus != nil && *r.ExecutionStatus == 0
}

// Returns the upload index (should only be called when IsOk() == true).
func (r *UploadResponse) GetUploadIndex() int64 {
	if r == nil || r.UploadIndex == nil {
		return 0
	}
	return *r.UploadIndex
}

// IsOk returns true if the message state if ok (processed, and can be
// downloaded).
func (r *GetMessageStateResponse) IsOk() bool {
	return r != nil && r.State == GetMessageStateCodeOk
}

// GetDownloadID returns the download ID (should only be called when IsOk() ==
// true).
func (r *GetMessageStateResponse) GetDownloadID() int64 {
	if r == nil {
		return 0
	}
	return r.DownloadID
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

func (t MessageFilterType) String() string {
	switch t {
	case MessageFilterErrors:
		return "E"
	case MessageFilterSent:
		return "T"
	case MessageFilterReceived:
		return "P"
	case MessageFilterCustomerMessage:
		return "R"
	}
	return ""
}

func (m Message) IsError() bool {
	return m.Type == "ERORI FACTURA"
}

func (m Message) IsSentInvoice() bool {
	return m.Type == "FACTURA TRIMISA"
}

func (m Message) IsReceivedInvoice() bool {
	return m.Type == "FACTURA PRIMITA"
}

func (m Message) IsBuyerMessage() bool {
	return m.Type == "MESAJ CUMPARATOR PRIMIT / MESAJ CUMPARATOR TRANSMIS"
}

// IsOk returns true if the response corresponding to a download was successful.
func (r *DownloadInvoiceResponse) IsOk() bool {
	return r != nil && r.Error == nil
}

// IsOk returns true if the response corresponding to fetching messages list
// was successful.
func (r *MessagesListResponse) IsOk() bool {
	return r != nil && (r.Error == "" || strings.HasPrefix(r.Error, "Nu exista mesaje in ultimele "))
}

// ValidateXML call the validate endpoint with the given standard and XML body
func (c *Client) ValidateXML(ctx context.Context, xml io.Reader, st ValidateStandard) (*ValidateResponse, error) {
	var response *ValidateResponse

	path := fmt.Sprintf(webserviceAppPathValidate, st)
	req, err := c.newApiPublicRequest(ctx, http.MethodPost, path, nil, xml)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if !responseBodyIsJSON(resp.Header) {
		return nil, newErrorResponse(resp,
			fmt.Errorf("expected application/json, got %s", responseMediaType(resp.Header)))
	}

	response = new(ValidateResponse)
	if err := jsonUnmarshalReader(resp.Body, response); err != nil {
		return nil, newErrorResponse(resp,
			fmt.Errorf("failed to decode JSON body: %v", err))
	}

	return response, nil
}

// ValidateInvoice validate the provided Invoice
func (c *Client) ValidateInvoice(ctx context.Context, invoice Invoice) (*ValidateResponse, error) {
	xmlReader, err := xmlMarshalReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.ValidateXML(ctx, xmlReader, ValidateStandardFACT1)
}

// XMLToPDF convert the given XML to PDF. To check if the generation is indeed
// successful and no validation or other invalid request error occured, check
// if response.IsOk() == true.
func (c *Client) XMLToPDF(ctx context.Context, xml io.Reader, st ValidateStandard, noValidate bool) (response *GeneratePDFResponse, err error) {
	path := fmt.Sprintf(webserviceAppPathXMLToPDF, st)
	if noValidate {
		path, _ = url.JoinPath(path, "DA")
	}
	req, er := c.newApiPublicRequest(ctx, http.MethodPost, path, nil, xml)
	if err = er; err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, er := c.do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err = er; err != nil {
		return
	}

	// If the response content type is application/json, then the validation
	// failed, otherwise we got the PDF in response body
	switch mediaType := responseMediaType(resp.Header); mediaType {
	case "application/json":
		response = &GeneratePDFResponse{
			Error: &GeneratePDFResponseError{},
		}
		if err = jsonUnmarshalReader(resp.Body, response.Error); err != nil {
			err = newErrorResponse(resp,
				fmt.Errorf("failed to unmarshal response body: %v", err))
			return
		}
	case "application/pdf":
		response = &GeneratePDFResponse{}
		if response.PDF, err = io.ReadAll(resp.Body); err != nil {
			err = newErrorResponse(resp,
				fmt.Errorf("failed to read body: %v", err))
			return
		}
	default:
		err = newErrorResponse(resp,
			fmt.Errorf("expected application/json or application/pdf, got %s", mediaType))
	}
	return
}

// InvoiceToPDF convert the given Invoice to PDF. See XMLToPDF for return
// values.
func (c *Client) InvoiceToPDF(ctx context.Context, invoice Invoice, noValidate bool) (response *GeneratePDFResponse, err error) {
	xmlReader, err := xmlMarshalReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.XMLToPDF(ctx, xmlReader, ValidateStandardFACT1, noValidate)
}

func ptrfyString(s string) *string {
	return &s
}

type uploadOptions struct {
	extern      *string
	autofactura *string
}

type uploadOption func(*uploadOptions)

// UploadOptionForeign is an upload option specifiying that the buyer is not a
// Romanian entity (no CUI or NIF).
func UploadOptionForeign() uploadOption {
	return func(o *uploadOptions) {
		o.extern = ptrfyString("DA")
	}
}

// UploadOptionSelfBilled is an upload option specifying that it's a
// self-billed invoice (the buyer is issuing the invoice on behalf of the
// supplier.
func UploadOptionSelfBilled() uploadOption {
	return func(o *uploadOptions) {
		o.autofactura = ptrfyString("DA")
	}
}

// UploadXML uploads and invoice or message XML. Optional upload options can be
// provided via call params.
func (c *Client) UploadXML(
	ctx context.Context, xml io.Reader, st UploadStandard, cif string, opts ...uploadOption,
) (response *UploadResponse, err error) {

	uploadOptions := uploadOptions{}
	for _, opt := range opts {
		opt(&uploadOptions)
	}

	query := url.Values{
		"standard": {st.String()},
		"cif":      {cif},
	}
	if uploadOptions.autofactura != nil {
		query.Set("extern", *uploadOptions.autofactura)
	}
	if uploadOptions.extern != nil {
		query.Set("extern", *uploadOptions.extern)
	}

	req, er := c.newApiRequest(ctx, http.MethodPost, apiPathUpload, query, xml)
	if err = er; err != nil {
		return
	}

	response = new(UploadResponse)
	err = c.doApiUnmarshalXML(req, response)
	return
}

// UploadInvoice uploads the given Invoice with the provided optional options.
func (c *Client) UploadInvoice(
	ctx context.Context, invoice Invoice, cif string, opts ...uploadOption,
) (response *UploadResponse, err error) {
	xmlReader, err := xmlMarshalReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.UploadXML(ctx, xmlReader, UploadStandardUBL, cif, opts...)
}

func (m MessageRASP) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// for stripping the custom MarshalXML method
	type messageRASP MessageRASP
	var xmlMsg struct {
		messageRASP

		XMLName xml.Name `xml:"header"`
	}

	xmlMsg.messageRASP = messageRASP(m)
	xmlMsg.XMLName.Space = XMLNSANAFreqMesajv1

	return e.EncodeElement(xmlMsg, start)
}

// UploadRASPMessage uploads the given MessageRASP with the provided optional
// options.
func (c *Client) UploadRASPMessage(
	ctx context.Context, msg MessageRASP, cif string, opts ...uploadOption,
) (response *UploadResponse, err error) {
	xmlReader, err := xmlMarshalReader(msg)
	if err != nil {
		return nil, err
	}
	return c.UploadXML(ctx, xmlReader, UploadStandardRASP, cif, opts...)
}

// GetMessageState fetch the state of a message. The uploadIndex must a result
// from an upload operation.
func (c *Client) GetMessageState(
	ctx context.Context, uploadIndex int,
) (response *GetMessageStateResponse, err error) {
	query := url.Values{
		"id_incarcare": {strconv.Itoa(uploadIndex)},
	}
	req, er := c.newApiRequest(ctx, http.MethodGet, apiPathMessageState, query, nil)
	if err = er; err != nil {
		return
	}

	response = new(GetMessageStateResponse)
	err = c.doApiUnmarshalXML(req, response)
	return
}

// GetMessages fetches the list of messages for a provided cif, number of days
// and a filter. For fetching all messages use MessageFilterAll as the value
// for msgType.
func (c *Client) GetMessagesList(
	ctx context.Context, cif string, numDays int, msgType MessageFilterType,
) (response *MessagesListResponse, err error) {
	query := url.Values{
		"cif":  {cif},
		"zile": {strconv.Itoa(numDays)},
	}
	if msgType != MessageFilterAll {
		query.Set("filter", msgType.String())
	}
	req, er := c.newApiRequest(ctx, http.MethodGet, apiPathMessageList, query, nil)
	if err = er; err != nil {
		return
	}

	response = new(MessagesListResponse)
	err = c.doApiUnmarshalJSON(req, response)
	return
}

// GetMessagesListPagination fetches the list of messages for a provided cif,
// start time (unix time in milliseconds), end time (unix time in milliseconds)
// and a filter. For fetching all messages use MessageFilterAll as the value
// for msgType.
func (c *Client) GetMessagesListPagination(
	ctx context.Context, cif string, startTs, endTs int64, msgType MessageFilterType,
) (response *MessagesListPaginationResponse, err error) {
	query := url.Values{
		"cif":       {cif},
		"startTime": {strconv.FormatInt(startTs, 10)},
		"endTime":   {strconv.FormatInt(endTs, 10)},
	}
	if msgType != MessageFilterAll {
		query.Set("filter", msgType.String())
	}

	req, er := c.newApiRequest(ctx, http.MethodGet, apiPathMessagePaginationList, query, nil)
	if err = er; err != nil {
		return
	}

	err = c.doApiUnmarshalXML(req, &response)
	return
}

// DownloadInvoice download an invoice zip for a given download index
func (c *Client) DownloadInvoice(
	ctx context.Context, downloadID int,
) (response *DownloadInvoiceResponse, err error) {
	query := url.Values{
		"id": {strconv.Itoa(downloadID)},
	}
	req, er := c.newApiRequest(ctx, http.MethodGet, apiPathDownload, query, nil)
	if err = er; err != nil {
		return
	}

	resp, er := c.do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err = er; err != nil {
		return
	}

	// If the response content type is application/json, then the download
	// failed, otherwise we got the zip in response body
	switch mediaType := responseMediaType(resp.Header); mediaType {
	case "application/json":
		response = &DownloadInvoiceResponse{
			Error: &DownloadInvoiceResponseError{},
		}
		if err = jsonUnmarshalReader(resp.Body, response.Error); err != nil {
			err = newErrorResponse(resp,
				fmt.Errorf("failed to unmarshal response body: %v", err))
			return
		}
	case "application/zip":
		response = &DownloadInvoiceResponse{}
		if response.Zip, err = io.ReadAll(resp.Body); err != nil {
			err = newErrorResponse(resp,
				fmt.Errorf("failed to read body: %v", err))
			return
		}
	default:
		err = newErrorResponse(resp,
			fmt.Errorf("expected application/json or application/pdf, got %s", mediaType))
	}
	return
}
