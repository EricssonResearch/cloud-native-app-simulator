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

var generateCmd = &cobra.Command{
	Use:   "generate [mode] [input-file]",
	Short: "This command can be run under two different modes: (i) 'random' mode which generates a random description file or (ii) 'preset' mode which generates Kubernetes manifest based on a description file in the input directory",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		mode := args[0]

		var inputFile string
		if mode == "random" {
			simpleMode := yesNoPrompt("Do you want to set simple configuration? (otherwise, extended)", true)

			// NOTE: Here we assume numerically consecutive names for clusters and namespaces
			clusterNamePrefix := stringPrompt("What is your cluster name prefix?")
			clusterNumber, _ := strconv.Atoi(stringPrompt("How many clusters do you have?"))
			nsNamePrefix := stringPrompt("What is your namespace prefix?")
			nsNumber, _ := strconv.Atoi(stringPrompt("How many namespaces do you have?"))

			svcMaxNumber := SvcMaxNumberDefault
			svcReplicaMaxNumber := SvcReplicaMaxNumberDefault
			svcEpMaxNumber := SvcEpMaxNumberDefault

			if !simpleMode {
				svcMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many services do you want to have? (influences fan-out)"))
				svcReplicaMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many service replicas do you want to have?"))
				svcEpMaxNumber, _ = strconv.Atoi(stringPrompt("Up to how many service endpoints do you want to have? (fan-in)"))
			}

			outputFileName := stringPrompt("What name you want the output file to have?")

			var clusters, namespaces []string

			for i := 1; i <= clusterNumber; i++ {
				clusters = append(clusters, clusterNamePrefix+strconv.Itoa(i))
			}

			for j := 1; j <= nsNumber; j++ {
				namespaces = append(namespaces, nsNamePrefix+strconv.Itoa(j))
			}

			userConfig := model.UserConfig{
				Clusters:            clusters,
				Namespaces:          namespaces,
				ClusterNamePrefix:   clusterNamePrefix,
				ClusterNumber:       clusterNumber,
				NsNamePrefix:        nsNamePrefix,
				NsNumber:            nsNumber,
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
