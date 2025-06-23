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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const optionsHeader = `.. list-table::
   :header-rows: 1
   :widths: 20 10 10 60

   * - Name
     - Type
     - Required
     - Description
`

var (
	ErrMissingDescription = errors.New("missing description")
	argsRegex             = regexp.MustCompile(`<[^>]+>|\[[^]]+]`)
)

func printArgs(buf *bytes.Buffer, cmd *cobra.Command) error {
	u := argsRegex.FindAllString(cmd.Use, -1)
	if len(u) == 0 {
		return nil
	}
	buf.WriteString("Arguments\n")
	buf.WriteString("---------\n\n")
	buf.WriteString(optionsHeader)
	for _, a := range u {
		value := a[1 : len(a)-1]
		description, hasDescription := cmd.Annotations[value+"Desc"]
		if !hasDescription {
			return fmt.Errorf("%w: %s - %s", ErrMissingDescription, cmd.CommandPath(), value)
		}
		required := strings.HasPrefix(a, "<")
		line := fmt.Sprintf("   * - %s\n     - string\n     - %v\n     - %s\n", value, required, description)
		buf.WriteString(line)
	}
	buf.WriteString("\n")

	return nil
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command) {
	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		buf.WriteString("Options\n")
		buf.WriteString("-------\n\n")
		buf.WriteString(optionsHeader)
		buf.WriteString(indentString(FlagUsages(flags), " "))
		buf.WriteString("\n")
	}

	parentFlags := cmd.InheritedFlags()
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("Inherited Options\n")
		buf.WriteString("-----------------\n\n")
		buf.WriteString(optionsHeader)
		buf.WriteString(indentString(FlagUsages(parentFlags), " "))
		buf.WriteString("\n")
	}
}
