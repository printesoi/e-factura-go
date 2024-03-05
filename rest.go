package efactura

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type (
	Code             string
	ValidateStandard string
	UploadStandard   string

	ErrorMessage struct {
		Message string `json:"message"`
	}

	// ValidationResponse is the response of the validate XML endpoint
	ValidationResponse struct {
		State    Code           `json:"stare"`
		TraceID  string         `json:"trace_id"`
		Messages []ErrorMessage `json:"Messages,omitempty"`
	}

	// GeneratePdfResponseError is the error response of the Xml-To-Pdf
	// endpoint
	GeneratePdfResponseError struct {
		State    Code           `json:"stare"`
		TraceID  string         `json:"trace_id"`
		Messages []ErrorMessage `json:"Messages,omitempty"`
	}

	GeneratePdfResponse struct {
		Error *GeneratePdfResponseError
		PDF   []byte
	}

	MessageRASP struct {
		UploadIndex int    `xml:"index_incarcare,attr"`
		Message     string `xml:"message,attr"`
	}

	ErrorMessageNode struct {
		ErrorMessage string `xml:"errorMessage,attr"`
	}

	// UploadResponse is a parsed response from the upload endpoint
	UploadResponse struct {
		ResponseDate    string             `xml:"dateResponse,attr,omitempty"`
		ExecutionStatus int                `xml:"ExecutionStatus,attr,omitempty"`
		UploadIndex     int                `xml:"index_incarcare,attr,omitempty"`
		Errors          []ErrorMessageNode `xml:"Errors,omitempty"`

		XMLName xml.Name `xml:"header"`
		XMLNS   string   `xml:"xmlns,attr"`
	}

	// GetMessageStateResponse is a parsed response from the get message state
	// endoint
	GetMessageStateResponse struct {
		State Code `json:"stare"`
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

	MessagesListResponse struct {
		Serial   string    `json:"serial"`
		CUI      string    `json:"cui"`
		Title    string    `json:"titlu"`
		Messages []Message `json:"message"`
	}

	MessagesListPaginationResponse struct {
		Serial   string    `json:"serial"`
		CUI      string    `json:"cui"`
		Title    string    `json:"titlu"`
		Messages []Message `json:"message"`

		RecordsInPage       int64 `json:"numar_inregistrari_in_pagina"`
		TotalRecordsPerPage int64 `json:"numar_total_inregistrari_per_pagina"`
		TotalRecords        int64 `json:"numar_total_inregistrari"`
		TotalPages          int64 `json:"numar_total_pagini"`
		CurrentPageIndex    int64 `json:"index_pagina_curenta"`
	}

	DownloadInvoiceResponseError struct {
		Error string `json:"eroare"`
		Title string `json:"titlu,omitempty"`
	}

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

func (r ValidationResponse) IsOk() bool {
	return r.State == CodeOk
}

func (r GeneratePdfResponse) IsOk() bool {
	return r.Error == nil
}

// IsOk returns true if the response corresponding to an upload was successful.
func (r UploadResponse) IsOk() bool {
	return r.ExecutionStatus == 0
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

// IsOk returns true if the response corresponding to a download was successful.
func (r DownloadInvoiceResponse) IsOk() bool {
	return r.Error == nil
}

// ValidateXML call the validate endpoint with the given standard and xml body
func (c *Client) ValidateXML(ctx context.Context, xml io.Reader, st ValidateStandard) (*ValidationResponse, error) {
	var response *ValidationResponse

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

	response = new(ValidationResponse)
	if err := jsonUnmarshalReader(resp.Body, response); err != nil {
		return nil, newErrorResponse(resp,
			fmt.Errorf("failed to decode JSON body: %v", err))
	}

	return response, nil
}

// ValidateInvoice validate the provided Invoice
func (c *Client) ValidateInvoice(ctx context.Context, invoice Invoice) (*ValidationResponse, error) {
	xmlReader, err := xmlMarshalReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.ValidateXML(ctx, xmlReader, ValidateStandardFACT1)
}

// XmlToPdf convert the given XML to PDF. To check if the generation is indeed
// successful and no validation or other invalid request error occured, check
// if response.IsOk() == true.
func (c *Client) XmlToPdf(ctx context.Context, xml io.Reader, st ValidateStandard, noValidate bool) (response *GeneratePdfResponse, err error) {
	path := fmt.Sprintf(webserviceAppPathXmlToPdf, st)
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
		response = &GeneratePdfResponse{
			Error: &GeneratePdfResponseError{},
		}
		if err = jsonUnmarshalReader(resp.Body, response.Error); err != nil {
			err = newErrorResponse(resp,
				fmt.Errorf("failed to unmarshal response body: %v", err))
			return
		}
	case "application/pdf":
		response = &GeneratePdfResponse{}
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

// InvoiceToPdf convert the given Invoice to PDF. See XmlToPdf for return
// values.
func (c *Client) InvoiceToPdf(ctx context.Context, invoice Invoice, noValidate bool) (response *GeneratePdfResponse, err error) {
	xmlReader, err := xmlMarshalReader(invoice)
	if err != nil {
		return nil, err
	}

	return c.XmlToPdf(ctx, xmlReader, ValidateStandardFACT1, noValidate)
}

func ptrfyString(s string) *string {
	return &s
}

type uploadOptions struct {
	extern      *string
	autofactura *string
}

type uploadOption func(*uploadOptions)

func UploadOptionExtern() uploadOption {
	return func(o *uploadOptions) {
		o.extern = ptrfyString("DA")
	}
}

func UploadOptionAutofactura() uploadOption {
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
		XMLNS   string   `xml:"xmlns,attr"`
	}

	xmlMsg.messageRASP = messageRASP(m)
	xmlMsg.XMLNS = XMLNSANAFreqMesajv1

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
	err = c.doApiUnmarshalXML(req, response)
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
	ctx context.Context, downloadIndex int,
) (response *DownloadInvoiceResponse, err error) {
	query := url.Values{
		"id": {strconv.Itoa(downloadIndex)},
	}
	req, er := c.newApiRequest(ctx, http.MethodGet, apiPathMessagePaginationList, query, nil)
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
