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

	"github.com/spf13/cobra"
)

func requiredRole(buf *bytes.Buffer, cmd *cobra.Command) {
	if cmd.Annotations["requiredRole"] != "" {
		buf.WriteString("\nTo use this command, the requesting user or API key must have the " + cmd.Annotations["requiredRole"] + " role.\n")
	}
}
