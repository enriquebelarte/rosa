/*
Copyright (c) 2020 Red Hat, Inc.

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

package version

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/cmd/verify/rosa"
	"github.com/openshift/rosa/pkg/info"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the tool",
	Long:  "Prints the version number of the tool.",
	Run:   run,
}

func run(cmd *cobra.Command, argv []string) {
	fmt.Fprintf(os.Stdout, "%s\n", info.Version)
	rosa.Cmd.Run(rosa.Cmd, []string{})
}
