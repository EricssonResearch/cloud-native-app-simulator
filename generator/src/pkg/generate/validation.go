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

package generate

import "application-generator/src/pkg/model"

// Occurrences returns the number of occurences of every value in a slice of strings
func Occurrences(strSlice []string) map[string]int {
	occurences := make(map[string]int)
	for _, entry := range strSlice {
		occurences[entry]++
	}
	return occurences
}

func IsDuplicateName(name string, allNames []string) bool {
	occurences := Occurrences(allNames)
	return occurences[name] > 1
}

// Validates an input JSON config provided by the user
func ValidateFileConfig(config *model.FileConfig) error {
	// TODO
	return nil
}
