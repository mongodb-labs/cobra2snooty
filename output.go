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
	"text/tabwriter"

	"github.com/spf13/cobra"
)

const (
	outputHeader = `Output
------
`
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
)

// This function can return the output for all commands when the output template is added as an annotation in the command file

func printOutputCreate(buf *bytes.Buffer, cmd *cobra.Command) {
	if cmd.Annotations["output"] == "" {
		return
	}

	output := strings.ReplaceAll(cmd.Annotations["output"], "{{range .Results}}", "")
	output = strings.ReplaceAll(output, "{{end}}", "")
	output = strings.ReplaceAll(output, "{{.", "<")
	output = strings.ReplaceAll(output, "}}", ">")
	output = strings.ReplaceAll(output, "%s", "<Name>")
	output = strings.Replace(output, "   ", "", 1)
	output = strings.ReplaceAll(output, "\n", "\n   ")
	w := new(tabwriter.Writer)
	w.Init(buf, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)

	buf.WriteString(outputHeader)
	buf.WriteString(`
If the command succeeds, the CLI returns output similar to the following sample. Values in brackets represent your values.

.. code-block::

`)
	fmt.Fprintln(w, "   "+output)
	w.Flush()
	buf.WriteString("\n")
}
