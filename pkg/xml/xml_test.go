// Copyright 2024-2025 Victor Dodon
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

package xml_test

import (
	"fmt"
	"strings"
	"testing"

	pkgxml "github.com/printesoi/e-factura-go/pkg/xml"
	"github.com/printesoi/xml-go"
	"github.com/stretchr/testify/assert"
)

const (
	xmlnsFoo       = "urn:oasis:names:specification:ubl:schema:xsd:Foo"
	xmlnsFooPrefix = "foo"
	xmlnsBar       = "urn:oasis:names:specification:ubl:schema:xsd:Bar"
	xmlnsBarPrefix = "bar"
)

type Doc struct {
	ID      string `xml:"urn:oasis:names:specification:ubl:schema:xsd:Foo ID"`
	Payload struct {
		Data string `xml:"urn:oasis:names:specification:ubl:schema:xsd:Foo Data,omitempty"`
	} `xml:"urn:oasis:names:specification:ubl:schema:xsd:Bar Payload"`
	// Name of node.
	XMLName xml.Name `xml:"Doc"`
	// xmlns:foo attr. Will be automatically set in MarshalXML
	NamespaceFoo string `xml:"xmlns:foo,attr"`
	// xmlns:bar attr. Will be automatically set in MarshalXML
	NamespaceBar string `xml:"xmlns:bar,attr"`
}

func setupEncoder(enc *xml.Encoder) *xml.Encoder {
	enc.AddNamespaceBinding(xmlnsFoo, xmlnsFooPrefix)
	enc.AddSkipNamespaceAttrForPrefix(xmlnsFoo, xmlnsFooPrefix)
	enc.AddNamespaceBinding(xmlnsBar, xmlnsBarPrefix)
	enc.AddSkipNamespaceAttrForPrefix(xmlnsBar, xmlnsBarPrefix)
	return enc
}

func (d Doc) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	d.NamespaceFoo = xmlnsFoo
	d.NamespaceBar = xmlnsBar
	setupEncoder(e)
	// This allows us to strip the MarshalXML method.
	type doc Doc
	return e.EncodeElement(doc(d), start)
}

func TestMarshal(t *testing.T) {
	assert := assert.New(t)

	var doc0 Doc
	doc0.ID = "123"
	doc0.Payload.Data = "FDT"
	xmlData, err := pkgxml.MarshalXML(doc0)

	if !assert.NoError(err) {
		return
	}

	xmlStr := string(xmlData)
	fmt.Printf("%s\n", xmlStr)

	assert.True(strings.Contains(xmlStr, fmt.Sprintf("xmlns:%s=\"%s\"", xmlnsFooPrefix, xmlnsFoo)))
	assert.True(strings.Contains(xmlStr, fmt.Sprintf("xmlns:%s=\"%s\"", xmlnsBarPrefix, xmlnsBar)))

	var doc1 Doc
	if !assert.NoError(pkgxml.UnmarshalXML(xmlData, &doc1)) {
		return
	}
	assert.Equal(doc0.ID, doc1.ID)
	assert.Equal(doc0.Payload.Data, doc1.Payload.Data)
}
