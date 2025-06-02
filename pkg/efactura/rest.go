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
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/printesoi/xml-go"

	ierrors "github.com/printesoi/e-factura-go/internal/errors"
	"github.com/printesoi/e-factura-go/internal/helpers"
	api_helpers "github.com/printesoi/e-factura-go/internal/helpers/api"
	iregexp "github.com/printesoi/e-factura-go/internal/regexp"
	"github.com/printesoi/e-factura-go/pkg/client"
	ptime "github.com/printesoi/e-factura-go/pkg/time"
	pxml "github.com/printesoi/e-factura-go/pkg/xml"
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

	RaspMessage struct {
		UploadIndex int64  `xml:"index_incarcare,attr"`
		Message     string `xml:"message,attr"`

		// Hardcode the namespace here so we don't need a custom marshaling
		// method.
		XMLName xml.Name `xml:"mfp:anaf:dgti:spv:reqMesaj:v1 header"`
	}

	// UploadResponse is the parsed response from the upload endpoint
	UploadResponse struct {
		ResponseDate    string `xml:"dateResponse,attr,omitempty"`
		ExecutionStatus *int   `xml:"ExecutionStatus,attr,omitempty"`
		UploadIndex     *int64 `xml:"index_incarcare,attr,omitempty"`
		Errors          []struct {
			ErrorMessage string `xml:"errorMessage,attr"`
		} `xml:"Errors,omitempty"`

		// Hardcode the namespace here so we don't need a custom marshaling
		// method.
		XMLName xml.Name `xml:"mfp:anaf:dgti:spv:respUploadFisier:v1 header"`
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

		// Hardcode the namespace here so we don't need a custom marshaling
		// method.
		XMLName xml.Name `xml:"mfp:anaf:dgti:efactura:stareMesajFactura:v1 header"`
	}

	MessageFilterType int

	Message struct {
		ID           string `json:"id"`
		Type         string `json:"tip"`
		UploadIndex  string `json:"id_solicitare"`
		CIF          string `json:"cif"`
		Details      string `json:"detalii"`
		CreationDate string `json:"data_creare"`
	}

	// MessagesListResponse is the parsed response from the list messages
	// endpoint.
	MessagesListResponse struct {
		Error    string    `json:"eroare,omitempty"`
		Title    string    `json:"titlu,omitempty"`
		Serial   string    `json:"serial,omitempty"`
		CUI      string    `json:"cui,omitempty"`
		Messages []Message `json:"mesaje,omitempty"`
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

	// DownloadInvoiceParseZipResponse is the type returned by the
	// DownloadInvoiceParseZip method. It includes the DownloadInvoiceResponse
	// (the zip archive as a []byte), the invoice and signature XML (as
	// []byte), and also a *Invoice or a *InvoiceErrorMessage (parsed Invoice
	// or InvoiceErrorMessage from InvoiceXML).
	DownloadInvoiceParseZipResponse struct {
		DownloadResponse *DownloadInvoiceResponse

		// InvoiceXML is the XML file corresponding to the
		// Invoice/InvoiceErrorMessage file from the ZIP archive.
		InvoiceXML ZipFile
		// Signature is the XML file corresponding to the Signature file from
		// the ZIP archive. This field is useful for manually parsing and
		// verifying the signature.
		SignatureXML ZipFile

		// Invoice is the parsed Invoice if the InvoiceXML is storing an
		// invoice.
		Invoice *Invoice
		// InvoiceError is the parse InvoiceErrorMessage if InvoiceXML is
		// storing an invoice error message.
		InvoiceError *InvoiceErrorMessage
	}

	// InvoiceErrorMessage is the type corresponding to an Invoice message
	// error from the download zip.
	InvoiceErrorMessage struct {
		UploadIndex int64  `xml:"Index_incarcare,attr,omitempty"`
		CIFSeller   string `xml:"Cif_emitent,attr,omitempty"`
		Errors      []struct {
			ErrorMessage string `xml:"errorMessage,attr"`
		} `xml:"Error,omitempty"`

		// Hardcode the namespace here so we don't need a custom marshaling
		// method.
		XMLName xml.Name `xml:"mfp:anaf:dgti:efactura:mesajEroriFactuta:v1 header"`
	}

	// ValidateSignatureResponse is the response returned by the validate
	// signature endpoint.
	ValidateSignatureResponse struct {
		Message string `json:"msg"`
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

	// A No-op filter that returns all messages
	MessageFilterAll MessageFilterType = iota
	// Filter that returns only the errors
	MessageFilterErrors
	// Filter that returns only the sent invoices
	MessageFilterSent
	// Filter that returns only the received invoices
	MessageFilterReceived
	// Filter that returns only the customer send or received messages
	MessageFilterBuyerMessage

	// MessageTypeError
	MessageTypeError           = "ERORI FACTURA"
	MessageTypeSentInvoice     = "FACTURA TRIMISA"
	MessageTypeReceivedInvoice = "FACTURA PRIMITA"
	MessageTypeBuyerMessage    = "MESAJ CUMPARATOR PRIMIT / MESAJ CUMPARATOR TRANSMIS"

	messageTimeLayout = "200601021504"
)

const (
	apiBase                      = "FCTEL/rest/"
	apiPathUpload                = apiBase + "upload"
	apiPathUploadB2C             = apiBase + "uploadb2c"
	apiPathMessageState          = apiBase + "stareMesaj"
	apiPathMessageList           = apiBase + "listaMesajeFactura"
	apiPathMessagePaginationList = apiBase + "listaMesajePaginatieFactura"
	apiPathDownload              = apiBase + "descarcare"
	apiPathValidateSignature     = "/api/validate/signature"

	publicApiBase         = "FCTEL/rest/"
	publicApiPathValidate = publicApiBase + "validare/%s"
	publicApiPathXMLToPDF = publicApiBase + "transformare/%s"
)

var (
	regexSellerCIF           = regexp.MustCompile("\\bcif_emitent=(\\d+)")
	regexBuyerCIF            = regexp.MustCompile("\\bcif_beneficiar=(\\d+)")
	regexErrTypeSelfBilled   = regexp.MustCompile("\\btip declarat=AUTOFACTURA\\b")
	regexTypeSelfBilled      = regexp.MustCompile(" ca autofactutra in numele cif=")
	regexSelfBilledSellerCIF = regexp.MustCompile("\\bin numele cif=(\\d+)")
	regexSelfBilledBuyerCIF  = regexp.MustCompile("\\btransmisa de cif=(\\d+)")

	regexZipFile          = regexp.MustCompile("^\\d+.xml$")
	regexZipSignatureFile = regexp.MustCompile("^semnatura_\\d+.xml$")
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

// GetFirstMessage returns the first message from the validate response. If no
// messages are set, empty string is returned.
func (r *ValidateResponse) GetFirstMessage() string {
	if r == nil || len(r.Messages) == 0 {
		return ""
	}

	return r.Messages[0].Message
}

// IsOk returns true if the XML-To-PDF response was successful.
func (r *GeneratePDFResponse) IsOk() bool {
	return r != nil && r.Error == nil
}

// GetError is a getter for the Error field.
func (r *GeneratePDFResponse) GetError() *GeneratePDFResponseError {
	if r == nil {
		return nil
	}
	return r.Error
}

// GetFirstMessage returns the first message from the validate response. If no
// messages are set, empty string is returned.
func (r *GeneratePDFResponseError) GetFirstMessage() string {
	if r == nil || len(r.Messages) == 0 {
		return ""
	}

	return r.Messages[0].Message
}

// IsOk returns true if the response corresponding to an upload was successful.
func (r *UploadResponse) IsOk() bool {
	return r != nil && r.ExecutionStatus != nil && *r.ExecutionStatus == 0
}

// GetUploadIndex returns the upload index (should only be called when
// IsOk() == true).
func (r *UploadResponse) GetUploadIndex() int64 {
	if r == nil || r.UploadIndex == nil {
		return 0
	}
	return *r.UploadIndex
}

// GetFirstErrorMessage returns the first error message. If no error messages
// are set for the upload response, empty string is returned.
func (r *UploadResponse) GetFirstErrorMessage() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}

	return r.Errors[0].ErrorMessage
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

// GetFirstErrorMessage returns the first error message. If no error messages
// are set for the response, empty string is returned.
func (r *GetMessageStateResponse) GetFirstErrorMessage() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}

	return r.Errors[0].ErrorMessage
}

func (t MessageFilterType) String() string {
	switch t {
	case MessageFilterErrors:
		return "E"
	case MessageFilterSent:
		return "T"
	case MessageFilterReceived:
		return "P"
	case MessageFilterBuyerMessage:
		return "R"
	}
	return ""
}

// IsError returns true if message type is ERORI FACTURA
func (m Message) IsError() bool {
	return m.Type == MessageTypeError
}

// IsSentInvoice returns true if message type is FACTURA TRIMISA
func (m Message) IsSentInvoice() bool {
	return m.Type == MessageTypeSentInvoice
}

// IsReceivedInvoice returns true if message type is FACTURA PRIMITA
func (m Message) IsReceivedInvoice() bool {
	return m.Type == MessageTypeReceivedInvoice
}

// IsBuyerMessage returns true if message type is MESAJ CUMPARATOR PRIMIT / MESAJ CUMPARATOR TRANSMIS
func (m Message) IsBuyerMessage() bool {
	return m.Type == MessageTypeBuyerMessage
}

// GetID parses and returns the message ID as int64 (since the API returns it
// as string).
func (m Message) GetID() int64 {
	n, _ := helpers.Atoi64(m.ID)
	return n
}

// GetUploadIndex parses and returns the upload index as int64 (since the API
// returns it as string).
func (m Message) GetUploadIndex() int64 {
	n, _ := helpers.Atoi64(m.UploadIndex)
	return n
}

// IsSelfBilledInvoice returns true if the message represents a self-billed
// invoice.
func (m Message) IsSelfBilledInvoice() bool {
	if m.IsError() {
		return regexErrTypeSelfBilled.MatchString(m.Details)
	}

	return regexTypeSelfBilled.MatchString(m.Details)
}

// GetSellerCIF parses message details and returns the seller CIF.
func (m Message) GetSellerCIF() (sellerCIF string) {
	if m.IsError() {
		return
	}
	if m.IsReceivedInvoice() {
		if m.IsSelfBilledInvoice() {
			sellerCIF, _ = iregexp.MatchFirstSubmatch(regexSelfBilledSellerCIF, m.Details)
			return
		}

		sellerCIF, _ = iregexp.MatchFirstSubmatch(regexSellerCIF, m.Details)
		return
	}
	sellerCIF = m.CIF
	return
}

// GetBuyerCIF parses message details and returns the buyer CIF.
func (m Message) GetBuyerCIF() (buyerCIF string) {
	if m.IsError() {
		return
	}
	if m.IsSelfBilledInvoice() {
		buyerCIF = m.CIF
		return
	}
	buyerCIF, _ = iregexp.MatchFirstSubmatch(regexBuyerCIF, m.Details)
	return
}

// GetCreationDate parsed CreationDate and returns a time.Time in
// RoZoneLocation.
func (m Message) GetCreationDate() (time.Time, bool) {
	t, err := ptime.ParseInRomania(messageTimeLayout, m.CreationDate)
	return t, err == nil
}

// IsOk returns true if the response corresponding to a download was successful.
func (r *DownloadInvoiceResponse) IsOk() bool {
	return r != nil && r.Error == nil
}

// IsOk returns true if the response corresponding to a download was successful.
func (r *DownloadInvoiceParseZipResponse) IsOk() bool {
	return r != nil && r.DownloadResponse.IsOk()
}

// IsOk returns true if the response corresponding to fetching messages list
// was successful.
func (r *MessagesListResponse) IsOk() bool {
	return r != nil && (r.Error == "" || strings.HasPrefix(r.Error, "Nu exista mesaje in "))
}

// ValidateXML call the validate endpoint with the given standard and XML body
// reader.
func (c *Client) ValidateXML(ctx context.Context, xml io.Reader, st ValidateStandard) (*ValidateResponse, error) {
	var response *ValidateResponse

	path := fmt.Sprintf(publicApiPathValidate, st)
	req, err := c.publicApiClient.NewRequest(ctx, http.MethodPost, path, nil, xml)
	if err != nil {
		return nil, err
	}

	// This is explicitly requested in the docs.
	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.publicApiClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if !api_helpers.ResponseBodyIsJSON(resp.Header) {
		return nil, ierrors.NewErrorResponse(resp,
			fmt.Errorf("expected %s, got %s", api_helpers.MediaTypeApplicationJSON, api_helpers.ResponseMediaType(resp.Header)))
	}

	response = new(ValidateResponse)
	if err := api_helpers.UnmarshalReaderJSON(resp.Body, response); err != nil {
		return nil, ierrors.NewErrorResponseParse(resp,
			fmt.Errorf("failed to decode JSON body: %v", err), false)
	}

	return response, nil
}

// ValidateInvoice validate the provided Invoice
func (c *Client) ValidateInvoice(ctx context.Context, invoice Invoice) (*ValidateResponse, error) {
	xmlReader, err := pxml.MarshalXMLToReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.ValidateXML(ctx, xmlReader, ValidateStandardFACT1)
}

// XMLToPDF converts the given XML to PDF. To check if the generation is indeed
// successful and no validation or other invalid request error occurred, check
// if response.IsOk() == true.
func (c *Client) XMLToPDF(ctx context.Context, xml io.Reader, st ValidateStandard, noValidate bool) (response *GeneratePDFResponse, err error) {
	path := fmt.Sprintf(publicApiPathXMLToPDF, st)
	if noValidate {
		path, _ = url.JoinPath(path, "DA")
	}
	req, er := c.publicApiClient.NewRequest(ctx, http.MethodPost, path, nil, xml)
	if err = er; err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, er := c.publicApiClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err = er; err != nil {
		return
	}

	// If the response content type is application/json, then the validation
	// failed, otherwise we got the PDF in response body
	switch mediaType := api_helpers.ResponseMediaType(resp.Header); mediaType {
	case api_helpers.MediaTypeApplicationJSON:
		resError := new(GeneratePDFResponseError)
		if err = api_helpers.UnmarshalReaderJSON(resp.Body, resError); err != nil {
			err = ierrors.NewErrorResponseParse(resp,
				fmt.Errorf("failed to unmarshal response body: %v", err), false)
			return
		}
		response = &GeneratePDFResponse{Error: resError}
	case api_helpers.MediaTypeApplicationPDF:
		response = &GeneratePDFResponse{}
		if response.PDF, err = io.ReadAll(resp.Body); err != nil {
			err = ierrors.NewErrorResponseParse(resp,
				fmt.Errorf("failed to read body: %v", err), false)
			return
		}
	default:
		err = ierrors.NewErrorResponse(resp,
			fmt.Errorf("expected %s or %s, got %s", api_helpers.MediaTypeApplicationJSON,
				api_helpers.MediaTypeApplicationPDF, mediaType))
	}
	return
}

// InvoiceToPDF convert the given Invoice to PDF. See XMLToPDF for return
// values.
func (c *Client) InvoiceToPDF(ctx context.Context, invoice Invoice, noValidate bool) (response *GeneratePDFResponse, err error) {
	xmlReader, err := pxml.MarshalXMLToReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.XMLToPDF(ctx, xmlReader, ValidateStandardFACT1, noValidate)
}

type uploadOptions struct {
	query *url.Values
	b2c   bool
}

type UploadOption func(*uploadOptions)

// UploadOptionForeign is an upload option specifying that the buyer is not a
// Romanian entity (no CUI or NIF).
func UploadOptionForeign() UploadOption {
	return func(o *uploadOptions) {
		o.query.Set("extern", "DA")
	}
}

// UploadOptionSelfBilled is an upload option specifying that it's a
// self-billed invoice (the buyer is issuing the invoice on behalf of the
// supplier.
func UploadOptionSelfBilled() UploadOption {
	return func(o *uploadOptions) {
		o.query.Set("extern", "DA")
	}
}

// UploadOptionEnforcement is an upload option specifying that the invoice it's
// uploaded by the enforcement authority on behalf of the debtor.
func UploadOptionEnforcement() UploadOption {
	return func(o *uploadOptions) {
		o.query.Set("executare", "DA")
	}
}

// UploadOptionB2C is an upload options specifying it's a B2C upload.
func UploadOptionB2C() UploadOption {
	return func(o *uploadOptions) {
		o.b2c = true
	}
}

// UploadXML uploads and invoice or message XML. Optional upload options can be
// provided via call params.
func (c *Client) UploadXML(
	ctx context.Context, xml io.Reader, st UploadStandard, cif string, opts ...UploadOption,
) (response *UploadResponse, err error) {
	query := url.Values{
		"standard": {st.String()},
		"cif":      {cif},
	}
	uploadOptions := uploadOptions{query: &query}
	for _, opt := range opts {
		opt(&uploadOptions)
	}

	path := apiPathUpload
	if uploadOptions.b2c {
		path = apiPathUploadB2C
	}
	req, er := c.apiClient.NewRequest(ctx, http.MethodPost, path, query, xml)
	if err = er; err != nil {
		return
	}

	res := new(UploadResponse)
	if err = c.apiClient.DoUnmarshalXML(req, res); err == nil {
		response = res
	}
	return
}

// UploadInvoice uploads the given Invoice with the provided optional options.
func (c *Client) UploadInvoice(
	ctx context.Context, invoice Invoice, cif string, opts ...UploadOption,
) (response *UploadResponse, err error) {
	xmlReader, err := pxml.MarshalXMLToReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.UploadXML(ctx, xmlReader, UploadStandardUBL, cif, opts...)
}

// UploadRaspMessage uploads the given RaspMessage.
func (c *Client) UploadRaspMessage(
	ctx context.Context, msg RaspMessage, cif string,
) (response *UploadResponse, err error) {
	xmlReader, err := pxml.MarshalXMLToReader(msg)
	if err != nil {
		return nil, err
	}
	return c.UploadXML(ctx, xmlReader, UploadStandardRASP, cif)
}

// GetMessageState fetch the state of a message. The uploadIndex must a result
// from an upload operation.
func (c *Client) GetMessageState(
	ctx context.Context, uploadIndex int64,
) (response *GetMessageStateResponse, err error) {
	query := url.Values{
		"id_incarcare": {strconv.FormatInt(uploadIndex, 10)},
	}
	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, apiPathMessageState, query, nil)
	if err = er; err != nil {
		return
	}

	res := new(GetMessageStateResponse)
	if err = c.apiClient.DoUnmarshalXML(req, res); err == nil {
		response = res
	}
	return
}

// GetMessagesList fetches the list of messages for a provided cif, number of days
// and a filter. For fetching all messages use MessageFilterAll as the value
// for msgType.
// NOTE: If there are no messages for the given interval, ANAF APIs
// return an error. For this case, the response.IsOk() returns true and the
// Messages slice is empty, since I don't think that no messages should result
// in an error.
func (c *Client) GetMessagesList(
	ctx context.Context, cif string, numDays int, msgType MessageFilterType,
) (response *MessagesListResponse, err error) {
	query := url.Values{
		"cif":  {cif},
		"zile": {strconv.Itoa(numDays)},
	}
	if msgType != MessageFilterAll {
		query.Set("filtru", msgType.String())
	}
	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, apiPathMessageList, query, nil)
	if err = er; err != nil {
		return
	}

	res := new(MessagesListResponse)
	if err = c.apiClient.DoUnmarshalJSON(req, res, func(r *http.Response, _ any) error {
		if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(res.Error); ok {
			return ierrors.NewLimitExceededError(r, limit, fmt.Errorf("%s: %s", res.Title, res.Error))
		}
		return nil
	}); err == nil {
		response = res
	}
	return
}

// GetMessagesListPagination fetches the list of messages for a provided cif,
// start time (unix time in milliseconds), end time (unix time in milliseconds)
// and a filter. For fetching all messages use MessageFilterAll as the value
// for msgType.
// NOTE: If there are no messages for the given interval, ANAF APIs
// return an error. For this case, the response.IsOk() returns true,
// response.TotalRecords = 0, and the Messages slice is empty, since I don't
// think that no messages should result in an error.
func (c *Client) GetMessagesListPagination(
	ctx context.Context, cif string, startTs, endTs time.Time, page int64, msgType MessageFilterType,
) (response *MessagesListPaginationResponse, err error) {
	query := url.Values{
		"cif":       {cif},
		"startTime": {helpers.Itoa64(startTs.UnixMilli())},
		"endTime":   {helpers.Itoa64(endTs.UnixMilli())},
		"pagina":    {helpers.Itoa64(page)},
	}
	if f := msgType.String(); f != "" {
		query.Set("filtru", f)
	}

	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, apiPathMessagePaginationList, query, nil)
	if err = er; err != nil {
		return
	}

	res := new(MessagesListPaginationResponse)
	if err = c.apiClient.DoUnmarshalJSON(req, res, func(r *http.Response, _ any) error {
		if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(res.Error); ok {
			return ierrors.NewLimitExceededError(r, limit, fmt.Errorf("%s: %s", res.Title, res.Error))
		}
		return nil
	}); err == nil {
		response = res
	}
	return
}

// DownloadInvoice downloads an invoice zip for a given download index.
func (c *Client) DownloadInvoice(
	ctx context.Context, downloadID int64,
) (response *DownloadInvoiceResponse, err error) {
	query := url.Values{
		"id": {strconv.FormatInt(downloadID, 10)},
	}
	req, er := c.apiClient.NewRequest(ctx, http.MethodGet, apiPathDownload, query, nil)
	if err = er; err != nil {
		return
	}

	resp, er := c.apiClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err = er; err != nil {
		return
	}

	// If the response content type is application/json, then the download
	// failed, otherwise we got the zip in response body
	switch mediaType := api_helpers.ResponseMediaType(resp.Header); mediaType {
	case api_helpers.MediaTypeApplicationJSON:
		resError := new(DownloadInvoiceResponseError)
		if err = api_helpers.UnmarshalReaderJSON(resp.Body, resError); err != nil {
			err = ierrors.NewErrorResponseParse(resp, err, false)
			return
		}
		if limit, ok := ierrors.ErrorMessageMatchLimitExceeded(resError.Error); ok {
			err = ierrors.NewLimitExceededError(resp, limit, fmt.Errorf("%s: %s", resError.Title, resError.Error))
			return
		}
		response = &DownloadInvoiceResponse{Error: resError}
	case api_helpers.MediaTypeApplicationZIP:
		response = &DownloadInvoiceResponse{}
		if response.Zip, err = io.ReadAll(resp.Body); err != nil {
			err = ierrors.NewErrorResponseParse(resp, err, false)
			return
		}
	case api_helpers.MediaTypeTextPlain:
		err = ierrors.NewErrorResponseDetectType(resp)
	default:
		err = ierrors.NewErrorResponse(resp,
			fmt.Errorf("expected %s or %s, got %s", api_helpers.MediaTypeApplicationJSON,
				api_helpers.MediaTypeApplicationPDF, mediaType))
	}
	return
}

// DownloadInvoiceParseZip same as DownloadInvoice but also parses the zip
// archive. If the response is not nil, the DownloadResponse will always be
// set. If there was an error parsing the zip archive, the response will
// contain the download response, and an error is returned. This method is not
// validating the signature.
func (c *Client) DownloadInvoiceParseZip(
	ctx context.Context, downloadID int64,
) (response *DownloadInvoiceParseZipResponse, err error) {
	dres, er := c.DownloadInvoice(ctx, downloadID)
	if er != nil {
		return nil, er
	}

	response = new(DownloadInvoiceParseZipResponse)
	response.DownloadResponse = dres
	if !dres.IsOk() {
		return
	}

	invoiceXML, signatureXML, err := ParseInvoiceZip(ctx, dres.Zip)
	if err != nil {
		return
	}

	response.InvoiceXML, response.SignatureXML = invoiceXML, signatureXML

	var invoice *Invoice
	var invoiceError *InvoiceErrorMessage
	invoice, invoiceError, err = UnmarshalDownloadedInvoiceXML(ctx, response.InvoiceXML.Data)
	if err != nil {
		return
	}

	response.Invoice, response.InvoiceError = invoice, invoiceError
	return
}

func (c *Client) doValidateSignature(
	ctx context.Context, body io.Reader, contentType string,
) (response *ValidateSignatureResponse, err error) {
	req, er := c.publicApiClient.NewRequest(ctx, http.MethodPost, apiPathValidateSignature, nil, body,
		client.RequestOptionHeader("Content-Type", contentType))
	if err = er; err != nil {
		return
	}

	eres := new(ValidateSignatureResponse)
	if err = c.publicApiClient.DoUnmarshalJSON(req, eres, func(r *http.Response, _ any) error {
		// TODO: check rate limiting
		return nil
	}); err == nil {
		response = eres
	}
	return
}

// ValidateSignatureFiles validate the detached signature applied to an invoice
// XML. This method accepts the paths corresponding the invoice XML file and to
// the detached signature XML file (these can be extracted from the invoice
// zip archive).
func (c *Client) ValidateSignatureFiles(
	ctx context.Context, invoiceXmlPath, signatureXmlPath string,
) (response *ValidateSignatureResponse, err error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	if err := client.MultipartFormFile(mw, "file", invoiceXmlPath, api_helpers.MediaTypeTextXML); err != nil {
		return nil, err
	}
	if err := client.MultipartFormFile(mw, "signature", signatureXmlPath, api_helpers.MediaTypeTextXML); err != nil {
		return nil, err
	}
	if err = mw.Close(); err != nil {
		return
	}
	return c.doValidateSignature(ctx, body, mw.FormDataContentType())
}

// ValidateSignature same as ValidateSignatureFiles but provide the binary data
// of the invoice XML and of the detached signature file (useful if obtaining
// the files with DownloadInvoiceParseZip).
func (c *Client) ValidateSignature(
	ctx context.Context, invoiceXmlData, signatureXmlData []byte,
) (response *ValidateSignatureResponse, err error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	if err := client.MultipartFormFileData(mw, "file", "invoice.xml", invoiceXmlData,
		api_helpers.MediaTypeTextXML); err != nil {
		return nil, err
	}
	if err := client.MultipartFormFileData(mw, "signature", "signature.xml", signatureXmlData,
		api_helpers.MediaTypeTextXML); err != nil {
		return nil, err
	}
	if err = mw.Close(); err != nil {
		return
	}
	return c.doValidateSignature(ctx, body, mw.FormDataContentType())
}

// ValidateSignatureZipData validate a zip archive (given by binary data)
// containing an invoice and a detached signature.
func (c *Client) ValidateSignatureZipData(
	ctx context.Context, zipData []byte,
) (response *ValidateSignatureResponse, err error) {
	invoiceXml, signatureXml, err := ParseInvoiceZip(ctx, zipData)
	if err != nil {
		return
	}
	return c.ValidateSignature(ctx, invoiceXml.Data, signatureXml.Data)
}

// ValidateSignatureZipFile validate a zip archive (given by path) containing
// an invoice and a detached signature.
func (c *Client) ValidateSignatureZipFile(
	ctx context.Context, zipPath string,
) (response *ValidateSignatureResponse, err error) {
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		return
	}
	return c.ValidateSignatureZipData(ctx, zipData)
}

// ZipFile is a parsed file from a ZIP archive downloaded from e-factura
// and stores the name of the file and data contents.
type ZipFile struct {
	Name string
	Data []byte
}

// ParseInvoiceZip parses a ZIP archive downloaded from e-factura and returns
// the invoice/error XML file and the signature XML file.
func ParseInvoiceZip(ctx context.Context, zipBody []byte) (invoiceXml, signatureXml ZipFile, err error) {
	var zr *zip.Reader
	zr, err = zip.NewReader(bytes.NewReader(zipBody), int64(len(zipBody)))
	if err != nil {
		return
	}

	if len(zr.File) != 2 {
		err = fmt.Errorf("expected exactly 2 files in the archive, got %v", len(zr.File))
		return
	}

	readAllZipFile := func(f *zip.File) ([]byte, error) {
		zof, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer zof.Close()
		return io.ReadAll(zof)
	}

	var data []byte
	for _, f := range zr.File {
		if regexZipFile.MatchString(f.Name) {
			data, err = readAllZipFile(f)
			if err != nil {
				return
			}
			invoiceXml = ZipFile{Data: data, Name: f.Name}

		} else if regexZipSignatureFile.MatchString(f.Name) {
			data, err = readAllZipFile(f)
			if err != nil {
				return
			}
			signatureXml = ZipFile{Data: data, Name: f.Name}
		}
	}

	if invoiceXml.Data == nil || signatureXml.Data == nil {
		err = fmt.Errorf("invoice archive is not complete")
		return
	}

	return
}

// UnmarshalDownloadedInvoiceXML unmarshals a downloaded invoice XML file data
// to either an Invoice or InvoiceErrorMessage.
func UnmarshalDownloadedInvoiceXML(ctx context.Context, invoiceXML []byte) (invoice *Invoice, invoiceError *InvoiceErrorMessage, err error) {
	// This is a trick for optimizing the unmarshaling: since the xml
	// can be either an Invoice or an InvoiceErrorMessage, we create a
	// struct with just an xml.Name, and based on the namespace we
	// unmarshal one or the other.
	type docName struct {
		XMLName xml.Name
	}
	var doc docName
	if err = pxml.UnmarshalXML(invoiceXML, &doc); err != nil {
		return
	}
	switch doc.XMLName.Space {
	case xmlnsUBLInvoice2:
		iv := new(Invoice)
		if err = pxml.UnmarshalXML(invoiceXML, iv); err != nil {
			return
		}
		invoice = iv

	case xmlnsMsgErrorV1:
		ie := new(InvoiceErrorMessage)
		if err = pxml.UnmarshalXML(invoiceXML, &ie); err != nil {
			return
		}
		invoiceError = ie

	default:
		err = fmt.Errorf("invalid namespace for invoice/message: %q", doc.XMLName.Space)
		return
	}

	return
}
