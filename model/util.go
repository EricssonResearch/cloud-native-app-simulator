/*
Copyright 2023 Telefonaktiebolaget LM Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"strings"
	"unicode"
)

// Translate K8s name into Go friendly name (example: endpoint-1 -> Endpoint1)
func GoName(name string) string {
	builder := strings.Builder{}
	nextUpper := true

	for _, c := range name {
		if c == '-' {
			nextUpper = true
			continue
		} else if nextUpper {
			c = unicode.ToUpper(c)
			nextUpper = false
		} else {
			c = unicode.ToLower(c)
		}

		builder.WriteRune(c)
	}

	return builder.String()
}
