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
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const examplesHeader = `Examples
--------

`

func printExamples(buf *bytes.Buffer, cmd *cobra.Command) error {
	// Create example substrings
	examplestrimmed := strings.TrimLeft(cmd.Example, "  #")
	examples := strings.Split(examplestrimmed, "# ")
	buf.WriteString(examplesHeader)
	// If it has an example, print the header, then print each in a code block.
	for _, example := range examples[0:] {
		if !strings.Contains(cmd.Example, "#") {
			buf.WriteString(`.. code-block::
			`)
			buf.WriteString(fmt.Sprintf("\n   %s\n", indentString(example, " ")))
		}
		if strings.Contains(cmd.Example, "#") {

			buf.WriteString(`.. code-block::
	`)
			buf.WriteString(fmt.Sprintf("\n   #%s\n", indentString(example, " ")))
		}
	}

	return nil
}