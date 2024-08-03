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
	"github.com/shopspring/decimal"

	ixml "github.com/printesoi/e-factura-go/xml"
	"github.com/printesoi/xml-go"
)

type PostingDeclaration struct {
	// CUI/CIF/CNP
	DeclarantCode string `xml:"codDeclarant,attr"`
	//
	DeclarantRef     string               `xml:"refDeclarant,attr,omitempty"`
	DeclPostIncident DeclPostIncidentType `xml:"declPostAvarie,attr,omitempty"`

	Notification  *PostingDeclarationNotification  `xml:"notificare,omitempty"`
	Deletion      *PostingDeclarationDeletion      `xml:"stergere,omitempty"`
	Confirmation  *PostingDeclarationConfirmation  `xml:"confirmare,omitempty"`
	VehicleChange *PostingDeclarationVehicleChange `xml:"modifVehicul,omitempty"`

	// Name of node.
	XMLName xml.Name `xml:"eTransport"`
	// xmlns attr. Will be automatically set in MarshalXML
	Namespace string `xml:"xmlns,attr"`
	// xmlns:xsi attr. Will be automatically set in MarshalXML
	NamespaceXSI string `xml:"xmlns:xsi,attr"`
	// generated with... Will be automatically set in MarshalXML if empty.
	Comment string `xml:",comment"`
}

// Prefill sets the  NS, NScac, NScbc and Comment properties for ensuring that
// the required attributes and properties are set for a valid UBL XML.
func (pd *PostingDeclaration) Prefill() {
	pd.Namespace = "mfp:anaf:dgti:eTransport:declaratie:v2"
	pd.NamespaceXSI = "http://www.w3.org/2001/XMLSchema-instance"
}

func (pd PostingDeclaration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type eTransport PostingDeclaration
	pd.Prefill()
	start.Name.Local = "eTransport"
	return e.EncodeElement(eTransport(pd), start)
}

// XML returns the XML encoding of the PostingDeclaration
func (pd PostingDeclaration) XML() ([]byte, error) {
	return ixml.MarshalXMLWithHeader(pd)
}

// XMLIndent works like XML, but each XML element begins on a new
// indented line that starts with prefix and is followed by one or more
// copies of indent according to the nesting depth.
func (pd PostingDeclaration) XMLIndent(prefix, indent string) ([]byte, error) {
	return ixml.MarshalIndentXMLWithHeader(pd, prefix, indent)
}

type UitType string

type PostingDeclarationNotification struct {
	OpType OpType `xml:"codTipOperatiune,attr"`

	// Cardinality: 0..1
	Correction *PostingDeclarationNotificationCorrection `xml:"corectie,omitempty"`
	// Cardinality: 1..n
	TransportedGoods []PostingDeclarationNotificationTransportedGood `xml:"bunuriTransportate"`
	// Cardinality: 1..1
	CommercialPartner PostingDeclarationNotificationCommercialPartner `xml:"partenerComercial"`
	// Cardinality: 1..1
	TransportData PostingDeclarationNotificationTransportData `xml:"dateTransport"`
	// Cardinality: 1..1
	RouteStartPlace PostingDeclationPlace `xml:"locStartTraseuRutier"`
	// Cardinality: 1..1
	RouteEndPlace PostingDeclationPlace `xml:"locFinalTraseuRutier"`
	// Cardinality: 1..n
	TransportDocuments []PostingDeclarationTransportDocument `xml:"documenteTransport"`
	// Cardinality: 0..n
	PrevNotifications []string `xml:"notificareAnterioara,omitempty"`
}

type PostingDeclarationNotificationCorrection struct {
	Uit UitType `xml:"uit,attr"`
}

type PostingDeclarationNotificationTransportedGood struct {
	OpPurposeCode   OpPurposeCodeType   `xml:"codScopOperatiune,attr"`
	TariffCode      string              `xml:"codTarifar,attr,omitempty"`
	GoodName        string              `xml:"denumireMarfa,attr"`
	Quantity        decimal.Decimal     `xml:"cantitate,attr"`
	UnitMeasureCode UnitMeasureCodeType `xml:"codUnitateMasura,attr"`
	NetWeight       *decimal.Decimal    `xml:"greutateNeta,attr,omitempty"`
	GrossWeight     decimal.Decimal     `xml:"greutateBruta,attr"`
	LeiValueNoVAT   *decimal.Decimal    `xml:"valoareLeiFaraTva,attr,omitempty"`
	DeclarantRef    string              `xml:"refDeclaranti,attr,omitempty"`
}

type PostingDeclarationNotificationCommercialPartner struct {
	CountryCode CountryCodeType `xml:"codTara,attr"`
	Code        string          `xml:"cod,attr,omitempty"`
	Name        string          `xml:"denumire,attr"`
}

type PostingDeclarationNotificationTransportData struct {
	LicensePlate            string          `xml:"nrVehicul,attr"`
	TrailerLicensePlate1    string          `xml:"nrRemorca1,attr,omitempty"`
	TrailerLicensePlate2    string          `xml:"nrRemorca2,attr,omitempty"`
	TransportOrgCountryCode CountryCodeType `xml:"codTaraOrgTransport,attr"`
	TransportOrgCode        string          `xml:"codOrgTransport,attr,omitempty"`
	TransportOrgName        string          `xml:"denumireOrgTransport,attr"`
	// TODO: type
	TransportDate string `xml:"dataTransport,attr"`
}

type PostingDeclationPlace struct {
	Location          *PostingDeclationLocation `xml:"locatie,omitempty"`
	BCPCode           BCPCodeType               `xml:"codPtf,attr,omitempty"`
	CustomsOfficeCode CustomsOfficeCodeType     `xml:"codBirouVamal,attr,omitempty"`
}

type PostingDeclationLocation struct {
	CountyCode       CountyCodeType `xml:"codJudet,attr"`
	LocalityName     string         `xml:"denumireLocalitate,attr"`
	StreetName       string         `xml:"denumireStrada,attr"`
	StreetNo         string         `xml:"numar,attr,omitempty"`
	Building         string         `xml:"bloc,attr,omitempty"`
	BuildingEntrance string         `xml:"scara,attr,omitempty"`
	Floor            string         `xml:"etaj,attr,omitempty"`
	Apartment        string         `xml:"apartament,attr,omitempty"`
	OtherInfo        string         `xml:"alteInfo,attr,omitempty"`
	ZipCode          string         `xml:"codPostal,attr,omitempty"`
}

type PostingDeclarationTransportDocument struct {
	DocumentType DocumentType `xml:"tipDocument,attr"`
	DocumentNo   string       `xml:"numarDocument,attr,omitempty"`
	DocumentDate string       `xml:"dataDocument,attr"`
	Remarks      string       `xml:"observatii,attr,omitempty"`
}

type PostingDeclarationDeletion struct {
}

type PostingDeclarationConfirmation struct {
}

type PostingDeclarationVehicleChange struct {
}
