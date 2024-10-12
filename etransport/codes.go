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

	"github.com/printesoi/e-factura-go/text"
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
	// EL - Greece
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

// "1" Petea (HU)
// "2" Borș (HU)
// "3" Vărșand (HU)
// "4" Nădlac (HU)
// "5" Calafat (BG)
// "6" Bechet (BG)
// "7" Turnu Măgurele (BG)
// "8" Zimnicea (BG)
// "9" Giurgiu (BG)
// "10" Ostrov (BG)
// "11" Negru Vodă (BG)
// "12" Vama Veche (BG)
// "13" Călărași (BG)
// "14" Corabia (BG)
// "15" Oltenița (BG)
// "16" Carei (HU)
// "17" Cenad (HU)
// "18" Episcopia Bihor (HU)
// "19" Salonta (HU)
// "20" Săcuieni (HU)
// "21" Turnu (HU)
// "22" Urziceni (HU)
// "23" Valea lui Mihai (HU)
// "24" Vladimirescu (HU)
// "25" Porțile de Fier 1 (RS)
// "26" Naidăș (RS)
// "27" Stamora Moravița (RS)
// "28" Jimbolia (RS)
// "29" Halmeu (UA)
// "30" Stânca Costești (MD)
// "31" Sculeni (MD)
// "32" Albița (MD)
// "33" Oancea (MD)
// "34" Galați Giurgiulești (MD)
// "35" Constanța Sud Agigea (-)
// "36" Siret (UA)
// "37" Nădlac 2 - A1 (HU)
// "38" Borș 2 - A3 (HU)
// TODO: enum values
type BCPCodeType string

// Valori posibile pentru câmpul codBirouVamal (BVI/F - Birou Vamal de Interior/Frontieră):
// "12801" BVI Alba Iulia (ROBV0300)
// "22801" BVI Arad (ROTM0200)
// "22901" BVF Arad Aeroport (ROTM0230)
// "22902" BVF Zona Liberă Curtici (ROTM2300)
// "32801" BVI Pitești (ROCR7000)
// "42801" BVI Bacău (ROIS0600)
// "42901" BVF Bacău Aeroport (ROIS0620)
// "52801" BVI Oradea (ROCJ6570)
// "52901" BVF Oradea Aeroport (ROCJ6580)
// "62801" BVI Bistriţa-Năsăud (ROCJ0400)
// "72801" BVI Botoşani (ROIS1600)
// "72901" BVF Stanca Costeşti (ROIS1610)
// "72902" BVF Rădăuţi Prut (ROIS1620)
// "82801" BVI Braşov (ROBV0900)
// "92901" BVF Zona Liberă Brăila (ROGL0710)
// "92902" BVF Brăila (ROGL0700)
// "102801" BVI Buzău (ROGL1500)
// "112801" BVI Reșița (ROTM7600)
// "112901" BVF Naidăș (ROTM6100)
// "122801" BVI Cluj Napoca (ROCJ1800)
// "122901" BVF Cluj Napoca Aero (ROCJ1810)
// "132901" BVF Constanţa Sud Agigea (ROCT1900)
// "132902" BVF Mihail Kogălniceanu (ROCT5100)
// "132903" BVF Mangalia (ROCT5400)
// "132904" BVF Constanţa Port (ROCT1970)
// "142801" BVI Sfântu Gheorghe (ROBV7820)
// "152801" BVI Târgoviște (ROBU8600)
// "162801" BVI Craiova (ROCR2100)
// "162901" BVF Craiova Aeroport (ROCR2110)
// "162902" BVF Bechet (ROCR1720)
// "162903" BVF Calafat (ROCR1700)
// "172901" BVF Zona Liberă Galaţi (ROGL3810)
// "172902" BVF Giurgiuleşti (ROGL3850)
// "172903" BVF Oancea (ROGL3610)
// "172904" BVF Galaţi (ROGL3800)
// "182801" BVI Târgu Jiu (ROCR8810)
// "192801" BVI Miercurea Ciuc (ROBV5600)
// "202801" BVI Deva (ROTM8100)
// "212801" BVI Slobozia (ROCT8220)
// "222901" BVF Iaşi Aero (ROIS4660)
// "222902" BVF Sculeni (ROIS4990)
// "222903" BVF Iaşi (ROIS4650)
// "232801" BVI Antrepozite/Ilfov (ROBU1200)
// "232901" BVF Otopeni Călători (ROBU1030)
// "242801" BVI Baia Mare (ROCJ0500)
// "242901" BVF Aero Baia Mare (ROCJ0510)
// "242902" BVF Sighet (ROCJ8000)
// "252901" BVF Orşova (ROCR7280)
// "252902" BVF Porţile De Fier I (ROCR7270)
// "252903" BVF Porţile De Fier II (ROCR7200)
// "252904" BVF Drobeta Turnu Severin (ROCR9000)
// "262801" BVI Târgu Mureş (ROBV8800)
// "262901" BVF Târgu Mureş Aeroport (ROBV8820)
// "272801" BVI Piatra Neamţ (ROIS7400)
// "282801" BVI Corabia (ROCR2000)
// "282802" BVI Olt (ROCR8210)
// "292801" BVI Ploiești (ROBU7100)
// "302801" BVI Satu-Mare (ROCJ7810)
// "302901" BVF Halmeu (ROCJ4310)
// "302902" BVF Aeroport Satu Mare (ROCJ7830)
// "312801" BVI Zalău (ROCJ9700)
// "322801" BVI Sibiu (ROBV7900)
// "322901" BVF Sibiu Aeroport (ROBV7910)
// "332801" BVI Suceava (ROIS8230)
// "332901" BVF Dorneşti (ROIS2700)
// "332902" BVF Siret (ROIS8200)
// "332903" BVF Suceava Aero (ROIS8250)
// "332904" BVF Vicovu De Sus (ROIS9620)
// "342801" BVI Alexandria (ROCR0310)
// "342901" BVF Turnu Măgurele (ROCR9100)
// "342902" BVF Zimnicea (ROCR5800)
// "352802" BVI Timişoara Bază (ROTM8720)
// "352901" BVF Jimbolia (ROTM5010)
// "352902" BVF Moraviţa (ROTM5510)
// "352903" BVF Timişoara Aeroport (ROTM8730)
// "362901" BVF Sulina (ROCT8300)
// "362902" BVF Aeroport Delta Dunării Tulcea (ROGL8910)
// "362903" BVF Tulcea (ROGL8900)
// "362904" BVF Isaccea (ROGL8920)
// "372801" BVI Vaslui (ROIS9610)
// "372901" BVF Fălciu (-)
// "372902" BVF Albiţa (ROIS0100)
// "382801" BVI Râmnicu Vâlcea (ROCR7700)
// "392801" BVI Focșani (ROGL3600)
// "402801" BVI Bucureşti Poştă (ROBU1380)
// "402802" BVI Târguri și Expoziții (ROBU1400)
// "402901" BVF Băneasa (ROBU1040)
// "512801" BVI Călăraşi (ROCT1710)
// "522801" BVI Giurgiu (ROBU3910)
// "522901" BVF Zona Liberă Giurgiu (ROBU3980)
// TODO: enum values
type CustomsOfficeCodeType string

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

// Câmpul codScopOperatiune ia valori diferite, în funcţie de valoarea câmpului codTipOperatiune, astfel:
// pentru codTipOperatiune = "10" AIC - Achiziţie intracomunitară:
//
//	"101" Comercializare
//	"201" Producție
//	"301" Gratuități
//	"401" Echipament comercial
//	"501" Mijloace fixe
//	"601" Consum propriu
//	"703" Operațiuni de livrare cu instalare
//	"801" Leasing financiar/operațional
//	"802" Bunuri în garanție
//	"901" Operațiuni scutite
//	"1001" Investiție in curs
//	"1101" Donații, ajutoare
//	"9901" Altele
//
// pentru codTipOperatiune = "12" LHI - Operațiuni în sistem lohn (UE) - intrare:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "14" SCI - Stocuri la dispoziția clientului (Call-off stock) - intrare:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "20" LIC - Livrare intracomunitară:
//
//	"101" Comercializare
//	"301" Gratuități
//	"703" Operațiuni de livrare cu instalare
//	"801" Leasing financiar/operațional
//	"802" Bunuri în garanție
//	"9901" Altele
//
// pentru codTipOperatiune = "22" LHE - Operațiuni în sistem lohn (UE) - ieşire:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "24" SCE - Stocuri la dispoziția clientului (Call-off stock) - ieşire:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "30" TTN - Transport pe teritoriul naţional:
//
//	"101" Comercializare
//	"704" Transfer între gestiuni
//	"705" Bunuri puse la dispoziția clientului
//	"9901" Altele
//
// pentru codTipOperatiune = "40" IMP - Import:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "50" EXP - Export:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "60" DIN - Tranzactie intracomunitara - Intrare pentru depozitare/formare nou transport:
//
//	"9999" Același cu operațiunea
//
// pentru codTipOperatiune = "70" DIE - Tranzactie intracomunitara - Iesire dupa depozitare/formare nou transport:
//
//	"9999" Același cu operațiunea
type OpPurposeCodeType string

const (
	// "101" Comercializare
	OpPurposeCodeTypeCommercialization OpPurposeCodeType = "101"
	// "201" Producție
	OpPurposeCodeTypeProduction OpPurposeCodeType = "201"
	// "301" Gratuități
	OpPurposeCodeTypeGratuities OpPurposeCodeType = "301"
	// "401" Echipament comercial
	OpPurposeCodeTypeCommercialEquipment OpPurposeCodeType = "401"
	// "501" Mijloace fixe
	OpPurposeCodeTypeFixedAssets OpPurposeCodeType = "501"
	// "601" Consum propriu
	OpPurposeCodeOwnConsumption OpPurposeCodeType = "601"
	// "703" Operațiuni de livrare cu instalare
	OpPurposeCodeDeliveryWithInstallation OpPurposeCodeType = "703"
	// "704" Transfer între gestiuni
	OpPurposeCodeTransfer OpPurposeCodeType = "704"
	// "705" Bunuri puse la dispoziția clientului
	OpPurposeGoodsAsCustomersDisposal OpPurposeCodeType = "705"
	// "801" Leasing financiar/operațional
	OpPurposeCodeLeasing OpPurposeCodeType = "801"
	// "802" Bunuri în garanție
	OpPurposeCodeGoodsUnderWarranty OpPurposeCodeType = "802"
	// "901" Operațiuni scutite
	OpPurposeCodeExemptOperations OpPurposeCodeType = "901"
	// "1001" Investiție in curs
	OpPurposeCodeOngoingInvestment OpPurposeCodeType = "1001"
	// "1101" Donații, ajutoare
	OpPurposeCodeDonations OpPurposeCodeType = "1101"
	// "9901" Altele
	OpPurposeCodeOthers OpPurposeCodeType = "9901"
	// 9999" Același cu operațiunea
	OpPurposeSameAsOperation OpPurposeCodeType = "9999"
)

// Valori posibile pentru câmpul codUnitateMasura: UN/ECE Recommendation N°20 and UN/ECE Recommendation N°21 — Unit codes
type UnitMeasureCodeType string

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
