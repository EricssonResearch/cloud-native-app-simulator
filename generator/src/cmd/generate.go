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
	"strconv"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [chain-file] [cluster-file] [k8s-readiness-probe]",
	Short: "This commands generates Kubernetes manifest files by using example files in chains and clusters directory and also uses k8s readiness probe time",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		chain := args[0]
		cluster := args[1]
		readinessProbe, err := strconv.Atoi(args[2])
		exitIfError(err)

		m, placement := generate.Parse(chain, cluster)
		generate.Create(m, placement, readinessProbe)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
