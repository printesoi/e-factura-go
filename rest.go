package efactura

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type (
	Code     string
	Standard string

	ErrorMessage struct {
		Message string `json:"message"`
	}

	// ValidationResponse is the response of the validate XML endpoint
	ValidationResponse struct {
		Stare    Code           `json:"stare"`
		TraceID  string         `json:"trace_id"`
		Messages []ErrorMessage `json:"Messages,omitempty"`
	}

	// GeneratePdfResponseError is the error response of the Xml-To-Pdf
	// endpoint
	GeneratePdfResponseError struct {
		Stare    Code           `json:"stare"`
		TraceID  string         `json:"trace_id"`
		Messages []ErrorMessage `json:"Messages,omitempty"`
	}

	GeneratePdfResponse struct {
		Error *GeneratePdfResponseError
		PDF   []byte
	}
)

const (
	CodeOk  Code = "ok"
	CodeNok Code = "nok"

	StandardFACT1 Standard = "FACT1"
	StandardFCN   Standard = "FCN"
)

func (r ValidationResponse) IsOk() bool {
	return r.Stare == CodeOk
}

func (r GeneratePdfResponse) IsOk() bool {
	return r.Error == nil
}

// ValidateXML call the validate endpoint with the given standard and xml body
func (c *Client) ValidateXML(ctx context.Context, xml []byte, st Standard) (ValidationResponse, error) {
	var response ValidationResponse

	path := fmt.Sprintf(webserviceAppPathValidate, st)
	url, err := buildParseUrl(webserviceSpBaseProd, path, nil)
	if err != nil {
		return response, err
	}
	req, err := newRequest(ctx, http.MethodPost, url, xml)
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.apiClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return response, err
	}
	if !responseBodyIsJSON(resp.Header) {
		return response, newErrorResponse(resp,
			fmt.Errorf("expected application/json, got %s", responseMediaType(resp.Header)))
	}
	if err := jsonUnmarshalReader(resp.Body, &response); err != nil {
		return response, newErrorResponse(resp,
			fmt.Errorf("failed to decode JSON body: %v", err))
	}

	return response, nil
}

// ValidateInvoice validate the provided Invoice
func (c *Client) ValidateInvoice(ctx context.Context, invoice Invoice) (ValidationResponse, error) {
	xmlInvoice, err := xml.Marshal(invoice)
	if err != nil {
		return ValidationResponse{}, err
	}

	return c.ValidateXML(ctx, xmlInvoice, StandardFACT1)
}

// XmlToPdf convert the given XML to PDF. To check if the generation is indeed
// successful and no validation or other invalid request error occured, check
// if response.IsOk() == true.
func (c *Client) XmlToPdf(ctx context.Context, xml []byte, st Standard, noValidate bool) (response GeneratePdfResponse, err error) {
	path := fmt.Sprintf(webserviceAppPathXmlToPdf, st)
	if noValidate {
		path, _ = url.JoinPath(path, "DA")
	}
	url, er := buildParseUrl(webserviceSpBaseProd, path, nil)
	if err = er; err != nil {
		return
	}
	req, er := newRequest(ctx, http.MethodPost, url, xml)
	if err = er; err != nil {
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	body, _, headers, er := c.doApi(req)
	if err = er; err != nil {
		return
	}

	// If the response content type is application/json, then the validation
	// failed, otherwise we got the PDF in response body
	if responseBodyIsJSON(headers) {
		response.Error = new(GeneratePdfResponseError)
		if err = json.Unmarshal(body, response.Error); err != nil {
			return
		}
	}

	response.PDF = body
	return
}

// InvoiceToPdf convert the given Invoice to PDF. See XmlToPdf for return
// values.
func (c *Client) InvoiceToPdf(ctx context.Context, invoice Invoice, noValidate bool) (response GeneratePdfResponse, err error) {
	xmlInvoice, er := xml.Marshal(invoice)
	if err = er; err != nil {
		return
	}

	return c.XmlToPdf(ctx, xmlInvoice, StandardFACT1, noValidate)
}
