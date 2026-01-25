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

package vat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	ierrors "github.com/printesoi/e-factura-go/internal/errors"
	"github.com/printesoi/e-factura-go/pkg/errors"
	"github.com/printesoi/e-factura-go/pkg/types"
)

var funcNow = time.Now

const (
	// Current maximum number of CIFs per request
	VatRequestCIFLimit = 100

	apiPathCheckVatV9 = "/api/PlatitorTvaRest/v9/tva"
)

// CIF represents a Romanian numeric VAT ID ("Cod de Identificare Fiscală)
type CIF int64

// MakeCIF just wraps an int64 into a CIF
func MakeCIF(n int64) CIF {
	return CIF(n)
}

// MakeCIFFromString parses a string as a numeric CIF. If the input string has
// the "RO" prefix (eg you stored in a database), this strips it.
func MakeCIFFromString(s string) (CIF, error) {
	s = strings.TrimPrefix(s, "RO")
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return MakeCIF(n), nil
}

// String returns the string representation of the CIF, useful in context where
// a string is needed.
func (c CIF) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// NullDate is a wrapper over the types.Date as the v9 TVA API does not return
// a JSON null, but rather an empty string.
type NullDate struct {
	types.Date
	Valid bool
}

func makeNullDate(d types.Date) NullDate {
	return NullDate{Date: d, Valid: true}
}

func (d *NullDate) UnmarshalJSON(data []byte) error {
	var str string
	d.Valid = false
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str = strings.TrimSpace(str); str == "" {
		return nil
	}

	date, err := types.MakeDateFromString(str)
	if err != nil {
		return err
	}

	d.Date, d.Valid = date, true
	return nil
}

func (d NullDate) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return json.Marshal(d.Date)
	}
	return json.Marshal("")
}

type CompanyGeneralData struct {
	CIF                      CIF        `json:"cui"`
	Date                     types.Date `json:"data"`
	Name                     string     `json:"denumire"`
	Address                  string     `json:"adresa"`
	RegComNo                 string     `json:"nrRegCom"`
	Phone                    string     `json:"telefon"`
	Fax                      string     `json:"fax"`
	PostalCode               string     `json:"codPostal"`
	AuthorizationDocument    string     `json:"act"`
	RegistrationState        string     `json:"stare_inregistrare"`
	RegistrationDate         NullDate   `json:"data_inregistrare"`
	CAEN                     string     `json:"cod_CAEN"`
	IBAN                     string     `json:"iban"`
	ROeFacturaStatus         bool       `json:"statusRO_e_Factura"`
	CompetentFiscalAuthority string     `json:"organFiscalCompetent"`
	OwnershipForm            string     `json:"forma_de_proprietate"`
	OrganizationalForm       string     `json:"forma_organizare"`
	LegalForm                string     `json:"forma_juridica"`
}

type CompanyVatPeriod struct {
	StartDate  NullDate `json:"data_inceput_ScpTVA"`
	EndDate    NullDate `json:"data_sfarsit_ScpTVA"`
	CancelDate NullDate `json:"data_anul_imp_ScpTVA"`
}

// Company data about registration for VAT purposes ("plătitor în scopuri de tva")
type CompanyVatRegistration struct {
	Vat        bool               `json:"scpTVA"`
	VatPeriods []CompanyVatPeriod `json:"perioade_TVA"`
	VatMessage string             `json:"mesaj_ScpTVA,omitempty"`
}

// Company data about split VAT payment ("plata defalcată a TVA")
type CompanySplitVatRegistration struct {
	StartDate NullDate `json:"dataInceputSplitTVA"`
	EndDate   NullDate `json:"dataAnulareSplitTVA"`
	Status    bool     `json:"statusSplitTVA"`
}

// Company data about cash-accounting VAT scheme ("sistemul TVA la încasare")
type CompanyVatCashVatSchemeRegistration struct {
	StartDate       NullDate `json:"dataInceputTvaInc"`
	EndDate         NullDate `json:"dataSfarsitTvaInc"`
	LastUpdateDate  NullDate `json:"dataActualizareTvaInc"`
	PublicationDate NullDate `json:"dataPublicareTvaInc"`
	DocumentType    string   `json:"tipActTvaInc"`
	Status          bool     `json:"statusTvaIncasare"`
}

// Refers to company data from the ANAF register of inactive or reactivated taxpayers
// "Registrul contribuabililor inactivi sau reactivați" -> https://www.anaf.ro/inactivi/
type CompanyInactiveState struct {
	InactivationDate   NullDate `json:"dataInactivare"`
	ReactivationDate   NullDate `json:"dataReactivare"`
	PublicationDate    NullDate `json:"dataPublicare"`
	DeregistrationDate NullDate `json:"dataRadiere"`
	InactiveStatus     bool     `json:"statusInactivi"`
}

type CompanyRegisteredOfficeAddress struct {
	StreetName     string `json:"sdenumire_Strada"`
	StreetNo       string `json:"snumar_Strada"`
	LocalityName   string `json:"sdenumire_Localitate"`
	LocalityCode   string `json:"scod_Localitate"`
	County         string `json:"sdenumire_Judet"`
	CountyCode     string `json:"scod_Judet"`
	CountyAutoCode string `json:"scod_JudetAuto"`
	Country        string `json:"stara"`
	AdressDetails  string `json:"sdetalii_Adresa"`
	PostalCode     string `json:"scod_Postal"`
}

type CompanyFiscalResidenceAddress struct {
	StreetName     string `json:"ddenumire_Strada"`
	StreetNo       string `json:"dnumar_Strada"`
	LocalityName   string `json:"ddenumire_Localitate"`
	LocalityCode   string `json:"dcod_Localitate"`
	County         string `json:"ddenumire_Judet"`
	CountyCode     string `json:"dcod_Judet"`
	CountyAutoCode string `json:"dcod_JudetAuto"`
	Country        string `json:"dtara"`
	AdressDetails  string `json:"ddetalii_Adresa"`
	PostalCode     string `json:"dcod_Postal"`
}

type CompanyData struct {
	GeneralData                  CompanyGeneralData                  `json:"date_generale"`
	VatRegistration              CompanyVatRegistration              `json:"inregistrare_scop_Tva"`
	VatCashVatSchemeRegistration CompanyVatCashVatSchemeRegistration `json:"inregistrare_RTVAI"`
	InactiveState                CompanyInactiveState                `json:"stare_inactiv"`
	SplitVatRegistration         CompanySplitVatRegistration         `json:"inregistrare_SplitTVA"`
	RegisteredOfficeAddress      CompanyRegisteredOfficeAddress      `json:"adresa_sediu_social"`
	FiscalResidenceAddress       CompanyFiscalResidenceAddress       `json:"adresa_domiciliu_fiscal"`
}

func (d CompanyData) GetCIF() CIF {
	return d.GeneralData.CIF
}

func (d CompanyData) HasVat() bool {
	return d.VatRegistration.Vat
}

func (d CompanyData) GetLastVatPeriod() (p CompanyVatPeriod, ok bool) {
	if len(d.VatRegistration.VatPeriods) == 0 {
		return p, false
	}

	return d.VatRegistration.VatPeriods[0], true
}

func (d CompanyData) GetVatEnrollDate() (date NullDate) {
	if latestVatPeriod, ok := d.GetLastVatPeriod(); ok {
		return latestVatPeriod.StartDate
	}

	return
}

func (d CompanyData) GetVatEndDate() (date NullDate) {
	if latestVatPeriod, ok := d.GetLastVatPeriod(); ok {
		return latestVatPeriod.EndDate
	}

	return
}

type CheckVatResponse struct {
	Code    int    `json:"cod"`
	Message string `json:"message"`
	// Found contains the company data for the valid CIFs.
	Found []CompanyData `json:"found"`
	// NotFound contains the CIFs that were not found by the API.
	NotFound []CIF `json:"notFound"`
}

type CheckVatRequestItem struct {
	CIF  CIF        `json:"cui"`
	Date types.Date `json:"data"`
}

// MakeCheckVatRequestItem is a convenience method that creates a
// CheckVatRequestItem for the given cif and date if you prefer the method call
// syntax more than the literal one.
func MakeCheckVatRequestItem(cif CIF, date types.Date) CheckVatRequestItem {
	return CheckVatRequestItem{CIF: cif, Date: date}
}

// CheckVatRequest represents a request to the ANAF TVA v9 endpoint.
type CheckVatRequest []CheckVatRequestItem

// MakeCheckVatRequest is a convenience method that creates CheckVatRequest for
// the provided cifs for the current date. For specifying a different date than
// current date, check MakeCheckVatRequestFromItems and AddCIFWithDate.
func MakeCheckVatRequest(cifs ...CIF) CheckVatRequest {
	r := make(CheckVatRequest, len(cifs))
	for _, cif := range cifs {
		r.AddCIF(cif)
	}
	return r
}

// MakeCheckVatRequestFromItems wraps items in a CheckVatRequest.
func MakeCheckVatRequestFromItems(items ...CheckVatRequestItem) CheckVatRequest {
	return CheckVatRequest(items)
}

func (r *CheckVatRequest) append(i CheckVatRequestItem) *CheckVatRequest {
	if r == nil {
		return nil
	}

	*r = append(*r, i)
	return r
}

func (r *CheckVatRequest) len() int {
	if r == nil {
		return 0
	}

	return len(*r)
}

// AddCIF add an item to the request for the given cif and current date. If you
// need to check a different date, see AddCIFWithDate.
func (r *CheckVatRequest) AddCIF(cif CIF) *CheckVatRequest {
	return r.append(CheckVatRequestItem{
		CIF:  cif,
		Date: types.MakeDateFromTime(funcNow()),
	})
}

// AddCIFWithDate add an item to the request for the given cif and date.
func (r *CheckVatRequest) AddCIFWithDate(cif CIF, date types.Date) *CheckVatRequest {
	return r.append(CheckVatRequestItem{
		CIF:  cif,
		Date: date,
	})
}

// CheckVatV9 check VAT status using the v9 TVA API for the CIFs in the
// checkVatRequest. The request must have at least a CIF and must not have more
// than 100 CIFs (otherwise a CheckVatValidationError error is return an no
// actual call is made).
func (c *Client) CheckVatV9(ctx context.Context, checkVatRequest CheckVatRequest) (response *CheckVatResponse, err error) {
	if checkVatRequest.len() == 0 {
		return nil, NewCheckVatValidationErrorf("empty request, expected a least a CIF")
	} else if checkVatRequest.len() > VatRequestCIFLimit {
		return nil, NewCheckVatValidationErrorf("request too big, expected at most %d CIFs, got %d", VatRequestCIFLimit, checkVatRequest.len())
	}

	body, err := json.Marshal(checkVatRequest)
	if err != nil {
		return
	}

	req, err := c.PublicApiClient.NewRequest(ctx, http.MethodPost, apiPathCheckVatV9, nil, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res := new(CheckVatResponse)
	if err = c.PublicApiClient.DoUnmarshalJSON(req, res, nil); err == nil {
		response = res
	} else {
		// ANAF returns HTTP status 404 if none of the CIFs are found, but 200
		// is at least one of them is found. We try to normalize this
		// behaviour, so don't return an error on 404 and parse the response
		// and leave the task to the caller.
		var errResp *errors.ErrorResponse
		if ierrors.As(err, &errResp) && errResp.StatusCode == 404 {
			if err := json.Unmarshal(errResp.ResponseBody, &res); err != nil {
				return nil, err
			}
			if len(res.NotFound) == len(checkVatRequest) && len(res.Found) == 0 {
				return res, nil
			}
		}
	}
	return
}
