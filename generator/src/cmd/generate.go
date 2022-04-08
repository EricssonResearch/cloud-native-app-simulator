/*
Copyright 2021 Telefonaktiebolaget LM Ericsson AB

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
package cmd

import (
	"application-generator/src/pkg/generate"
	"application-generator/src/pkg/model"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

const (
	SvcMaxNumberDefault        = 10
	SvcReplicaMaxNumberDefault = 10
	SvcEpMaxNumberDefault      = 5
)

func yesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func expandInputString(names string) []string {
	list := strings.Split(names, ",")

	var subList []string
	for index, item := range list {
		item = strings.TrimSpace(item)
		clusterRange := strings.Split(item, ":")

		if len(clusterRange) == 2 {
			prefix := getPrefix(clusterRange[0])
			min, _ := strconv.Atoi(strings.TrimPrefix(clusterRange[0], prefix))
			max, _ := strconv.Atoi(strings.TrimPrefix(clusterRange[1], prefix))

			for i := min; i <= max; i++ {
				subList = append(subList, prefix+strconv.Itoa(i))
			}

			// Remove element from list since it was expanded
			copy(list[index:], list[index+1:])
			list[len(list)-1] = ""
			list = list[:len(list)-1]
		}
	}
	list = append(list, subList...)

	return list
}

func getPrefix(s string) string {
	prefix := ""
	for _, c := range s {
		if !unicode.IsDigit(c) {
			prefix = prefix + string(c)
		} else {
			break
		}
	}

	return prefix
}

var generateCmd = &cobra.Command{
	Use:   "generate [mode] [input-file]",
	Short: "This command can be run under two different modes: (i) 'random' mode which generates a random description file or (ii) 'preset' mode which generates Kubernetes manifest based on a description file in the input directory",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		mode := args[0]

		var inputFile string
		if mode == "random" {
			simpleMode := yesNoPrompt("Do you want to set simple configuration? (otherwise, extended)", true)

			clusterNames := stringPrompt("Enter your cluster names (separated by commas, consecutive numbers can be given as range using ':')")
			nsNames := stringPrompt("Enter your namespaces (separated by commas, consecutive names with same prefix can be given as range using ':')")

			svcMaxNumber := SvcMaxNumberDefault
			svcReplicaMaxNumber := SvcReplicaMaxNumberDefault
			svcEpMaxNumber := SvcEpMaxNumberDefault

			if !simpleMode {
				svcMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many services do you want to have? (influences fan-out)"))
				svcReplicaMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many service replicas do you want to have?"))
				svcEpMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many service endpoints do you want to have? (fan-in)"))
			}

			outputFileName := stringPrompt("What name you want the output file to have?")

			clusterList := expandInputString(clusterNames)
			nsList := expandInputString(nsNames)

			userConfig := model.UserConfig{
				Clusters:            clusterList,
				Namespaces:          nsList,
				SvcMaxNumber:        svcMaxNumber,
				SvcReplicaMaxNumber: svcReplicaMaxNumber,
				SvcEpMaxNumber:      svcEpMaxNumber,
				OutputFileName:      outputFileName,
			}

			inputFile = generate.CreateJsonInput(userConfig)
		} else if mode == "preset" {
			inputFile = args[1]
		}

		config, clusters := generate.Parse(inputFile)
		generate.CreateK8sYaml(config, clusters)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
