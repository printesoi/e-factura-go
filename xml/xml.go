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

package xml

import (
	"bytes"
	"io"

	"github.com/printesoi/xml-go"
)

// MarshalXML returns the XML encoding of v in Canonical XML form [XML-C14N].
// This method must be used for marshaling objects from this library, instead
// of encoding/xml. This method does NOT include the XML header declaration.
func MarshalXML(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// MarshalXMLWithHeader same as MarshalXML, but also add the XML header
// declaration.
func MarshalXMLWithHeader(v any) ([]byte, error) {
	data, err := MarshalXML(v)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), data...), nil
}

// MarshalIndentXML works like MarshalXML, but each XML element begins on a new
// indented line that starts with prefix and is followed by one or more
// copies of indent according to the nesting depth. This method does NOT
// include the XML header declaration.
func MarshalIndentXML(v any, prefix, indent string) ([]byte, error) {
	return xml.MarshalIndent(v, prefix, indent)
}

// MarshalIndentXMLWithHeader same as MarshalIndentXML, but also add the XML
// header declaration.
func MarshalIndentXMLWithHeader(v any, prefix, indent string) ([]byte, error) {
	data, err := MarshalIndentXML(v, prefix, indent)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), data...), nil
}

// Unmarshal parses the XML-encoded data and stores the result in
// the value pointed to by v, which must be an arbitrary struct,
// slice, or string. Well-formed data that does not fit into v is
// discarded. This method must be used for unmarshaling objects from this
// library, instead of encoding/xml.
func UnmarshalXML(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

// UnmarshalReaderXML reads all the content from the given reader r and
// unmarshals the data as XML into the value v.
func UnmarshalReaderXML(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return UnmarshalXML(data, v)
}

// MarshalXMLToReader returns the XML encoding of v as a io.Reader.
func MarshalXMLToReader(v any) (r io.Reader, err error) {
	var b bytes.Buffer
	if _, err := b.WriteString(xml.Header); err != nil {
		return nil, err
	}
	enc := xml.NewEncoder(&b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return &b, nil
}
