// Copyright 2021 MongoDB Inc
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

// Package cobra2snooty was mostly inspired by https://github.com/spf13/cobra/tree/master/doc
// but with some changes to match the expected formats and styles of our writers and tools.
package cobra2snooty

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	separator        = "-"
	defaultExtension = ".txt"
)

// GenTreeDocs generates the docs for the full tree of commands.
func GenTreeDocs(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenTreeDocs(c, dir); err != nil {
			return err
		}
	}

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", separator) + defaultExtension
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return GenDocs(cmd, f)
}

const toc = `
.. default-domain:: mongodb

.. contents:: On this page
   :local:
   :backlinks: none
   :depth: 1
   :class: singlecol
`

const syntaxHeader = `Syntax
------

.. code-block::
`

const examplesHeader = `Examples
--------

.. code-block::
`

const tocHeader = `
.. toctree::
   :titlesonly:
`

// GenDocs creates snooty help output.
// Adapted from https://github.com/spf13/cobra/tree/master/doc to match MongoDB tooling and style.
func GenDocs(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	ref := strings.ReplaceAll(name, " ", separator)

	buf.WriteString(".. _" + ref + ":\n\n")
	buf.WriteString(strings.Repeat("=", len(name)) + "\n")
	buf.WriteString(name + "\n")
	buf.WriteString(strings.Repeat("=", len(name)) + "\n")
	buf.WriteString(toc)
	buf.WriteString("\n" + cmd.Short + "\n")
	if long := cmd.Long; long != "" {
		if strings.Contains(name, "completion bash") {
			long = bashCompletionLong(cmd)
		}
		buf.WriteString("\n" + long + "\n")
	}
	buf.WriteString("\n")

	if cmd.Runnable() {
		buf.WriteString(syntaxHeader)
		buf.WriteString(fmt.Sprintf("\n   %s\n\n", strings.ReplaceAll(cmd.UseLine(), "[flags]", "[options]")))
	}
	if err := printArgs(buf, cmd); err != nil {
		return err
	}
	printOptions(buf, cmd)

	if len(cmd.Example) > 0 {
		buf.WriteString(examplesHeader)
		buf.WriteString(fmt.Sprintf("\n%s\n\n", indentString(cmd.Example, " ")))
	}

	if hasRelatedCommands(cmd) {
		buf.WriteString("Related Commands\n")
		buf.WriteString("----------------\n\n")

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			ref = strings.ReplaceAll(cname, " ", separator)
			buf.WriteString(fmt.Sprintf("* :ref:`%s` - %s\n", ref, child.Short))
		}
		buf.WriteString("\n")
	}
	if _, ok := cmd.Annotations["toc"]; ok || !cmd.Runnable() {
		buf.WriteString(tocHeader)
		buf.WriteString("\n")
		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			ref = strings.ReplaceAll(cname, " ", separator)
			buf.WriteString(fmt.Sprintf("   %s </command/%s>\n", child.Name(), ref))
		}
		buf.WriteString("\n")
	}

	if !cmd.DisableAutoGenTag {
		buf.WriteString("*Auto generated by cobra2snooty on " + time.Now().Format("2-Jan-2006") + "*\n")
	}
	_, err := buf.WriteTo(w)
	return err
}

func bashCompletionLong(cmd *cobra.Command) string {
	return fmt.Sprintf(`
Generate the autocompletion script for the bash shell.
This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.
To load completions in your current shell session:
$ source <(%[1]s completion bash)
To load completions for every new session, execute once:
Linux:
$ %[1]s completion bash > /etc/bash_completion.d/%[1]s
MacOS:
$ %[1]s completion bash > /usr/local/etc/bash_completion.d/%[1]s
You will need to start a new shell for this setup to take effect.
`, cmd.Root().Name())
}

// Test to see if we have a reason to print See Also information in docs
// Basically this is a test for a parent command or a subcommand which is
// both not deprecated and not the autogenerated help command.
func hasRelatedCommands(cmd *cobra.Command) bool {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
