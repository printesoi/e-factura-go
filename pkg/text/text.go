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

// Package text provides useful functions for processing text used in payloads.
package text

import (
	"github.com/alexsergivan/transliterator"
)

var (
	tr *transliterator.Transliterator
)

func init() {
	tr = transliterator.NewTransliterator(nil)
}

// Transliterate performs Unicode -> ASCII transliteration of the input text,
// basically transforms diacritics and characters from other alphabets to their
// ASCII variant (ă -> a, â -> a, é -> e, 池 -> Chi, ы -> y, etc).
// This method is useful for strings that may contain diacritics that must be
// used in an Invoice, like an Item name or Address line1, since ANAF discards
// non-ASCII characters.
func Transliterate(s string) string {
	return tr.Transliterate(s, "ro")
}
