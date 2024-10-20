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
	"strings"

	"github.com/printesoi/e-factura-go/pkg/text"
	"github.com/printesoi/e-factura-go/pkg/units"
)

type DeclPostIncidentType string

const (
	// "D" - Singura valoare posibilă pentru atributul "declPostAvarie"
	// (declarare până la sfârșitul următoarei zile lucrătoare repunerii în
	// funcțiune a sistemului, cf. art. 8, alin. 1^3 al OUG41/2022 aşa cum este
	// modificată prin OUG132/2022), dacă este cazul. Pentru declarare normală
	// - anterioară punerii în mişcare a vehiculului pe parcurs rutier naţional
	// - atributul nu trebuie specificat.
	DeclPostIncident DeclPostIncidentType = "D"
)

type OpType string

const (
	// "10" AIC - Achiziţie intracomunitară
	OpTypeAIC OpType = "10"
	// "12" LHI - Operațiuni în sistem lohn (UE) - intrare
	OpTypeLHI OpType = "12"
	// "14" SCI - Stocuri la dispoziția clientului (Call-off stock) - intrare
	OpTypeSCI OpType = "14"
	// "20" LIC - Livrare intracomunitară
	OpTypeLIC OpType = "20"
	// "22" LHE - Operațiuni în sistem lohn (UE) - ieșire
	OpTypeLHE OpType = "22"
	// "24" SCE - Stocuri la dispoziția clientului (Call-off stock) - ieșire
	OpTypeSCE OpType = "24"
	// "30" TTN - Transport pe teritoriul naţional
	OpTypeTTN OpType = "30"
	// "40" IMP - Import
	OpTypeIMP OpType = "40"
	// "50" EXP - Export
	OpTypeEXP OpType = "50"
	// "60" DIN - Tranzacţie intracomunitară - Intrare pentru depozitare/formare nou transport
	OpTypeDIN OpType = "60"
	// "70" DIE - Tranzacţie intracomunitară - Ieşire după depozitare/formare nou transport
	OpTypeDIE OpType = "70"
)

type CountryCodeType string

const (
	// RO - Romania
	CountryCodeRO CountryCodeType = "RO"
	// AD - Andorra
	CountryCodeAD CountryCodeType = "AD"
	// AE - United Arab Emirates
	CountryCodeAE CountryCodeType = "AE"
	// AF - Afghanistan
	CountryCodeAF CountryCodeType = "AF"
	// AG - Antigua and Barbuda
	CountryCodeAG CountryCodeType = "AG"
	// AI - Anguilla
	CountryCodeAI CountryCodeType = "AI"
	// AL - Albania
	CountryCodeAL CountryCodeType = "AL"
	// AM - Armenia
	CountryCodeAM CountryCodeType = "AM"
	// AO - Angola
	CountryCodeAO CountryCodeType = "AO"
	// AQ - Antarctica
	CountryCodeAQ CountryCodeType = "AQ"
	// AR - Argentina
	CountryCodeAR CountryCodeType = "AR"
	// AS - American Samoa
	CountryCodeAS CountryCodeType = "AS"
	// AT - Austria
	CountryCodeAT CountryCodeType = "AT"
	// AU - Australia
	CountryCodeAU CountryCodeType = "AU"
	// AW - Aruba
	CountryCodeAW CountryCodeType = "AW"
	// AX - Aland Islands
	CountryCodeAX CountryCodeType = "AX"
	// AZ - Azerbaijan
	CountryCodeAZ CountryCodeType = "AZ"
	// BA - Bosnia and Herzegovina
	CountryCodeBA CountryCodeType = "BA"
	// BB - Barbados
	CountryCodeBB CountryCodeType = "BB"
	// BD - Bangladesh
	CountryCodeBD CountryCodeType = "BD"
	// BE - Belgium
	CountryCodeBE CountryCodeType = "BE"
	// BF - Burkina Faso
	CountryCodeBF CountryCodeType = "BF"
	// BG - Bulgaria
	CountryCodeBG CountryCodeType = "BG"
	// BH - Bahrain
	CountryCodeBH CountryCodeType = "BH"
	// BI - Burundi
	CountryCodeBI CountryCodeType = "BI"
	// BJ - Benin
	CountryCodeBJ CountryCodeType = "BJ"
	// BL - Saint Barth&#xE9;lemy
	CountryCodeBL CountryCodeType = "BL"
	// BM - Bermuda
	CountryCodeBM CountryCodeType = "BM"
	// BN - Brunei Darussalam
	CountryCodeBN CountryCodeType = "BN"
	// BO - Bolivia, Plurinational State of
	CountryCodeBO CountryCodeType = "BO"
	// BQ - Bonaire, Sint Eustatius and Saba
	CountryCodeBQ CountryCodeType = "BQ"
	// BR - Brazil
	CountryCodeBR CountryCodeType = "BR"
	// BS - Bahamas
	CountryCodeBS CountryCodeType = "BS"
	// BT - Bhutan
	CountryCodeBT CountryCodeType = "BT"
	// BV - Bouvet Island
	CountryCodeBV CountryCodeType = "BV"
	// BW - Botswana
	CountryCodeBW CountryCodeType = "BW"
	// BY - Belarus
	CountryCodeBY CountryCodeType = "BY"
	// BZ - Belize
	CountryCodeBZ CountryCodeType = "BZ"
	// CA - Canada
	CountryCodeCA CountryCodeType = "CA"
	// CC - Cocos (Keeling) Islands
	CountryCodeCC CountryCodeType = "CC"
	// CD - Congo, the Democratic Republic of the
	CountryCodeCD CountryCodeType = "CD"
	// CF - Central African Republic
	CountryCodeCF CountryCodeType = "CF"
	// CG - Congo
	CountryCodeCG CountryCodeType = "CG"
	// CH - Switzerland
	CountryCodeCH CountryCodeType = "CH"
	// CI - C&#xF4;te d'Ivoire
	CountryCodeCI CountryCodeType = "CI"
	// CK - Cook Islands
	CountryCodeCK CountryCodeType = "CK"
	// CL - Chile
	CountryCodeCL CountryCodeType = "CL"
	// CM - Cameroon
	CountryCodeCM CountryCodeType = "CM"
	// CN - China
	CountryCodeCN CountryCodeType = "CN"
	// CO - Colombia
	CountryCodeCO CountryCodeType = "CO"
	// CR - Costa Rica
	CountryCodeCR CountryCodeType = "CR"
	// CU - Cuba
	CountryCodeCU CountryCodeType = "CU"
	// CV - Cabo Verde
	CountryCodeCV CountryCodeType = "CV"
	// CW - Cura&#xE7;ao
	CountryCodeCW CountryCodeType = "CW"
	// CX - Christmas Island
	CountryCodeCX CountryCodeType = "CX"
	// CY - Cyprus
	CountryCodeCY CountryCodeType = "CY"
	// CZ - Czechia
	CountryCodeCZ CountryCodeType = "CZ"
	// DE - Germany
	CountryCodeDE CountryCodeType = "DE"
	// DJ - Djibouti
	CountryCodeDJ CountryCodeType = "DJ"
	// DK - Denmark
	CountryCodeDK CountryCodeType = "DK"
	// DM - Dominica
	CountryCodeDM CountryCodeType = "DM"
	// DO - Dominican Republic
	CountryCodeDO CountryCodeType = "DO"
	// DZ - Algeria
	CountryCodeDZ CountryCodeType = "DZ"
	// EC - Ecuador
	CountryCodeEC CountryCodeType = "EC"
	// EE - Estonia
	CountryCodeEE CountryCodeType = "EE"
	// EG - Egypt
	CountryCodeEG CountryCodeType = "EG"
	// EH - Western Sahara
	CountryCodeEH CountryCodeType = "EH"
	// EL - Greece. For ease of use both CountryCodeGR and CountryCodeEL
	// represent Greece.
	CountryCodeEL CountryCodeType = "EL"
	// ER - Eritrea
	CountryCodeER CountryCodeType = "ER"
	// ES - Spain
	CountryCodeES CountryCodeType = "ES"
	// ET - Ethiopia
	CountryCodeET CountryCodeType = "ET"
	// FI - Finland
	CountryCodeFI CountryCodeType = "FI"
	// FJ - Fiji
	CountryCodeFJ CountryCodeType = "FJ"
	// FK - Falkland Islands (Malvinas)
	CountryCodeFK CountryCodeType = "FK"
	// FM - Micronesia, Federated States of
	CountryCodeFM CountryCodeType = "FM"
	// FO - Faroe Islands
	CountryCodeFO CountryCodeType = "FO"
	// FR - France
	CountryCodeFR CountryCodeType = "FR"
	// GA - Gabon
	CountryCodeGA CountryCodeType = "GA"
	// GB - United Kingdom of Great Britain and Northern Ireland
	CountryCodeGB CountryCodeType = "GB"
	// GD - Grenada
	CountryCodeGD CountryCodeType = "GD"
	// GE - Georgia
	CountryCodeGE CountryCodeType = "GE"
	// GF - French Guiana
	CountryCodeGF CountryCodeType = "GF"
	// GG - Guernsey
	CountryCodeGG CountryCodeType = "GG"
	// GH - Ghana
	CountryCodeGH CountryCodeType = "GH"
	// GI - Gibraltar
	CountryCodeGI CountryCodeType = "GI"
	// GL - Greenland
	CountryCodeGL CountryCodeType = "GL"
	// GM - Gambia
	CountryCodeGM CountryCodeType = "GM"
	// GN - Guinea
	CountryCodeGN CountryCodeType = "GN"
	// GP - Guadeloupe
	CountryCodeGP CountryCodeType = "GP"
	// EL - Greece. For ease of use both CountryCodeGR and CountryCodeEL
	// represent Greece.
	CountryCodeGR CountryCodeType = CountryCodeEL
	// GQ - Equatorial Guinea
	CountryCodeGQ CountryCodeType = "GQ"
	// GS - South Georgia and the South Sandwich Islands
	CountryCodeGS CountryCodeType = "GS"
	// GT - Guatemala
	CountryCodeGT CountryCodeType = "GT"
	// GU - Guam
	CountryCodeGU CountryCodeType = "GU"
	// GW - Gu inea-Bissau
	CountryCodeGW CountryCodeType = "GW"
	// GY - Guyana
	CountryCodeGY CountryCodeType = "GY"
	// HK - Hong Kong
	CountryCodeHK CountryCodeType = "HK"
	// HM - Heard Island and McDonald Islands
	CountryCodeHM CountryCodeType = "HM"
	// HN - Honduras
	CountryCodeHN CountryCodeType = "HN"
	// HR - Croatia
	CountryCodeHR CountryCodeType = "HR"
	// HT - Haiti
	CountryCodeHT CountryCodeType = "HT"
	// HU - Hungary
	CountryCodeHU CountryCodeType = "HU"
	// ID - Indonesia
	CountryCodeID CountryCodeType = "ID"
	// IE - Ireland
	CountryCodeIE CountryCodeType = "IE"
	// IL - Israel
	CountryCodeIL CountryCodeType = "IL"
	// IM - Isle of Man
	CountryCodeIM CountryCodeType = "IM"
	// IN - India
	CountryCodeIN CountryCodeType = "IN"
	// IO - British Indian Ocean Territory
	CountryCodeIO CountryCodeType = "IO"
	// IQ - Iraq
	CountryCodeIQ CountryCodeType = "IQ"
	// IR - Iran, Islamic Republic of
	CountryCodeIR CountryCodeType = "IR"
	// IS - Iceland
	CountryCodeIS CountryCodeType = "IS"
	// IT - Italy
	CountryCodeIT CountryCodeType = "IT"
	// JE - Jersey
	CountryCodeJE CountryCodeType = "JE"
	// JM - Jamaica
	CountryCodeJM CountryCodeType = "JM"
	// JO - Jordan
	CountryCodeJO CountryCodeType = "JO"
	// JP - Japan
	CountryCodeJP CountryCodeType = "JP"
	// KE - Kenya
	CountryCodeKE CountryCodeType = "KE"
	// KG - Kyrgyzstan
	CountryCodeKG CountryCodeType = "KG"
	// KH - Cambodia
	CountryCodeKH CountryCodeType = "KH"
	// KI - Kiribati
	CountryCodeKI CountryCodeType = "KI"
	// KM - Comoros
	CountryCodeKM CountryCodeType = "KM"
	// KN - Saint Kitts and Nevis
	CountryCodeKN CountryCodeType = "KN"
	// KP - Korea, Democratic People's Republic of
	CountryCodeKP CountryCodeType = "KP"
	// KR - Korea, Republic of
	CountryCodeKR CountryCodeType = "KR"
	// KW - Kuwait
	CountryCodeKW CountryCodeType = "KW"
	// KY - Cayman Islands
	CountryCodeKY CountryCodeType = "KY"
	// KZ - Kazakhstan
	CountryCodeKZ CountryCodeType = "KZ"
	// LA - Lao People's Democratic Republic
	CountryCodeLA CountryCodeType = "LA"
	// LB - Lebanon
	CountryCodeLB CountryCodeType = "LB"
	// LC - Saint Lucia
	CountryCodeLC CountryCodeType = "LC"
	// LI - Liechtenstein
	CountryCodeLI CountryCodeType = "LI"
	// LK - Sri Lanka
	CountryCodeLK CountryCodeType = "LK"
	// LR - Liberia
	CountryCodeLR CountryCodeType = "LR"
	// LS - Lesotho
	CountryCodeLS CountryCodeType = "LS"
	// LT - Lithuania
	CountryCodeLT CountryCodeType = "LT"
	// LU - Luxembourg
	CountryCodeLU CountryCodeType = "LU"
	// LV - Latvia
	CountryCodeLV CountryCodeType = "LV"
	// LY - Libya
	CountryCodeLY CountryCodeType = "LY"
	// MA - Morocco
	CountryCodeMA CountryCodeType = "MA"
	// MC - Monaco
	CountryCodeMC CountryCodeType = "MC"
	// MD - Moldova, Republic of
	CountryCodeMD CountryCodeType = "MD"
	// ME - Montenegro
	CountryCodeME CountryCodeType = "ME"
	// MF - Saint Martin (French part)
	CountryCodeMF CountryCodeType = "MF"
	// MG - Madagascar
	CountryCodeMG CountryCodeType = "MG"
	// MH - Marshall Islands
	CountryCodeMH CountryCodeType = "MH"
	// MK - Macedonia, the former Yugoslav Republic of
	CountryCodeMK CountryCodeType = "MK"
	// ML - Mali
	CountryCodeML CountryCodeType = "ML"
	// MM - Myanmar
	CountryCodeMM CountryCodeType = "MM"
	// MN - Mongolia
	CountryCodeMN CountryCodeType = "MN"
	// MO - Macao
	CountryCodeMO CountryCodeType = "MO"
	// MP - Northern Mariana Islands
	CountryCodeMP CountryCodeType = "MP"
	// MQ - Martinique
	CountryCodeMQ CountryCodeType = "MQ"
	// MR - Mauritania
	CountryCodeMR CountryCodeType = "MR"
	// MS - Montserrat
	CountryCodeMS CountryCodeType = "MS"
	// MT - Malta
	CountryCodeMT CountryCodeType = "MT"
	// MU - Mauritius
	CountryCodeMU CountryCodeType = "MU"
	// MV - Maldives
	CountryCodeMV CountryCodeType = "MV"
	// MW - Malawi
	CountryCodeMW CountryCodeType = "MW"
	// MX - Mexico
	CountryCodeMX CountryCodeType = "MX"
	// MY - Malaysia
	CountryCodeMY CountryCodeType = "MY"
	// MZ - Mozambique
	CountryCodeMZ CountryCodeType = "MZ"
	// NA - Namibia
	CountryCodeNA CountryCodeType = "NA"
	// NC - New Caledonia
	CountryCodeNC CountryCodeType = "NC"
	// NE - Niger
	CountryCodeNE CountryCodeType = "NE"
	// NF - Norfolk Island
	CountryCodeNF CountryCodeType = "NF"
	// NG - Nigeria
	CountryCodeNG CountryCodeType = "NG"
	// NI - Nicaragua
	CountryCodeNI CountryCodeType = "NI"
	// NL - Netherlands
	CountryCodeNL CountryCodeType = "NL"
	// NO - Norway
	CountryCodeNO CountryCodeType = "NO"
	// NP - Nepal
	CountryCodeNP CountryCodeType = "NP"
	// NR - Nauru
	CountryCodeNR CountryCodeType = "NR"
	// NU - Niue
	CountryCodeNU CountryCodeType = "NU"
	// NZ - New Zealand
	CountryCodeNZ CountryCodeType = "NZ"
	// OM - Oman
	CountryCodeOM CountryCodeType = "OM"
	// PA - Panama
	CountryCodePA CountryCodeType = "PA"
	// PE - Peru
	CountryCodePE CountryCodeType = "PE"
	// PF - French Polynesia
	CountryCodePF CountryCodeType = "PF"
	// PG - Papua New Guinea
	CountryCodePG CountryCodeType = "PG"
	// PH - Philippines
	CountryCodePH CountryCodeType = "PH"
	// PK - Pakistan
	CountryCodePK CountryCodeType = "PK"
	// PL - Poland
	CountryCodePL CountryCodeType = "PL"
	// PM - Saint Pierre and Miquelon
	CountryCodePM CountryCodeType = "PM"
	// PN - Pitcairn
	CountryCodePN CountryCodeType = "PN"
	// PR - Puerto Rico
	CountryCodePR CountryCodeType = "PR"
	// PS - Palestine, State of
	CountryCodePS CountryCodeType = "PS"
	// PT - Portugal
	CountryCodePT CountryCodeType = "PT"
	// PW - Palau
	CountryCodePW CountryCodeType = "PW"
	// PY - Paraguay
	CountryCodePY CountryCodeType = "PY"
	// QA - Qatar
	CountryCodeQA CountryCodeType = "QA"
	// RE - R&#xE9;union
	CountryCodeRE CountryCodeType = "RE"
	// RS - Serbia
	CountryCodeRS CountryCodeType = "RS"
	// RU - Russian Federation
	CountryCodeRU CountryCodeType = "RU"
	// RW - Rwanda
	CountryCodeRW CountryCodeType = "RW"
	// SA - Saudi Arabia
	CountryCodeSA CountryCodeType = "SA"
	// SB - Solomon Islands
	CountryCodeSB CountryCodeType = "SB"
	// SC - Seychelles
	CountryCodeSC CountryCodeType = "SC"
	// SD - Sudan
	CountryCodeSD CountryCodeType = "SD"
	// SE - Sweden
	CountryCodeSE CountryCodeType = "SE"
	// SG - Singapore
	CountryCodeSG CountryCodeType = "SG"
	// SH - Saint Helena, Ascension and Tristan da Cunha
	CountryCodeSH CountryCodeType = "SH"
	// SI - Slovenia
	CountryCodeSI CountryCodeType = "SI"
	// SJ - Svalbard and Jan Mayen
	CountryCodeSJ CountryCodeType = "SJ"
	// SK - Slovakia
	CountryCodeSK CountryCodeType = "SK"
	// SL - Sierra Leone
	CountryCodeSL CountryCodeType = "SL"
	// SM - San Marino
	CountryCodeSM CountryCodeType = "SM"
	// SN - Senegal
	CountryCodeSN CountryCodeType = "SN"
	// SO - Somalia
	CountryCodeSO CountryCodeType = "SO"
	// SR - Suriname
	CountryCodeSR CountryCodeType = "SR"
	// SS - South Sudan
	CountryCodeSS CountryCodeType = "SS"
	// ST - Sao Tome and Principe
	CountryCodeST CountryCodeType = "ST"
	// SV - El Salvador
	CountryCodeSV CountryCodeType = "SV"
	// SX - Sint Maarten (Dutch part)
	CountryCodeSX CountryCodeType = "SX"
	// SY - Syrian Arab Republic
	CountryCodeSY CountryCodeType = "SY"
	// SZ - Swaziland
	CountryCodeSZ CountryCodeType = "SZ"
	// TC - Turks and Caicos Islands
	CountryCodeTC CountryCodeType = "TC"
	// TD - Chad
	CountryCodeTD CountryCodeType = "TD"
	// TF - French Southern Territories
	CountryCodeTF CountryCodeType = "TF"
	// TG - Togo
	CountryCodeTG CountryCodeType = "TG"
	// TH - Thailand
	CountryCodeTH CountryCodeType = "TH"
	// TJ - Tajikistan
	CountryCodeTJ CountryCodeType = "TJ"
	// TK - Tokelau
	CountryCodeTK CountryCodeType = "TK"
	// TL - Timor-Leste
	CountryCodeTL CountryCodeType = "TL"
	// TM - Turkmenistan
	CountryCodeTM CountryCodeType = "TM"
	// TN - Tunisia
	CountryCodeTN CountryCodeType = "TN"
	// TO - Tonga
	CountryCodeTO CountryCodeType = "TO"
	// TR - Turkey
	CountryCodeTR CountryCodeType = "TR"
	// TT - Trinidad and Tobago
	CountryCodeTT CountryCodeType = "TT"
	// TV - Tuvalu
	CountryCodeTV CountryCodeType = "TV"
	// TW - Taiwan, Province of China
	CountryCodeTW CountryCodeType = "TW"
	// TZ - Tanzania, United Republic of
	CountryCodeTZ CountryCodeType = "TZ"
	// UA - Ukraine
	CountryCodeUA CountryCodeType = "UA"
	// UG - Uganda
	CountryCodeUG CountryCodeType = "UG"
	// UM - United States Minor Outlying Islands
	CountryCodeUM CountryCodeType = "UM"
	// US - United States of America
	CountryCodeUS CountryCodeType = "US"
	// UY - Uruguay
	CountryCodeUY CountryCodeType = "UY"
	// UZ - Uzbekistan
	CountryCodeUZ CountryCodeType = "UZ"
	// VA - Holy See
	CountryCodeVA CountryCodeType = "VA"
	// VC - Saint Vincent and the Grenadines
	CountryCodeVC CountryCodeType = "VC"
	// VE - Venezuela, Bolivarian Republic of
	CountryCodeVE CountryCodeType = "VE"
	// VG - Virgin Islands, British
	CountryCodeVG CountryCodeType = "VG"
	// VI - Virgin Islands, U.S.
	CountryCodeVI CountryCodeType = "VI"
	// VN - Viet Nam
	CountryCodeVN CountryCodeType = "VN"
	// VU - Vanuatu
	CountryCodeVU CountryCodeType = "VU"
	// WF - Wallis and Futuna
	CountryCodeWF CountryCodeType = "WF"
	// WS - Samoa
	CountryCodeWS CountryCodeType = "WS"
	// YE - Yemen
	CountryCodeYE CountryCodeType = "YE"
	// YT - Mayotte
	CountryCodeYT CountryCodeType = "YT"
	// XK - Kosovo
	CountryCodeXK CountryCodeType = "XK"
	// ZA - South Africa
	CountryCodeZA CountryCodeType = "ZA"
	// ZM - Zambia
	CountryCodeZM CountryCodeType = "ZM"
	// ZW - Zimbabwe
	CountryCodeZW CountryCodeType = "ZW"
	// 1A - Kosovo
	CountryCode1A CountryCodeType = "1A"
)

// BCPCodeType represent the code of a border crossing point (BCP)
type BCPCodeType string

const (
	// "1" Petea (HU)
	BCPPetea BCPCodeType = "1"
	// "2" Borș (HU)
	BCPBors BCPCodeType = "2"
	// "3" Vărșand (HU)
	BCPVarsand BCPCodeType = "3"
	// "4" Nădlac (HU)
	BCPNadlac BCPCodeType = "4"
	// "5" Calafat (BG)
	BCPCalafat BCPCodeType = "5"
	// "6" Bechet (BG)
	BCPBechet BCPCodeType = "6"
	// "7" Turnu Măgurele (BG)
	BCPTurnuMagurele BCPCodeType = "7"
	// "8" Zimnicea (BG)
	BCPZimnicea BCPCodeType = "8"
	// "9" Giurgiu (BG)
	BCPGiurgiu BCPCodeType = "9"
	// "10" Ostrov (BG)
	BCPOstrov BCPCodeType = "10"
	// "11" Negru Vodă (BG)
	BCPNegruVoda BCPCodeType = "11"
	// "12" Vama Veche (BG)
	BCPVamaVeche BCPCodeType = "12"
	// "13" Călărași (BG)
	BCPCalarasi BCPCodeType = "13"
	// "14" Corabia (BG)
	BCPCorabia BCPCodeType = "14"
	// "15" Oltenița (BG)
	BCPOltenita BCPCodeType = "15"
	// "16" Carei (HU)
	BCPCarei BCPCodeType = "16"
	// "17" Cenad (HU)
	BCPCenad BCPCodeType = "17"
	// "18" Episcopia Bihor (HU)
	BCPEpiscopiaBihor BCPCodeType = "18"
	// "19" Salonta (HU)
	BCPSalonta BCPCodeType = "19"
	// "20" Săcuieni (HU)
	BCPSacuieni BCPCodeType = "20"
	// "21" Turnu (HU)
	BCPTurnu BCPCodeType = "21"
	// "22" Urziceni (HU)
	BCPUrziceni BCPCodeType = "22"
	// "23" Valea lui Mihai (HU)
	BCPValeaLuiMihai BCPCodeType = "23"
	// "24" Vladimirescu (HU)
	BCPVladimirescu BCPCodeType = "24"
	// "25" Porțile de Fier 1 (RS)
	BCPPortileDeFier1 BCPCodeType = "25"
	// "26" Naidăș (RS)
	BCPNaidas BCPCodeType = "26"
	// "27" Stamora Moravița (RS)
	BCPStamoraMoravita BCPCodeType = "27"
	// "28" Jimbolia (RS)
	BCPJimbolia BCPCodeType = "28"
	// "29" Halmeu (UA)
	BCPHalmeu BCPCodeType = "29"
	// "30" Stânca Costești (MD)
	BCPStancaCostești BCPCodeType = "30"
	// "31" Sculeni (MD)
	BCPSculeni BCPCodeType = "31"
	// "32" Albița (MD)
	BCPAlbita BCPCodeType = "32"
	// "33" Oancea (MD)
	BCPOancea BCPCodeType = "33"
	// "34" Galați Giurgiulești (MD)
	BCPGalatiGiurgiulești BCPCodeType = "34"
	// "35" Constanța Sud Agigea (-)
	BCPConstantaSudAgigea BCPCodeType = "35"
	// "36" Siret (UA)
	BCPSiret BCPCodeType = "36"
	// "37" Nădlac 2 - A1 (HU)
	BCPNadlac2 BCPCodeType = "37"
	// "38" Borș 2 - A3 (HU)
	BCPBors2 BCPCodeType = "38"
)

// CustomsOfficeCodeType represents the customs office code type
// Possible values (COI/X - Customs Office Inland/Exit):
// Valori posibile (BVI/F - Birou Vamal de Interior/Frontieră):
type CustomsOfficeCodeType string

const (
	// "12801" BVI Alba Iulia (ROBV0300)
	CustomsOfficeROBV0300 CustomsOfficeCodeType = "12801"
	// "22801" BVI Arad (ROTM0200)
	CustomsOfficeROTM0200 CustomsOfficeCodeType = "22801"
	// "22901" BVF Arad Aeroport (ROTM0230)
	CustomsOfficeROTM0230 CustomsOfficeCodeType = "22901"
	// "22902" BVF Zona Liberă Curtici (ROTM2300)
	CustomsOfficeROTM2300 CustomsOfficeCodeType = "22902"
	// "32801" BVI Pitești (ROCR7000)
	CustomsOfficeROCR7000 CustomsOfficeCodeType = "32801"
	// "42801" BVI Bacău (ROIS0600)
	CustomsOfficeROIS0600 CustomsOfficeCodeType = "42801"
	// "42901" BVF Bacău Aeroport (ROIS0620)
	CustomsOfficeROIS0620 CustomsOfficeCodeType = "42901"
	// "52801" BVI Oradea (ROCJ6570)
	CustomsOfficeROCJ6570 CustomsOfficeCodeType = "52801"
	// "52901" BVF Oradea Aeroport (ROCJ6580)
	CustomsOfficeROCJ6580 CustomsOfficeCodeType = "52901"
	// "62801" BVI Bistriţa-Năsăud (ROCJ0400)
	CustomsOfficeROCJ0400 CustomsOfficeCodeType = "62801"
	// "72801" BVI Botoşani (ROIS1600)
	CustomsOfficeROIS1600 CustomsOfficeCodeType = "72801"
	// "72901" BVF Stanca Costeşti (ROIS1610)
	CustomsOfficeROIS1610 CustomsOfficeCodeType = "72901"
	// "72902" BVF Rădăuţi Prut (ROIS1620)
	CustomsOfficeROIS1620 CustomsOfficeCodeType = "72902"
	// "82801" BVI Braşov (ROBV0900)
	CustomsOfficeROBV0900 CustomsOfficeCodeType = "82801"
	// "92901" BVF Zona Liberă Brăila (ROGL0710)
	CustomsOfficeROGL0710 CustomsOfficeCodeType = "92901"
	// "92902" BVF Brăila (ROGL0700)
	CustomsOfficeROGL0700 CustomsOfficeCodeType = "92902"
	// "102801" BVI Buzău (ROGL1500)
	CustomsOfficeROGL1500 CustomsOfficeCodeType = "102801"
	// "112801" BVI Reșița (ROTM7600)
	CustomsOfficeROTM7600 CustomsOfficeCodeType = "112801"
	// "112901" BVF Naidăș (ROTM6100)
	CustomsOfficeROTM6100 CustomsOfficeCodeType = "112901"
	// "122801" BVI Cluj Napoca (ROCJ1800)
	CustomsOfficeROCJ1800 CustomsOfficeCodeType = "122801"
	// "122901" BVF Cluj Napoca Aero (ROCJ1810)
	CustomsOfficeROCJ1810 CustomsOfficeCodeType = "122901"
	// "132901" BVF Constanţa Sud Agigea (ROCT1900)
	CustomsOfficeROCT1900 CustomsOfficeCodeType = "132901"
	// "132902" BVF Mihail Kogălniceanu (ROCT5100)
	CustomsOfficeROCT5100 CustomsOfficeCodeType = "132902"
	// "132903" BVF Mangalia (ROCT5400)
	CustomsOfficeROCT5400 CustomsOfficeCodeType = "132903"
	// "132904" BVF Constanţa Port (ROCT1970)
	CustomsOfficeROCT1970 CustomsOfficeCodeType = "132904"
	// "142801" BVI Sfântu Gheorghe (ROBV7820)
	CustomsOfficeROBV7820 CustomsOfficeCodeType = "142801"
	// "152801" BVI Târgoviște (ROBU8600)
	CustomsOfficeROBU8600 CustomsOfficeCodeType = "152801"
	// "162801" BVI Craiova (ROCR2100)
	CustomsOfficeROCR2100 CustomsOfficeCodeType = "162801"
	// "162901" BVF Craiova Aeroport (ROCR2110)
	CustomsOfficeROCR2110 CustomsOfficeCodeType = "162901"
	// "162902" BVF Bechet (ROCR1720)
	CustomsOfficeROCR1720 CustomsOfficeCodeType = "162902"
	// "162903" BVF Calafat (ROCR1700)
	CustomsOfficeROCR1700 CustomsOfficeCodeType = "162903"
	// "172901" BVF Zona Liberă Galaţi (ROGL3810)
	CustomsOfficeROGL3810 CustomsOfficeCodeType = "172901"
	// "172902" BVF Giurgiuleşti (ROGL3850)
	CustomsOfficeROGL3850 CustomsOfficeCodeType = "172902"
	// "172903" BVF Oancea (ROGL3610)
	CustomsOfficeROGL3610 CustomsOfficeCodeType = "172903"
	// "172904" BVF Galaţi (ROGL3800)
	CustomsOfficeROGL3800 CustomsOfficeCodeType = "172904"
	// "182801" BVI Târgu Jiu (ROCR8810)
	CustomsOfficeROCR8810 CustomsOfficeCodeType = "182801"
	// "192801" BVI Miercurea Ciuc (ROBV5600)
	CustomsOfficeROBV5600 CustomsOfficeCodeType = "192801"
	// "202801" BVI Deva (ROTM8100)
	CustomsOfficeROTM8100 CustomsOfficeCodeType = "202801"
	// "212801" BVI Slobozia (ROCT8220)
	CustomsOfficeROCT8220 CustomsOfficeCodeType = "212801"
	// "222901" BVF Iaşi Aero (ROIS4660)
	CustomsOfficeROIS4660 CustomsOfficeCodeType = "222901"
	// "222902" BVF Sculeni (ROIS4990)
	CustomsOfficeROIS4990 CustomsOfficeCodeType = "222902"
	// "222903" BVF Iaşi (ROIS4650)
	CustomsOfficeROIS4650 CustomsOfficeCodeType = "222903"
	// "232801" BVI Antrepozite/Ilfov (ROBU1200)
	CustomsOfficeROBU1200 CustomsOfficeCodeType = "232801"
	// "232901" BVF Otopeni Călători (ROBU1030)
	CustomsOfficeROBU1030 CustomsOfficeCodeType = "232901"
	// "242801" BVI Baia Mare (ROCJ0500)
	CustomsOfficeROCJ0500 CustomsOfficeCodeType = "242801"
	// "242901" BVF Aero Baia Mare (ROCJ0510)
	CustomsOfficeROCJ0510 CustomsOfficeCodeType = "242901"
	// "242902" BVF Sighet (ROCJ8000)
	CustomsOfficeROCJ8000 CustomsOfficeCodeType = "242902"
	// "252901" BVF Orşova (ROCR7280)
	CustomsOfficeROCR7280 CustomsOfficeCodeType = "252901"
	// "252902" BVF Porţile De Fier I (ROCR7270)
	CustomsOfficeROCR7270 CustomsOfficeCodeType = "252902"
	// "252903" BVF Porţile De Fier II (ROCR7200)
	CustomsOfficeROCR7200 CustomsOfficeCodeType = "252903"
	// "252904" BVF Drobeta Turnu Severin (ROCR9000)
	CustomsOfficeROCR9000 CustomsOfficeCodeType = "252904"
	// "262801" BVI Târgu Mureş (ROBV8800)
	CustomsOfficeROBV8800 CustomsOfficeCodeType = "262801"
	// "262901" BVF Târgu Mureş Aeroport (ROBV8820)
	CustomsOfficeROBV8820 CustomsOfficeCodeType = "262901"
	// "272801" BVI Piatra Neamţ (ROIS7400)
	CustomsOfficeROIS7400 CustomsOfficeCodeType = "272801"
	// "282801" BVI Corabia (ROCR2000)
	CustomsOfficeROCR2000 CustomsOfficeCodeType = "282801"
	// "282802" BVI Olt (ROCR8210)
	CustomsOfficeROCR8210 CustomsOfficeCodeType = "282802"
	// "292801" BVI Ploiești (ROBU7100)
	CustomsOfficeROBU7100 CustomsOfficeCodeType = "292801"
	// "302801" BVI Satu-Mare (ROCJ7810)
	CustomsOfficeROCJ7810 CustomsOfficeCodeType = "302801"
	// "302901" BVF Halmeu (ROCJ4310)
	CustomsOfficeROCJ4310 CustomsOfficeCodeType = "302901"
	// "302902" BVF Aeroport Satu Mare (ROCJ7830)
	CustomsOfficeROCJ7830 CustomsOfficeCodeType = "302902"
	// "312801" BVI Zalău (ROCJ9700)
	CustomsOfficeROCJ9700 CustomsOfficeCodeType = "312801"
	// "322801" BVI Sibiu (ROBV7900)
	CustomsOfficeROBV7900 CustomsOfficeCodeType = "322801"
	// "322901" BVF Sibiu Aeroport (ROBV7910)
	CustomsOfficeROBV7910 CustomsOfficeCodeType = "322901"
	// "332801" BVI Suceava (ROIS8230)
	CustomsOfficeROIS8230 CustomsOfficeCodeType = "332801"
	// "332901" BVF Dorneşti (ROIS2700)
	CustomsOfficeROIS2700 CustomsOfficeCodeType = "332901"
	// "332902" BVF Siret (ROIS8200)
	CustomsOfficeROIS8200 CustomsOfficeCodeType = "332902"
	// "332903" BVF Suceava Aero (ROIS8250)
	CustomsOfficeROIS8250 CustomsOfficeCodeType = "332903"
	// "332904" BVF Vicovu De Sus (ROIS9620)
	CustomsOfficeROIS9620 CustomsOfficeCodeType = "332904"
	// "342801" BVI Alexandria (ROCR0310)
	CustomsOfficeROCR0310 CustomsOfficeCodeType = "342801"
	// "342901" BVF Turnu Măgurele (ROCR9100)
	CustomsOfficeROCR9100 CustomsOfficeCodeType = "342901"
	// "342902" BVF Zimnicea (ROCR5800)
	CustomsOfficeROCR5800 CustomsOfficeCodeType = "342902"
	// "352802" BVI Timişoara Bază (ROTM8720)
	CustomsOfficeROTM8720 CustomsOfficeCodeType = "352802"
	// "352901" BVF Jimbolia (ROTM5010)
	CustomsOfficeROTM5010 CustomsOfficeCodeType = "352901"
	// "352902" BVF Moraviţa (ROTM5510)
	CustomsOfficeROTM5510 CustomsOfficeCodeType = "352902"
	// "352903" BVF Timişoara Aeroport (ROTM8730)
	CustomsOfficeROTM8730 CustomsOfficeCodeType = "352903"
	// "362901" BVF Sulina (ROCT8300)
	CustomsOfficeROCT8300 CustomsOfficeCodeType = "362901"
	// "362902" BVF Aeroport Delta Dunării Tulcea (ROGL8910)
	CustomsOfficeROGL8910 CustomsOfficeCodeType = "362902"
	// "362903" BVF Tulcea (ROGL8900)
	CustomsOfficeROGL8900 CustomsOfficeCodeType = "362903"
	// "362904" BVF Isaccea (ROGL8920)
	CustomsOfficeROGL8920 CustomsOfficeCodeType = "362904"
	// "372801" BVI Vaslui (ROIS9610)
	CustomsOfficeROIS9610 CustomsOfficeCodeType = "372801"
	// "372901" BVF Fălciu (-)
	// "372902" BVF Albiţa (ROIS0100)
	CustomsOfficeROIS0100 CustomsOfficeCodeType = "372902"
	// "382801" BVI Râmnicu Vâlcea (ROCR7700)
	CustomsOfficeROCR7700 CustomsOfficeCodeType = "382801"
	// "392801" BVI Focșani (ROGL3600)
	CustomsOfficeROGL3600 CustomsOfficeCodeType = "392801"
	// "402801" BVI Bucureşti Poştă (ROBU1380)
	CustomsOfficeROBU1380 CustomsOfficeCodeType = "402801"
	// "402802" BVI Târguri și Expoziții (ROBU1400)
	CustomsOfficeROBU1400 CustomsOfficeCodeType = "402802"
	// "402901" BVF Băneasa (ROBU1040)
	CustomsOfficeROBU1040 CustomsOfficeCodeType = "402901"
	// "512801" BVI Călăraşi (ROCT1710)
	CustomsOfficeROCT1710 CustomsOfficeCodeType = "512801"
	// "522801" BVI Giurgiu (ROBU3910)
	CustomsOfficeROBU3910 CustomsOfficeCodeType = "522801"
	// "522901" BVF Zona Liberă Giurgiu (ROBU3980)
	CustomsOfficeROBU3980 CustomsOfficeCodeType = "522901"
)

type CountyCodeType string

const (
	// "40" Municipiul Bucureşti
	CountyCodeB CountyCodeType = "40"
	// "1" Alba
	CountyCodeAB CountyCodeType = "1"
	// "2" Arad
	CountyCodeAR CountyCodeType = "2"
	// "3" Argeş
	CountyCodeAG CountyCodeType = "3"
	// "4" Bacău
	CountyCodeBC CountyCodeType = "4"
	// "5" Bihor
	CountyCodeBH CountyCodeType = "5"
	// "6" Bistriţa-Năsăud
	CountyCodeBN CountyCodeType = "6"
	// "7" Botoşani
	CountyCodeBT CountyCodeType = "7"
	// "8" Braşov
	CountyCodeBV CountyCodeType = "8"
	// "9" Brăila
	CountyCodeBR CountyCodeType = "9"
	// "10" Buzău
	CountyCodeBZ CountyCodeType = "10"
	// "11" Caraş-Severin
	CountyCodeCS CountyCodeType = "11"
	// "51" Călăraşi
	CountyCodeCL CountyCodeType = "51"
	// "12" Cluj
	CountyCodeCJ CountyCodeType = "12"
	// "13" Constanţa
	CountyCodeCT CountyCodeType = "13"
	// "14" Covasna
	CountyCodeCV CountyCodeType = "14"
	// "15" Dâmboviţa
	CountyCodeDB CountyCodeType = "15"
	// "16" Dolj
	CountyCodeDJ CountyCodeType = "16"
	// "17" Galaţi
	CountyCodeGL CountyCodeType = "17"
	// "52" Giurgiu
	CountyCodeGR CountyCodeType = "52"
	// "18" Gorj
	CountyCodeGJ CountyCodeType = "18"
	// "19" Harghita
	CountyCodeHR CountyCodeType = "19"
	// "20" Hunedoara
	CountyCodeHD CountyCodeType = "20"
	// "21" Ialomiţa
	CountyCodeIL CountyCodeType = "21"
	// "22" Iaşi
	CountyCodeIS CountyCodeType = "22"
	// "23" Ilfov
	CountyCodeIF CountyCodeType = "23"
	// "24" Maramureş
	CountyCodeMM CountyCodeType = "24"
	// "25" Mehedinţi
	CountyCodeMH CountyCodeType = "25"
	// "26" Mureş
	CountyCodeMS CountyCodeType = "26"
	// "27" Neamţ
	CountyCodeNT CountyCodeType = "27"
	// "28" Olt
	CountyCodeOT CountyCodeType = "28"
	// "29" Prahova
	CountyCodePH CountyCodeType = "29"
	// "30" Satu Mare
	CountyCodeSM CountyCodeType = "30"
	// "31" Sălaj
	CountyCodeSJ CountyCodeType = "31"
	// "32" Sibiu
	CountyCodeSB CountyCodeType = "32"
	// "33" Suceava
	CountyCodeSV CountyCodeType = "33"
	// "34" Teleorman
	CountyCodeTR CountyCodeType = "34"
	// "35" Timiş
	CountyCodeTM CountyCodeType = "35"
	// "36" Tulcea
	CountyCodeTL CountyCodeType = "36"
	// "37" Vaslui
	CountyCodeVS CountyCodeType = "37"
	// "38" Vâlcea
	CountyCodeVL CountyCodeType = "38"
	// "39" Vrancea
	CountyCodeVN CountyCodeType = "39"
)

// RoCountyNameToCountyCode returns the county code for a Romanian county name.
// eg. "bucurești" -> "40"
func RoCountyNameToCountyCode(name string) (_ CountyCodeType, ok bool) {
	switch strings.ToLower(text.Transliterate(name)) {
	case "bucuresti", "municipiul bucuresti":
		return CountyCodeB, true
	case "alba":
		return CountyCodeAB, true
	case "arad":
		return CountyCodeAR, true
	case "arges":
		return CountyCodeAG, true
	case "bacau":
		return CountyCodeBC, true
	case "bihor":
		return CountyCodeBH, true
	case "bistrita-nasaud":
		return CountyCodeBN, true
	case "botosani":
		return CountyCodeBT, true
	case "braila":
		return CountyCodeBR, true
	case "brasov":
		return CountyCodeBV, true
	case "buzau":
		return CountyCodeBZ, true
	case "calarasi":
		return CountyCodeCL, true
	case "caras-severin":
		return CountyCodeCS, true
	case "cluj":
		return CountyCodeCJ, true
	case "constanta":
		return CountyCodeCT, true
	case "covasna":
		return CountyCodeCV, true
	case "dambovita":
		return CountyCodeDB, true
	case "dolj":
		return CountyCodeDJ, true
	case "galati":
		return CountyCodeGL, true
	case "giurgiu":
		return CountyCodeGR, true
	case "gorj":
		return CountyCodeGJ, true
	case "harghita":
		return CountyCodeHR, true
	case "hunedoara":
		return CountyCodeHD, true
	case "ialomita":
		return CountyCodeIL, true
	case "iasi":
		return CountyCodeIS, true
	case "ilfov":
		return CountyCodeIF, true
	case "maramures":
		return CountyCodeMM, true
	case "mehedinti":
		return CountyCodeMH, true
	case "mures":
		return CountyCodeMS, true
	case "neamt":
		return CountyCodeNT, true
	case "olt":
		return CountyCodeOT, true
	case "prahova":
		return CountyCodePH, true
	case "salaj":
		return CountyCodeSJ, true
	case "satu mare":
		return CountyCodeSM, true
	case "sibiu":
		return CountyCodeSB, true
	case "suceava":
		return CountyCodeSV, true
	case "teleorman":
		return CountyCodeTR, true
	case "timis":
		return CountyCodeTM, true
	case "tulcea":
		return CountyCodeTL, true
	case "vaslui":
		return CountyCodeVS, true
	case "valcea":
		return CountyCodeVL, true
	case "vrancea":
		return CountyCodeVN, true
	}
	return
}

type OpPurposeCodeType string

const (
	// "101" Commercialization / Comercializare
	OpPurposeCodeTypeCommercialization OpPurposeCodeType = "101"
	// "201" Production / Producție
	OpPurposeCodeTypeProduction OpPurposeCodeType = "201"
	// "301" Gratuities / Gratuități
	OpPurposeCodeTypeGratuities OpPurposeCodeType = "301"
	// "401" Commercial equipment / Echipament comercial
	OpPurposeCodeTypeCommercialEquipment OpPurposeCodeType = "401"
	// "501" Fixed assets / Mijloace fixe
	OpPurposeCodeTypeFixedAssets OpPurposeCodeType = "501"
	// "601" Own consumption / Consum propriu
	OpPurposeCodeOwnConsumption OpPurposeCodeType = "601"
	// "703" Delivery ops with installation / Operațiuni de livrare cu instalare
	OpPurposeCodeDeliveryWithInstallation OpPurposeCodeType = "703"
	// "704" Transfer between management units / Transfer între gestiuni
	OpPurposeCodeTransfer OpPurposeCodeType = "704"
	// "705" Goods at customer disposal / Bunuri puse la dispoziția clientului
	OpPurposeGoodsAtCustomersDisposal OpPurposeCodeType = "705"
	// "801" Financial/operation leasing / Leasing financiar/operațional
	OpPurposeCodeLeasing OpPurposeCodeType = "801"
	// "802" Goods under warranty / Bunuri în garanție
	OpPurposeCodeGoodsUnderWarranty OpPurposeCodeType = "802"
	// "901" Exempted operations / Operațiuni scutite
	OpPurposeCodeExemptOperations OpPurposeCodeType = "901"
	// "1001" Ongoing / investment Investiție in curs
	OpPurposeCodeOngoingInvestment OpPurposeCodeType = "1001"
	// "1101" Donations, AID / Donații, ajutoare
	OpPurposeCodeDonations OpPurposeCodeType = "1101"
	// "9901" Others / Altele
	OpPurposeCodeOthers OpPurposeCodeType = "9901"
	// 9999" Same as the operation / Același cu operațiunea
	OpPurposeSameAsOperation OpPurposeCodeType = "9999"
)

type UnitCodeType = units.UnitCodeType

type DocumentType string

const (
	// "10" CMR
	DocumentTypeCMR DocumentType = "10"
	// "20" Factura
	DocumentTypeInvoice DocumentType = "20"
	// "30" Aviz de însoțire a mărfii
	DocumentTypeDeliveryNote DocumentType = "30"
	// "9999" Altele
	DocumentTypeOther DocumentType = "9999"
)

type ConfirmationType string

const (
	ConfirmationTypeConfirmed          ConfirmationType = "10"
	ConfirmationTypePartiallyConfirmed ConfirmationType = "20"
	ConfirmationTypeUnconfirmed        ConfirmationType = "30"
)
