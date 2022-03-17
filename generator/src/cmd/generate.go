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
	"github.com/spf13/cobra"
	"application-generator/src/pkg/model"
)

var generateCmd = &cobra.Command{
	Use:   "generate [mode] [input-file]",
	Short: "This commands can be run under two different modes: (i) 'random' mode which generates a random description file or (ii) 'preset' mode which generates Kubernetes manifest based on a description file in the input directory",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		mode := args[0]

		var inputFile string
		if mode == "random" {
			// TODO: Change this hard-coded cluster configuration for actual user inputs
			clusterConfig := model.ClusterConfig{
				Clusters: 	[]string{"cluster-1", "cluster-2", "cluster-3", "cluster-4", "cluster-5"},
				Namespaces: []string{"ns-1", "ns-2", "ns-3"},
			}

			inputFile = generate.CreateJsonInput(clusterConfig)
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
