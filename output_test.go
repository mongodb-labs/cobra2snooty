// Copyright 2023 MongoDB Inc
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
	"testing"

	"github.com/spf13/cobra"
)

func TestPrintOutputCreate(t *testing.T) {
	t.Run("replaceWithValueOrDefault", func(t *testing.T) {
		outputTemplate := `ID	NAME	DATABASE	COLLECTION	TYPE   
{{range .}}{{.IndexID}}	%s	{{.Database}}	{{.CollectionName}}	{{if .Type }}{{.Type}}{{else}}defaultValue{{end}}
{{end}}`

		expected := outputHeader + outputDescription +
			`   ID          NAME     DATABASE     COLLECTION         TYPE
   <IndexID>   <Name>   <Database>   <CollectionName>   <Type>
   

`

		cmd = &cobra.Command{
			Annotations: map[string]string{
				"output": outputTemplate,
			},
		}

		buf := new(bytes.Buffer)
		printOutputCreate(buf, cmd)
		result := buf.String()

		if result != expected {
			t.Errorf("expected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})
}
