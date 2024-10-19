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
	"errors"

	"github.com/printesoi/e-factura-go/pkg/types"
	ixml "github.com/printesoi/e-factura-go/pkg/xml"
	"github.com/printesoi/xml-go"
)

const (
	nsETransportDeclV2 = "mfp:anaf:dgti:eTransport:declaratie:v2"
	nsXSI              = "http://www.w3.org/2001/XMLSchema-instance"
)

// PostingDeclarationV2 is the object that represents an e-transport posting
// declaration payload for the upload v2 endpoint.
type PostingDeclarationV2 struct {
	// CUI/CIF/CNP
	DeclarantCode string `xml:"codDeclarant,attr"`
	// Reference for the declarant
	DeclarantRef string `xml:"refDeclarant,attr,omitempty"`
	// DeclPostIncident this must be set only if the posting declaration is
	// uploaded after the transport already had place, otherwise leave this
	// empty.
	DeclPostIncident DeclPostIncidentType `xml:"declPostAvarie,attr,omitempty"`

	declarationType    postingDeclarationType
	declarationPayload any
}

type postingDeclarationType int

const (
	postingDeclarationTypeNA postingDeclarationType = iota
	postingDeclarationTypeNotification
	postingDeclarationTypeDeletion
	postingDeclarationTypeConfirmation
	postingDeclarationTypeVehicleChange
)

// SetNotification set the given PostingDeclarationNotification as the
// PostingDeclarationV2 payload.
func (pd *PostingDeclarationV2) SetNotification(notification PostingDeclarationNotification) *PostingDeclarationV2 {
	pd.declarationType = postingDeclarationTypeNotification
	pd.declarationPayload = notification
	return pd
}

// SetDeletion set the given PostingDeclarationDeletion as the
// PostingDeclarationV2 payload.
func (pd *PostingDeclarationV2) SetDeletion(deletion PostingDeclarationDeletion) *PostingDeclarationV2 {
	pd.declarationType = postingDeclarationTypeDeletion
	pd.declarationPayload = deletion
	return pd
}

// SetConfirmation set the given PostingDeclarationConfirmation as the
// PostingDeclarationV2 payload.
func (pd *PostingDeclarationV2) SetConfirmation(confirmation PostingDeclarationConfirmation) *PostingDeclarationV2 {
	pd.declarationType = postingDeclarationTypeConfirmation
	pd.declarationPayload = confirmation
	return pd
}

// SetVehicleChange set the given PostingDeclarationVehicleChange as the
// PostingDeclarationV2 payload.
func (pd *PostingDeclarationV2) SetVehicleChange(vehicleChange PostingDeclarationVehicleChange) *PostingDeclarationV2 {
	pd.declarationType = postingDeclarationTypeVehicleChange
	pd.declarationPayload = vehicleChange
	return pd
}

func (pd PostingDeclarationV2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type postingDeclaration PostingDeclarationV2
	var eTransport struct {
		postingDeclaration

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
		// generated with... comment
		Comment string `xml:",comment"`
	}

	eTransport.postingDeclaration = postingDeclaration(pd)
	eTransport.Namespace = nsETransportDeclV2
	eTransport.NamespaceXSI = nsXSI
	switch pd.declarationType {
	case postingDeclarationTypeNotification:
		notification, _ := pd.declarationPayload.(PostingDeclarationNotification)
		eTransport.Notification = &notification

	case postingDeclarationTypeDeletion:
		deletion, _ := pd.declarationPayload.(PostingDeclarationDeletion)
		eTransport.Deletion = &deletion

	case postingDeclarationTypeConfirmation:
		confirmation, _ := pd.declarationPayload.(PostingDeclarationConfirmation)
		eTransport.Confirmation = &confirmation

	case postingDeclarationTypeVehicleChange:
		confirmation, _ := pd.declarationPayload.(PostingDeclarationVehicleChange)
		eTransport.VehicleChange = &confirmation

	default:
		return errors.New("payload not set for posting declaration")
	}

	start.Name.Local = "eTransport"
	return e.EncodeElement(eTransport, start)
}

// XML returns the XML encoding of the PostingDeclarationV2
func (pd PostingDeclarationV2) XML() ([]byte, error) {
	return ixml.MarshalXMLWithHeader(pd)
}

// XMLIndent works like XML, but each XML element begins on a new
// indented line that starts with prefix and is followed by one or more
// copies of indent according to the nesting depth.
func (pd PostingDeclarationV2) XMLIndent(prefix, indent string) ([]byte, error) {
	return ixml.MarshalIndentXMLWithHeader(pd, prefix, indent)
}

type UITType string

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
	PrevNotifications []PostingDeclarationNotificationPrevNotification `xml:"notificareAnterioara,omitempty"`
}

type PostingDeclarationNotificationCorrection struct {
	UIT UITType `xml:"uit,attr"`
}

type PostingDeclarationNotificationTransportedGood struct {
	OpPurposeCode   OpPurposeCodeType   `xml:"codScopOperatiune,attr"`
	TariffCode      string              `xml:"codTarifar,attr,omitempty"`
	GoodName        string              `xml:"denumireMarfa,attr"`
	Quantity        types.Decimal       `xml:"cantitate,attr"`
	UnitMeasureCode UnitMeasureCodeType `xml:"codUnitateMasura,attr"`
	NetWeight       *types.Decimal      `xml:"greutateNeta,attr,omitempty"`
	GrossWeight     types.Decimal       `xml:"greutateBruta,attr"`
	LeiValueNoVAT   *types.Decimal      `xml:"valoareLeiFaraTva,attr,omitempty"`
	DeclarantRef    string              `xml:"refDeclaranti,attr,omitempty"`
}

type PostingDeclarationNotificationCommercialPartner struct {
	CountryCode CountryCodeType `xml:"codTara,attr"`
	Code        string          `xml:"cod,attr,omitempty"`
	Name        string          `xml:"denumire,attr"`
}

type PostingDeclarationNotificationTransportData struct {
	LicensePlate            string          `xml:"nrVehicul,attr"`
	Trailer1LicensePlate    string          `xml:"nrRemorca1,attr,omitempty"`
	Trailer2LicensePlate    string          `xml:"nrRemorca2,attr,omitempty"`
	TransportOrgCountryCode CountryCodeType `xml:"codTaraOrgTransport,attr"`
	TransportOrgCode        string          `xml:"codOrgTransport,attr,omitempty"`
	TransportOrgName        string          `xml:"denumireOrgTransport,attr"`
	TransportDate           types.Date      `xml:"dataTransport,attr"`
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
	DocumentDate types.Date   `xml:"dataDocument,attr"`
	Remarks      string       `xml:"observatii,attr,omitempty"`
}

type PostingDeclarationNotificationPrevNotification struct {
	UIT          UITType `xml:"uit,attr"`
	Remarks      string  `xml:"observatii,attr,omitempty"`
	DeclarantRef string  `xml:"refDeclarant,attr,omitempty"`
}

type PostingDeclarationDeletion struct {
	UIT UITType `xml:"uit,attr"`
}

type PostingDeclarationConfirmation struct {
	UIT              UITType          `xml:"uit,attr"`
	ConfirmationType ConfirmationType `xml:"tipConfirmare,attr"`
	Remarks          string           `xml:"observatii,attr,omitempty"`
}

type PostingDeclarationVehicleChange struct {
	UIT                  UITType        `xml:"uit,attr"`
	LicensePlate         string         `xml:"nrVehicul,attr"`
	Trailer1LicensePlate string         `xml:"nrRemorca1,attr,omitempty"`
	Trailer2LicensePlate string         `xml:"nrRemorca2,attr,omitempty"`
	ChangeDate           types.DateTime `xml:"dataModificare,attr"`
	Remarks              string         `xml:"observatii,attr,omitempty"`
}
