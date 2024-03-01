package efactura

import (
	"encoding/xml"
)

type Invoice struct {
	// ID: BT-1
	// Term: Numărul facturii
	// Description: O identificare unică a facturii.
	// Cardinality: 1..1
	ID string `xml:"cbc:ID"`
	// ID: BT-2
	// Term: Data emiterii facturii
	// Description: Data la care a fost emisă factura.
	// Cardinality: 1..1
	IssueDate Date `xml:"cbc:IssueDate"`
	// ID: BT-9
	// Term: Data scadenţei
	// Description: Data până la care trebuie făcută plata.
	// Cardinality: 0..1
	DueDate *Date `xml:"cbc:DueDate,omitempty"`
	// ID: BT-3
	// Term: Codul tipului facturii
	// Description: Un cod care specifică tipul funcţional al facturii.
	// Cardinality: 1..1
	InvoiceTypeCode InvoiceTypeCodeType `xml:"cbc:InvoiceTypeCode"`
	// ID: BT-5
	// Term: Codul monedei facturii
	// Description: Moneda în care sunt exprimate toate sumele din factură,
	//    cu excepţia sumei totale a TVA care este în moneda de contabilizare.
	// Cardinality: 1..1
	InvoiceCurrencyCode CurrencyCodeType `xml:"cbc:DocumentCurrencyCode"`
	// ID: BT-6
	// Term: Codul monedei de contabilizare a TVA
	// Description: Moneda utilizată pentru contabilizarea şi declararea TVA
	//     aşa cum se acceptă sau se cere în ţara Vânzătorului.
	// Cardinality: 0..1
	TaxCurrencyCode CurrencyCodeType `xml:"cbc:TaxCurrencyCode,omitempty"`
	// ID: BT-10
	// Term: Referinţa Cumpărătorului
	// Description: Un identificator atribuit de către Cumpărător utilizat
	//     pentru circuitul intern al facturii.
	// Cardinality: 0..1
	BuyerReference string `xml:"cbc:BuyerReference,omitempty"`
	// Conditional / Free-form text pertinent to this document, conveying
	// information that is not contained explicitly in other structures.
	// Cardinality: 0..1
	Note string `xml:"cbc:Note,omitempty"`
	// ID: BG-14
	// Term: Perioada de facturare
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre perioada de facturare.
	// Cardinality: 0..1
	InvoicePeriod *InvoicePeriod `xml:"cac:InvoicePeriod,omitempty"`
	// ID: BG-4
	// Term: VÂNZĂTOR
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre Vânzător.
	// Cardinality: 1..1
	Supplier InvoiceSupplier `xml:"cac:AccountingSupplierParty>cac:Party"`
	// ID: BG-7
	// Term: CUMPĂRĂTOR
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre Cumpărător.
	// Cardinality: 1..1
	Customer InvoiceCustomer `xml:"cac:AccountingCustomerParty>cac:Party"`
	// ID: BG-16
	// Term: INSTRUCŢIUNI DE PLATĂ
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre plată.
	// Cardinality: 0..1
	PaymentMeans *InvoicePaymentMeans `xml:"cac:PaymentMeans,omitempty"`
	// ID: BG-20
	// Term: DEDUCERI LA NIVELUL DOCUMENTULUI
	// Cardinality: 0..n
	// TODO:
	TaxTotal *InvoiceTaxTotal `xml:"cac:TaxTotal"`

	// ID: BG-22
	// Term: TOTALURILE DOCUMENTULUI
	// Cardinality: 1..1
	LegalMonetaryTotal InvoiceLegalMonetaryTotal `xml:"cac:LegalMonetaryTotal"`

	// ID: BG-25
	// Term: LINIE A FACTURII
	// Cardinality: 1..n
	InvoiceLines []InvoiceLine `xml:"cac:InvoiceLine"`
}

func (iv Invoice) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type invoice Invoice
	var xmliv struct {
		// These need to be first fields, because apparently the validators care
		// about the order of xml nodes.
		// Conditional / Identifies the earliest version of the UBL 2 schema for
		// this document type that defines all of the elements that might be
		// encountered in the current instance.
		// Path: /Invoice/cbc:UBLVersionID
		UBLVersionID string `xml:"cbc:UBLVersionID"`
		// ID: BT-24
		// Term: Identificatorul specificaţiei
		// Description: O identificare a specificaţiei care conţine totalitatea
		//     regulilor privind conţinutul semantic, cardinalităţile şi regulile
		//     operaţionale cu care datele conţinute în instanţa de factură sunt
		//     conforme.
		// Cardinality: 1..1
		CustomizationID string `xml:"cbc:CustomizationID"`

		invoice

		XMLName  xml.Name `xml:"Invoice"`
		XMLNS    string   `xml:"xmlns,attr"`
		XMLNScac string   `xml:"xmlns:cac,attr"`
		XMLNScbc string   `xml:"xmlns:cbc,attr"`
	}

	xmliv.XMLNS = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
	xmliv.XMLNScac = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	xmliv.XMLNScbc = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
	xmliv.UBLVersionID = "2.1"
	xmliv.CustomizationID = "urn:cen.eu:en16931:2017#compliant#urn:efactura.mfinante.ro:CIUS-RO:1.0.1"
	xmliv.invoice = invoice(iv)

	return e.EncodeElement(xmliv, start)
}

type InvoiceSupplier struct {
	// ID: BT-29
	// Term: Identificatorul Vânzătorului
	// Cardinality: 0..n
	SellerID []InvoiceSupplierIdentification `xml:"cac:PartyIdentification,omitempty"`
	// ID: BT-28
	// Term: Denumirea comercială a Vânzătorului
	// Description: Un nume sub care este cunoscut Vânzătorul, altul decât
	//     numele Vânzătorului (cunoscut, de asemenea, ca denumirea comercială).
	// Cardinality: 0..1
	CommercialName *InvoicePartyName `xml:"cac:PartyName,omitempty"`
	// ID: BG-5
	// Term: Adresa poștală a vânzătorului
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre adresa Vânzătorului.
	// Cardinality: 1..1
	PostalAddress InvoiceSupplierPostalAddress `xml:"cac:PostalAddress"`
	// ID: BT-31
	// Term: Identificatorul de TVA al Vânzătorului
	// Description: Identificatorul de TVA al Vânzătorului (cunoscut, de
	//     asemenea, ca numărul de identificare de TVA al Vânzătorului).
	// Cardinality: 0..1
	TaxScheme   *InvoicePartyTaxScheme     `xml:"cac:PartyTaxScheme,omitempty"`
	LegalEntity InvoiceSupplierLegalEntity `xml:"cac:PartyLegalEntity"`
	// ID: BT-32
	// Term: Identificatorul de înregistrare fiscală a Vânzătorului
	// Description: Identificarea locală (definită prin adresa Vânzătorului)
	//     a Vânzătorului pentru scopuri fiscale sau o referinţă care-i permite
	//     Vânzătorului să demonstreze că este înregistrat la administraţia
	//     fiscală.
	// Cardinality: 0..1
	// TODO:
	// ID: BG-6
	// Term: CONTACTUL VÂNZĂTORULUI
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     de contact despre Vânzător.
	// Cardinality: 0..1
	Contact *InvoiceSupplierContact `xml:"cac:Contact,omitempty"`
}

type InvoiceSupplierLegalEntity struct {
	// ID: BT-27
	// Term: Numele vânzătorului
	// Description: Denumirea oficială completă sub care Vânzătorul este
	//     înscris în registrul naţional al persoanelor juridice sau în calitate
	//     de Contribuabil sau îşi exercită activităţile în calitate de persoană
	//     sau grup de persoane.
	// Cardinality: 1..1
	Name string `xml:"cbc:RegistrationName"`
	// ID: BT-30
	// Term: Identificatorul de înregistrare legală a Vânzătorului
	// Description: Un identificator emis de un organism oficial de
	//     înregistrare care identifică Vânzătorul ca o entitate sau persoană
	//     juridică.
	// Cardinality: 1..1
	CompanyID *ValueWithScheme `xml:"cbc:CompanyID,omitempty"`
	// ID: BT-33
	// Term: Informaţii juridice suplimentare despre Vânzător
	// Cardinality: 0..1
	CompanyLegalForm string `xml:"cbc:CompanyLegalForm,omitempty"`
}

type InvoiceSupplierIdentification struct {
	ID ValueWithScheme `xml:"cbc:ID"`
}

type InvoiceSupplierIdentification_ struct {
	// ID: BT-29
	// Term: Identificatorul Vânzătorului
	// Cardinality: 0..n
	ID string `xml:"cbc:ID"`
	// ID: BT-29
	// Term: Identificatorul schemei
	// Cardinality: 0..1
	SchemeID string `xml:"cbc:ID>schemeID,attr,omitempty"`

	XMLName xml.Name `xml:"Invoice"`
}

type InvoiceSupplierPostalAddress struct {
	// ID: BT-35
	// Term: Adresa Vânzătorului - Linia 1
	// Description: Linia principală a unei adrese.
	// Cardinality: 0..1
	Line1 string `xml:"cbc:StreetName,omitempty"`
	// ID: BT-36
	// Term: Adresa Vânzătorului - Linia 2
	// Description: O linie suplimentară într-o adresă care poate fi utilizată
	//     pentru informaţii suplimentare şi completări la linia principală.
	// Cardinality: 0..1
	Line2 string `xml:"cbc:AdditionalStreetName,omitempty"`
	// ID: BT-162
	// Term: Adresa Vânzătorului - Linia 3
	// Description: O linie suplimentară într-o adresă care poate fi utilizată
	//     pentru informaţii suplimentare şi completări la linia principală.
	// Cardinality: 0..1
	Line3 string `xml:"cbc:AddressLine,omitempty"`
	// ID: BT-37
	// Term: Localitatea Vânzătorului
	// Description: Numele uzual al municipiului, oraşului sau satului, în care
	//     se află adresa Vânzătorului.
	// Cardinality: 0..1
	CityName string `xml:"cbc:CityName,omitempty"`
	// ID: BT-38
	// Term: Codul poştal al Vânzătorului
	// Cardinality: 0..1
	PostalZone string `xml:"cbc:PostalZone,omitempty"`
	// ID: BT-39
	// Term: Subdiviziunea ţării Vânzătorului
	// Cardinality: 0..1
	CountrySubentity string `xml:"cbc:CountrySubentity,omitempty"`
	// ID: BT-40
	// Term: Subdiviziunea ţării Vânzătorului
	// Cardinality: 1..1
	CountryIdentificationCode string `xml:"cac:Country>cbc:IdentificationCode"`
}

type InvoiceSupplierContact struct {
	// ID: BT-41
	// Term: Punctul de contact al Vânzătorului
	// Description: Un punct de contact pentru o entitate sau persoană
	//     juridică, cum ar fi numele persoanei, identificarea unui contact,
	//     departament sau serviciu.
	// Cardinality: 0..1
	Name string `xml:"cbc:Name,omitempty"`
	// ID: BT-42
	// Term: Numărul de telefon al contactului Vânzătorului
	// Description: Un număr de telefon pentru punctul de contact.
	// Cardinality: 0..1
	Phone string `xml:"cbc:Telephone,omitempty"`
	// ID: BT-43
	// Term: Adresa de email a contactului Vânzătorului
	// Description: O adresă de e-mail pentru punctul de contact.
	// Cardinality: 0..1
	Email string `xml:"cbc:ElectronicMail,omitempty"`
}

type InvoiceCustomer struct {
	// ID: BT-46
	// Term: Identificatorul Cumpărătorului
	// Cardinality: 0..n
	SellerID []InvoiceSupplierIdentification `xml:"cac:PartyIdentification,omitempty"`
	// ID: BT-45
	// Term: Denumirea comercială a Cumpărătorului
	// Description: Un nume sub care este cunoscut Cumpărătorul, altul decât
	//     numele Cumpărătorului (cunoscut, de asemenea, ca denumirea comercială).
	// Cardinality: 0..1
	CommercialName *InvoicePartyName `xml:"cac:PartyName,omitempty"`
	// ID: BG-8
	// Term: Adresa poștală a Cumpărătorului
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     despre adresa Cumpărătorului.
	// Cardinality: 1..1
	PostalAddress InvoiceCustomerPostalAddress `xml:"cac:PostalAddress"`
	// ID: BT-48
	// Term: Identificatorul de TVA al Cumpărătorului
	// Description: Identificatorul de TVA al Cumpărătorului (cunoscut, de
	//     asemenea, ca numărul de identificare de TVA al Cumpărătorului).
	// Cardinality: 0..1
	TaxScheme   *InvoicePartyTaxScheme     `xml:"cac:PartyTaxScheme"`
	LegalEntity InvoiceCustomerLegalEntity `xml:"cac:PartyLegalEntity"`
	// ID: BG-9
	// Term: Contactul Cumpărătorului
	// Description: Un grup de termeni operaţionali care furnizează informaţii
	//     de contact despre Cumpărător.
	// Cardinality: 0..1
	Contact *InvoiceCustomerContact `xml:"cac:Contact,omitempty"`
}

type InvoiceCustomerLegalEntity struct {
	// ID: BT-44
	// Term: Numele cumpărătorului
	// Description: Numele complet al Cumpărătorului.
	// Cardinality: 1..1
	Name string `xml:"cbc:RegistrationName"`
	// ID: BT-47
	// Term: Identificatorul de înregistrare legală a Cumpărătorului
	// Description: Un identificator emis de un organism oficial de
	//     înregistrare care identifică Cumpărătorul ca o entitate sau persoană
	//     juridică.
	// Cardinality: 1..1
	CompanyID *ValueWithScheme `xml:"cbc:CompanyID,omitempty"`
}

type InvoiceCustomerPostalAddress struct {
	// ID: BT-50
	// Term: Adresa Cumpărătorului - Linia 1
	// Description: Linia principală a unei adrese.
	// Cardinality: 0..1
	Line1 string `xml:"cbc:StreetName,omitempty"`
	// ID: BT-51
	// Term: Adresa Cumpărătorului - Linia 2
	// Description: O linie suplimentară într-o adresă care poate fi utilizată
	//     pentru informaţii suplimentare şi completări la linia principală.
	// Cardinality: 0..1
	Line2 string `xml:"cbc:AdditionalStreetName,omitempty"`
	// ID: BT-163
	// Term: Adresa Cumpărătorului - Linia 3
	// Description: O linie suplimentară într-o adresă care poate fi utilizată
	//     pentru informaţii suplimentare şi completări la linia principală.
	// Cardinality: 0..1
	Line3 string `xml:"cbc:AddressLine,omitempty"`
	// ID: BT-52
	// Term: Localitatea Cumpărătorului
	// Description: Numele uzual al municipiului, oraşului sau satului, în care
	//     se află adresa Cumpărătorului.
	// Cardinality: 0..1
	CityName string `xml:"cbc:CityName,omitempty"`
	// ID: BT-53
	// Term: Codul poştal al Cumpărătorului
	// Cardinality: 0..1
	PostalZone string `xml:"cbc:PostalZone,omitempty"`
	// ID: BT-54
	// Term: Subdiviziunea ţării Cumpărătorului
	// Cardinality: 0..1
	CountrySubentity string `xml:"cbc:CountrySubentity,omitempty"`
	// ID: BT-55
	// Term: Subdiviziunea ţării Cumpărătorului
	// Cardinality: 1..1
	CountryIdentificationCode string `xml:"cac:Country>cbc:IdentificationCode"`
}

type InvoiceCustomerContact struct {
	// ID: BT-56
	// Term: Punctul de contact al Cumpărătorului
	// Description: Un punct de contact pentru o entitate sau persoană
	//     juridică, cum ar fi numele persoanei, identificarea unui contact,
	//     departament sau serviciu.
	// Cardinality: 0..1
	Name string `xml:"cbc:Name,omitempty"`
	// ID: BT-57
	// Term: Numărul de telefon al contactului Cumpărătorului
	// Description: Un număr de telefon pentru punctul de contact.
	// Cardinality: 0..1
	Phone string `xml:"cbc:Telephone,omitempty"`
	// ID: BT-58
	// Term: Adresa de email a contactului Vânzătorului
	// Description: O adresă de e-mail pentru punctul de contact.
	// Cardinality: 0..1
	Email string `xml:"cbc:ElectronicMail,omitempty"`
}

type InvoicePeriod struct {
	// ID: BT-73
	// Term: Data de început a perioadei de facturare
	// Description: Data la care începe perioada de facturare.
	// Cardinality: 0..1
	StartDate *Date `xml:"cbc:StartDate,omitempty"`
	// ID: BT-74
	// Term: Data de sfârșit a perioadei de facturare
	// Description: Data la care sfârșește perioada de facturare.
	// Cardinality: 0..1
	EndDate *Date `xml:"cbc:EndDate,omitempty"`
}

type InvoicePaymentMeans struct {
	// ID: BT-81
	// Term: Codul tipului de instrument de plată
	// Description: Cod care indică modul în care o platătrebuie să fie sau a
	//     fost efectuată.
	// Cardinality: 1..1
	PaymentMeansCode PaymentMeansCode `xml:"cbc:PaymentMeansCode"`
	// ID: BT-83
	// Term: Aviz de plată
	// Description: Valoare textuală utilizată pentru a stabili o legătură
	//     între plată şi Factură, emisă de Vânzător.
	// Cardinality: 0..1
	PaymentID string `xml:"cbc:PaymentID,omitempty"`
	// ID: BG-17
	// Term: VIRAMENT
	// Cardinality: 0..n
	PayeeFinancialAccount []PayeeFinancialAccount `xml:"cac:PayeeFinancialAccount,omitempty"`
}

type PaymentMeansCode struct {
	Code string `xml:",chardata"`
	// ID: BT-82
	// Term: Explicaţii privind instrumentul de plată
	// Description: Text care indică modul în care o plată trebuie să fie sau
	//     a fost efectuată.
	// Cardinality: 0..1
	Name string `xml:"name,attr,omitempty"`
}

type PayeeFinancialAccount struct {
	// ID: BT-84
	// Term: Identificatorul contului de plată
	// Description: Un identificator unic al contului bancar de plată, la un
	//     furnizor de servicii de plată la care se recomandă să se facă plata
	// Cardinality: 1..1
	ID string `xml:"cbc:ID"`
	// ID: BT-85
	// Term: Numele contului de plată
	// Cardinality: 0..1
	Name string `xml:"bc:Name,omitempty"`
	// ID: BT-86
	// Term: Identificatorul furnizorului de servicii de plată.
	// Cardinality: 0..1
	// TODO:
}

type InvoiceTaxTotal struct {
	// ID: BT-110
	// Term: Valoarea totală a TVA a facturii
	// Cardinality: 0..1
	// ID: BT-111
	// Term: Valoarea totală a TVA a facturii în moneda de contabilizare
	// Cardinality: 0..1
	TaxAmount *AmountWithCurrency `xml:"cbc:TaxAmount,omitempty"`
	// ID: BG-23
	// Term: DETALIEREA TVA
	// Cardinality: 1..n
	TaxSubtotal []InvoiceTaxSubtotal `xml:"cac:TaxSubtotal"`
}

type InvoiceTaxSubtotal struct {
	// ID: BT-116
	// Term: Baza de calcul pentru categoria de TVA
	// Cardinality: 1..1
	TaxableAmount AmountWithCurrency `xml:"cbc:TaxableAmount"`
	// ID: BT-117
	// Term: Valoarea TVA pentru fiecare categorie de TVA
	// Cardinality: 1..1
	TaxAmount   AmountWithCurrency `xml:"cbc:TaxAmount"`
	TaxCategory InvoiceTaxCategory `xml:"cac:TaxCategory"`
}

type InvoiceTaxCategory struct {
	// ID: BT-118
	// Term: Codul categoriei de TVA
	// Cardinality: 1..1
	ID TaxCategoryCodeType `xml:"cbc:ID"`
	// ID: BT-119
	// Term: Cota categoriei de TVA
	// Cardinality: 0..1
	Percent *Decimal `xml:"cbc:Percent,omitempty"`
	// ID: BT-120
	// Term: Motivul scutirii de TVA
	// Cardinality: 0..1
	TaxExemptionReason string `xml:"cbc:TaxExemptionReason,omitempty"`
	// ID: BT-121
	// Term: Codul motivului scutirii de TVA
	// Cardinality: 0..1
	TaxExemptionReasonCode TaxExemptionReasonCodeType `xml:"cbc:TaxExemptionReasonCode,omitempty"`
	TaxSchemeID            ValueWithScheme            `xml:"cac:TaxScheme>cbc:ID,omitempty"`
}

type InvoiceLegalMonetaryTotal struct {
	// ID: BT-106
	// Term: Suma valorilor nete ale liniilor facturii
	// Cardinality: 1..1
	LineExtensionAmount AmountWithCurrency `xml:"cbc:LineExtensionAmount"`
	// ID: BT-107
	// Term: Suma deducerilor la nivelul documentului
	// Cardinality: 0..1
	AllowanceTotalAmount *AmountWithCurrency `xml:"cbc:AllowanceTotalAmount"`
	// ID: BT-108
	// Term: Suma taxelor suplimentare la nivelul documentului
	// Cardinality: 0..1
	ChargeTotalAmount *AmountWithCurrency `xml:"cbc:ChargeTotalAmount"`
	// ID: BT-109
	// Term: Valoarea totală a facturii fără TVA
	// Cardinality: 1..1
	TaxExclusiveAmount AmountWithCurrency `xml:"cbc:TaxExclusiveAmount"`
	// ID: BT-112
	// Term: Valoarea totală a facturii cu TVA
	// Cardinality: 1..1
	TaxInclusiveAmount AmountWithCurrency `xml:"cbc:TaxInclusiveAmount"`
	// ID: BT-115
	// Term: Suma de plată
	// Cardinality: 1..1
	PayableAmount AmountWithCurrency `xml:"cbc:PayableAmount"`
}

type InvoiceLine struct {
	// ID: BT-126
	// Term: Identificatorul liniei facturii
	// Cardinality: 1..1
	ID string `xml:"cbc:ID"`
	// ID: BT-127
	// Term: Nota liniei facturii
	// Description: O notă textuală care furnizează o informaţie nestructurată
	//     care este relevantă pentru linia facturii.
	// Cardinality: 0..1
	Note string `xml:"cbc:Note,omitempty"`
	// ID: BT-129
	// Term: Cantitatea facturată
	// Description: Cantitatea articolelor (bunuri sau servicii) luate în
	//     considerare în linia din factură.
	// Cardinality: 1..1
	// ID: BT-130
	// Term: Codul unităţii de măsură a cantităţii facturate
	// Cardinality: 1..1
	InvoicedQuantity InvoicedQuantity `xml:"cbc:InvoicedQuantity"`
	// ID: BT-131
	// Term: Valoarea netă a liniei facturii
	// Cardinality: 1..1
	LineExtensionAmount AmountWithCurrency `xml:"cbc:LineExtensionAmount"`
	// ID: BG-26
	// Term: Perioada de facturare a liniei
	// Cardinality: 0..1
	InvoicePeriod *InvoiceLinePeriod `xml:"cac:InvoicePeriod,omitempty"`
	// ID: BG-27 / BG-28
	// Term: DEDUCERI LA LINIA FACTURII / TAXE SUPLIMENTARE LA LINIA FACTURII
	// Cardinality: 0..n / 0..n
	AllowanceCharge []InvoiceAllowanceCharge `xml:"cac:AllowanceCharge,omitempty"`
	Item            InvoiceLineItem          `xml:"cac:Item"`
	// ID: BG-29
	// Term: DETALII ALE PREŢULUI
	// Cardinality: 1..1
	Price InvoiceLinePrice `xml:"cac:Price"`
}

type InvoicedQuantity struct {
	Quantity string `xml:",chardata"`
	UnitCode string `xml:"unitCode,attr,omitempty"`
}

type InvoiceLinePeriod struct {
	// ID: BT-134
	// Term: Data de început a perioadei de facturare a liniei facturii
	// Cardinality: 0..1
	StartDate *Date `xml:"cbc:StartDate,omitempty"`
	// ID: BT-135
	// Term: Data de sfârșit a perioadei de facturare
	// Cardinality: 0..1
	EndDate *Date `xml:"cbc:EndDate,omitempty"`
}

type InvoiceAllowanceCharge struct {
	// cbc:ChargeIndicator = false  ==>  BG-27 deducere
	// cbc:ChargeIndicator = true   ==>  BG-28 taxă suplimentară
	ChargeIndicator bool `xml:"cbc:ChargeIndicator"`
	// cbc:ChargeIndicator = false  ==>  {{{
	// ID: BT-140
	// Term: Codul motivului deducerii la linia facturii
	// Cardinality: 0..1
	// }}}
	// cbc:ChargeIndicator = true  ==>  {{{
	// ID: BT-145
	// Term: Codul motivului taxei suplimentare la linia facturii
	// Cardinality: 0..1
	// }}}
	AllowanceChargeReasonCode string `xml:"cbc:AllowanceChargeReasonCode,omitempty"`
	// cbc:ChargeIndicator = false  ==>  {{{
	// ID: BT-139
	// Term: Motivul deducerii la linia facturii
	// Cardinality: 0..1
	// }}}
	// cbc:ChargeIndicator = true  ==>  {{{
	// ID: BT-144
	// Term: Motivul taxei suplimentare la linia facturii
	// Cardinality: 0..1
	// }}}
	AllowanceChargeReason string `xml:"cbc:AllowanceChargeReason,omitempty"`
	// cbc:ChargeIndicator = false  ==>  {{{
	// ID: BT-136
	// Term: Valoarea deducerii la linia facturii
	// Description: fără TVA
	// Cardinality: 1..1
	// }}}
	// cbc:ChargeIndicator = true  ==>  {{{
	// ID: BT-141
	// Term: Valoarea taxei suplimentare la linia facturii
	// Description: fără TVA
	// Cardinality: 1..1
	// }}}
	Amount AmountWithCurrency `xml:"cbc:Amount"`
	// cbc:ChargeIndicator = false  ==>  {{{
	// ID: BT-137
	// Term: Valoarea de bază a deducerii la linia facturii
	// Description: Valoarea de bază care poate fi utilizată, împreună cu
	//     procentajul deducerii la linia facturii, pentru a calcula valoarea
	//     deducerii la linia facturii.
	// Cardinality: 0..1
	// }}}
	// cbc:ChargeIndicator = true  ==>  {{{
	// ID: BT-137
	// Term: Valoarea de bază a taxei suplimentare la linia facturii
	// Description: Valoarea de bază care poate fi utilizată, împreună cu
	//     procentajul taxei suplimentare la linia facturii, pentru a calcula
	//     valoarea taxei suplimentare la linia facturii.
	// Cardinality: 0..1
	// }}}
	BaseAmount *AmountWithCurrency `xml:"cbc:BaseAmount"`
}

type InvoiceLinePrice struct {
	// ID: BT-146
	// Term: Preţul net al articolului
	// Description: Preţul unui articol, exclusiv TVA, după aplicarea reducerii
	//     la preţul articolului.
	// Cardinality: 1..1
	PriceAmount AmountWithCurrency `xml:"cbc:PriceAmount"`
	// ID: BT-149
	// Term: Cantitatea de bază a preţului articolului
	// Cardinality: 0..1
	// ID: BT-150
	// Term: Codul unităţii de măsură a cantităţii de bază a preţului articolului
	// Cardinality: 0..1
	BaseQuantity *InvoicedQuantity `xml:"cbc:BaseQuantity,omitempty"`
}

type InvoiceLineItem struct {
	// ID: BG-31
	// Term: INFORMAȚII PRIVIND ARTICOLUL

	// ID: BT-153
	// Term: Numele articolului
	// Cardinality: 1..1
	Name string `xml:"cbc:Name"`
	// ID: BT-154
	// Term: Descrierea articolului
	// Cardinality: 0..1
	Description string `xml:"cbc:Description,omitempty"`
	// ID: BT-155
	// Term: BT-155
	// Cardinality: 0..1
	SellerItemID string `xml:"cac:SellersItemIdentification>cbc:ID"`
	// ID: BT-157/BT-157-1
	// Term: Identificatorul standard al articolului / Identificatorul schemei
	StandardItemIdentification *ItemStandardIdentificationCode `xml:"cac:StandardItemIdentification,omitempty"`
	// ID: BT-158/BT-158-1
	// Term: Identificatorul clasificării articolului / Identificatorul schemei
	ItemClassificationCode *ItemClassificationCode `xml:"cac:CommodityClassification>cbc:ItemClassificationCode,omitempty"`
	// ID: BG-30
	// Term: INFORMAŢII PRIVIND TVA A LINIEI
	ClassifiedTaxCategory *InvoiceClassifiedTaxCategory `xml:"cac:ClassifiedTaxCategory,omitempty"`
}

type ItemStandardIdentificationCode struct {
	Code     string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr"`
}

type ItemClassificationCode struct {
	Code   string `xml:",chardata"`
	ListID string `xml:"listID,attr,omitempty"`
}

type InvoicePartyName struct {
	Name string `xml:"cbc:Name"`
}

type InvoicePartyTaxScheme struct {
	CompanyID   string          `xml:"cbc:CompanyID,omitempty"`
	TaxSchemeID ValueWithScheme `xml:"cac:TaxScheme>cbc:ID,omitempty"`
}

// TODO: document
type InvoiceClassifiedTaxCategory struct {
	ID          TaxCategoryCodeType `xml:"cbc:ID"`
	Percent     *Decimal            `xml:"cbc:Percent,omitempty"`
	TaxSchemeID ValueWithScheme     `xml:"cac:TaxScheme>cbc:ID,omitempty"`
}
