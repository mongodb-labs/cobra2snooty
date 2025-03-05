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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/spf13/cobra"
)

func emptyRun(*cobra.Command, []string) {}

var rootCmd *cobra.Command
var echoCmd *cobra.Command

func Root() *cobra.Command {
	if rootCmd != nil {
		return rootCmd
	}
	rootCmd = &cobra.Command{
		Use:   "root",
		Short: "Root short description",
		Long:  "Root long description",
		Run:   emptyRun,
	}
	rootCmd.PersistentFlags().StringP("rootflag", "r", "two", "")
	rootCmd.PersistentFlags().StringP("strtwo", "t", "two", "help message for parent flag strtwo")

	printCmd := &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Long:  `an absolutely utterly useless command for testing.`,
	}
	printCmd.PersistentFlags().StringP("strthree", "s", "three", "help message for flag strthree")
	printCmd.Flags().IntP("intthree", "i", 345, "help message for flag intthree")
	printCmd.Flags().BoolP("boolthree", "b", true, "help message for flag boolthree")

	dummyCmd := &cobra.Command{
		Use:   "dummy [action]",
		Short: "Performs a dummy action",
	}

	rootCmd.AddCommand(printCmd, Echo(), dummyCmd)
	return rootCmd
}

func Echo() *cobra.Command {
	if echoCmd != nil {
		return echoCmd
	}
	echoCmd = &cobra.Command{
		Use:     "echo <string to print> [test param]",
		Aliases: []string{"say"},
		Short:   "Echo anything to the screen",
		Long:    "an utterly useless command for testing",
		Example: "# Example with intro text\n atlas command no intro text\n",
		Annotations: map[string]string{
			"string to printDesc": "A string to print",
			"test paramDesc":      "just for testing",
		},
	}
	echoCmd.PersistentFlags().StringP("strone", "s", "one", "help message for flag strone")
	echoCmd.PersistentFlags().BoolP("persistentbool", "p", false, "help message for flag persistentbool")
	echoCmd.Flags().IntP("intone", "i", 123, "help message for flag intone")
	echoCmd.Flags().BoolP("boolone", "b", true, "help message for flag boolone")
	echoCmd.Flags().StringToStringP("stringtostring", "x", nil, "help message for flag stringtostring")

	timesCmd := &cobra.Command{
		Use:        "times [# times] [string to echo]",
		SuggestFor: []string{"counts"},
		Short:      "Echo anything to the screen more times",
		Long:       `a slightly useless command for testing.`,
		Run:        emptyRun,
	}
	timesCmd.PersistentFlags().StringP("strtwo", "t", "2", "help message for child flag strtwo")
	timesCmd.Flags().IntP("inttwo", "j", 234, "help message for flag inttwo")
	timesCmd.Flags().BoolP("booltwo", "c", false, "help message for flag booltwo")

	echoCmd.AddCommand(timesCmd, EchoSubCmd(), deprecatedCmd)
	return echoCmd
}

var echoSubCmd *cobra.Command

func EchoSubCmd() *cobra.Command {
	if echoSubCmd != nil {
		return echoSubCmd
	}
	echoSubCmd = &cobra.Command{
		Use:   "echosub [string to print]",
		Short: "second sub command for echo",
		Long:  "an absolutely utterly useless command for testing gendocs!.",
		Run:   emptyRun,
	}
	return echoSubCmd
}

var deprecatedCmd = &cobra.Command{
	Use:        "deprecated [can't do anything here]",
	Short:      "A command which is deprecated",
	Long:       `an absolutely utterly useless command for testing deprecation!.`,
	Deprecated: "Please use echo instead",
}

func TestGenDocs(t *testing.T) {
	// We generate on a subcommand, so we have both subcommands and parents
	buf := new(bytes.Buffer)
	Root() // init root
	if err := GenDocs(Echo(), buf); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	checkStringContains(t, output, Echo().Long)
	checkStringContains(t, output, `# Example with intro text
  atlas command no intro text
`)
	checkStringContains(t, output, "boolone")
	checkStringContains(t, output, "rootflag")
	//
	checkStringContains(t, output, fmt.Sprintf("   * - string to print\n     - string\n     - true\n     - %s\n", Echo().Annotations["string to printDesc"]))
	checkStringContains(t, output, fmt.Sprintf("   * - test param\n     - string\n     - false\n     - %s\n", Echo().Annotations["test paramDesc"]))
	checkStringOmits(t, output, Root().Short)
	checkStringContains(t, output, EchoSubCmd().Short)
	checkStringOmits(t, output, deprecatedCmd.Short)

	// Verify that the text "This value defaults to" is not printed when the default value is provided to StringToStringP
	checkStringContains(t, output, "* - -x, --stringtostring\n     - key=value\n     - false\n     - help message for flag stringtostring\n   *")
}

func TestGenDocsNoHiddenParents(t *testing.T) {
	// We generate on a subcommand so we have both subcommands and parents
	for _, name := range []string{"rootflag", "strtwo"} {
		f := Root().PersistentFlags().Lookup(name)
		f.Hidden = true
		t.Cleanup(func() { f.Hidden = false })
	}
	buf := new(bytes.Buffer)
	if err := GenDocs(Echo(), buf); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	checkStringContains(t, output, Echo().Long)
	checkStringContains(t, output, `# Example with intro text
  atlas command no intro text
`)
	checkStringContains(t, output, "boolone")
	checkStringOmits(t, output, "rootflag")
	checkStringOmits(t, output, Root().Short)
	checkStringContains(t, output, Echo().Short)
	checkStringOmits(t, output, deprecatedCmd.Short)
	checkStringOmits(t, output, "Options inherited from parent commands")
}

func TestGenDocsNoTag(t *testing.T) {
	Root().DisableAutoGenTag = true
	defer func() { Root().DisableAutoGenTag = false }()

	buf := new(bytes.Buffer)
	if err := GenDocs(Root(), buf); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	unexpected := "Auto generated"
	checkStringOmits(t, output, unexpected)
}

func TestGenTreeDocs(t *testing.T) {
	c := &cobra.Command{
		Use: "do <arg1> [arg2]",
		Annotations: map[string]string{
			"arg1Desc": "desc",
			"arg2Desc": "desc",
		},
	}

	tmpdir, err := os.MkdirTemp("", "test-gen-rst-tree")
	if err != nil {
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)

	if err := GenTreeDocs(c, tmpdir); err != nil {
		t.Fatalf("GenTreeDocs failed: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(tmpdir, "do.txt")); err != nil {
		t.Fatalf("Expected file 'do.txt' to exist")
	}
}

func BenchmarkGenDocsToFile(b *testing.B) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := GenDocs(Root(), file); err != nil {
			b.Fatal(err)
		}
	}
}

func checkStringContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

func checkStringOmits(t *testing.T, got, expected string) {
	t.Helper()
	if strings.Contains(got, expected) {
		t.Errorf("Expected to not contain: \n %v\nGot: %v", expected, got)
	}
}

func TestArgsRegex(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		result := argsRegex.FindAllString("<arg1> [arg2]", -1)
		expected := []string{"<arg1>", "[arg2]"}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("expected: %s, got: %s\n", expected[i], result[i])
			}
		}
	})
	t.Run("with spaces", func(t *testing.T) {
		result := argsRegex.FindAllString("<this arg1> [that arg2]", -1)
		expected := []string{"<this arg1>", "[that arg2]"}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("expected: %s, got: %s\n", expected[i], result[i])
			}
		}
	})
	t.Run("repeating", func(t *testing.T) {
		result := argsRegex.FindAllString("<arg1>... [arg2]...", -1)
		expected := []string{"<arg1>", "[arg2]"}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("expected: %s, got: %s\n", expected[i], result[i])
			}
		}
	})
	t.Run("empty", func(t *testing.T) {
		result := argsRegex.FindAllString("<> []", -1)
		if len(result) != 0 {
			t.Fatalf("expected no matches\n")
		}
	})
	t.Run("complex", func(t *testing.T) {
		result := argsRegex.FindAllString("<this arg1> <that arg2> [optional] [long option]", -1)
		expected := []string{"<this arg1>", "<that arg2>", "[optional]", "[long option]"}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("expected: %s, got: %s\n", expected[i], result[i])
			}
		}
	})
}

// To update snapshots run: UPDATE_SNAPSHOTS=true go test ./...
func TestGenDocsSnapshots(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		cmd     *cobra.Command
		options []GenDocsOption
	}{
		{
			name: "default_example",
			cmd: &cobra.Command{
				Use:     "default_example",
				Long:    "Testing example output using the default example formatter",
				Example: "example --test -v",
			},
			options: []GenDocsOption{
				WithCustomTimeGetter(func() time.Time {
					return time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC)
				}),
			},
		},
		{
			name: "custom_formatter_example",
			cmd: &cobra.Command{
				Use:     "custom_example",
				Long:    "Testing example output using a custom example formatter",
				Example: "example --test -v",
			},
			options: []GenDocsOption{
				WithCustomTimeGetter(func() time.Time {
					return time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC)
				}),
				WithCustomExampleFormatter(func(buf *bytes.Buffer, cmd *cobra.Command) {
					_, _ = fmt.Fprintf(buf, "custom example for %s\n", cmd.Use)
					_, _ = buf.WriteString(cmd.Example)
				}),
			},
		},
		{
			name: "custom_formatter_surround_default_example",
			cmd: &cobra.Command{
				Use:     "custom_formatter_surround_default_example",
				Long:    "Testing example output using a custom example formatter which calls the default formatter",
				Example: "example --test -v",
			},
			options: []GenDocsOption{
				WithCustomTimeGetter(func() time.Time {
					return time.Date(2025, 3, 5, 17, 0, 0, 0, time.UTC)
				}),
				WithCustomExampleFormatter(func(buf *bytes.Buffer, cmd *cobra.Command) {
					_, _ = buf.WriteString("-- before example --\n")
					DefaultExampleFormatter(buf, cmd)
					_, _ = buf.WriteString("-- after example --\n")
				}),
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshotter := cupaloy.New(cupaloy.SnapshotFileExtension(".txt"))
			// Create buffer to capture output
			buf := new(bytes.Buffer)

			// Execute function
			err := GenDocs(tt.cmd, buf, tt.options...)
			if err != nil {
				t.Errorf("GenDocs() error = %v", err)
				return
			}

			// Compare with snapshot
			err = snapshotter.SnapshotWithName(tt.name, buf.String())
			if err != nil {
				t.Errorf("Snapshot comparison failed: %v", err)
			}
		})
	}
}
