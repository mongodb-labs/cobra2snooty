package cobra2snooty

import (
	"bytes"
	"errors"
	"fmt"
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

var ErrMissingDescription = errors.New("missing description")

func printArgs(buf *bytes.Buffer, cmd *cobra.Command) error {
	if args, ok := cmd.Annotations["args"]; ok {
		buf.WriteString("Arguments\n")
		buf.WriteString("---------\n\n")
		buf.WriteString(optionsHeader)
		var requiredSlice []string
		if requiredArgs, hasRequired := cmd.Annotations["requiredArgs"]; hasRequired {
			requiredSlice = strings.Split(requiredArgs, ",")
		}

		for _, arg := range strings.Split(args, ",") {
			required := stringInSlice(requiredSlice, arg)
			if description, hasDescription := cmd.Annotations[arg+"Desc"]; hasDescription {
				line := fmt.Sprintf("   * - %s\n     - string\n     - %v\n     - %s", arg, required, description)
				buf.WriteString(line)
			} else {
				return fmt.Errorf("%w: %s", ErrMissingDescription, arg)
			}
		}
		buf.WriteString("\n\n")
	}
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
