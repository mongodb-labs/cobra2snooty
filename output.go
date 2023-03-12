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

// This function can return the output for all commands when the output template is added as an annotation in the command file

func printOutputCreate(buf *bytes.Buffer, cmd *cobra.Command) {
	removerange := strings.ReplaceAll(cmd.Annotations["output"], "{{range .Results}}", "")
	removeend := strings.ReplaceAll(removerange, "{{end}}", "")
	bracketsremoved1 := strings.ReplaceAll(removeend, "{{.", "<")
	bracketsremoved2 := strings.ReplaceAll(bracketsremoved1, "}}", ">")
	replacevariable := strings.ReplaceAll(bracketsremoved2, "%s", "<Name>")
	replacevariable2 := strings.ReplaceAll(replacevariable, "\n", "\n   ")
	w := new(tabwriter.Writer)
	w.Init(buf, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)

	if cmd.Annotations["output"] == "" {
		return
	}

	buf.WriteString(outputHeader)
	buf.WriteString(`
If the command succeeds, the CLI prints a message similar to the following and replaces the values in brackets with your values. The | symbol represents a horizontal tab.

.. code-block::

`)
	if strings.HasSuffix(cmd.Annotations["output"], "}"+"\n") {
		fmt.Fprintln(w, replacevariable2)
		w.Flush()
	} else {
		buf.WriteString("   ")
		fmt.Fprintln(w, replacevariable2)
		w.Flush()
	}
	buf.WriteString("\n")
}
