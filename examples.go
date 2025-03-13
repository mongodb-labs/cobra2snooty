// Copyright 2022 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cobra2snooty

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	examplesHeader = `Examples
--------

`
	identChar = " "
)

type ExampleFormatter func(w io.Writer, cmd *cobra.Command)

func DefaultExampleFormatter(w io.Writer, cmd *cobra.Command) {
	if cmd.Example != "" {
		printExamples(w, cmd)
	}
}

func WithCustomExampleFormatter(customFormatter ExampleFormatter) func(options *GenDocsOptions) {
	return func(options *GenDocsOptions) {
		options.exampleFormatter = customFormatter
	}
}

func printExamples(w io.Writer, cmd *cobra.Command) {
	// Create example substrings
	examplestrimmed := strings.TrimLeft(cmd.Example, " #")
	examples := strings.Split(examplestrimmed, "# ")
	_, _ = w.Write([]byte(examplesHeader))
	// If it has an example, print the header, then print each in a code block.
	for _, example := range examples[0:] {
		comment := ""
		if strings.Contains(cmd.Example, "#") {
			comment = " #"
		}
		_, _ = w.Write([]byte(`.. code-block::
   :copyable: false
`))
		_, _ = fmt.Fprintf(w, "\n  %s%s\n", comment, indentString(example, identChar))
	}
}
