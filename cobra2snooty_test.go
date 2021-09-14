package cobra2snooty

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
		Use:     "echo [string to echo]",
		Aliases: []string{"say"},
		Short:   "Echo anything to the screen",
		Long:    "an utterly useless command for testing",
		Example: "Just run root echo",
		Annotations: map[string]string{
			"args":                "string to print",
			"string to printDesc": "A string to print",
		},
	}
	echoCmd.PersistentFlags().StringP("strone", "s", "one", "help message for flag strone")
	echoCmd.PersistentFlags().BoolP("persistentbool", "p", false, "help message for flag persistentbool")
	echoCmd.Flags().IntP("intone", "i", 123, "help message for flag intone")
	echoCmd.Flags().BoolP("boolone", "b", true, "help message for flag boolone")

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
	// We generate on a subcommand so we have both subcommands and parents
	buf := new(bytes.Buffer)
	Root() // init root
	if err := GenDocs(Echo(), buf); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	checkStringContains(t, output, Echo().Long)
	checkStringContains(t, output, Echo().Example)
	checkStringContains(t, output, "boolone")
	checkStringContains(t, output, "rootflag")
	checkStringContains(t, output, Echo().Annotations["string to printDesc"])
	checkStringOmits(t, output, Root().Short)
	checkStringContains(t, output, EchoSubCmd().Short)
	checkStringOmits(t, output, deprecatedCmd.Short)
}

func TestGenDocsNoHiddenParents(t *testing.T) {
	// We generate on a subcommand so we have both subcommands and parents
	for _, name := range []string{"rootflag", "strtwo"} {
		f := Root().PersistentFlags().Lookup(name)
		f.Hidden = true
		defer func() { f.Hidden = false }()
	}
	buf := new(bytes.Buffer)
	if err := GenDocs(Echo(), buf); err != nil {
		t.Fatal(err)
	}
	output := buf.String()

	checkStringContains(t, output, Echo().Long)
	checkStringContains(t, output, Echo().Example)
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
	c := &cobra.Command{Use: "do [OPTIONS] arg1 arg2"}

	tmpdir, err := ioutil.TempDir("", "test-gen-rst-tree")
	if err != nil {
		t.Fatalf("Failed to create tmpdir: %s", err.Error())
	}
	defer os.RemoveAll(tmpdir)

	if err := GenTreeDocs(c, tmpdir); err != nil {
		t.Fatalf("GenReSTTree failed: %s", err.Error())
	}

	if _, err := os.Stat(filepath.Join(tmpdir, "do.txt")); err != nil {
		t.Fatalf("Expected file 'do.rst' to exist")
	}
}

func BenchmarkGenDocsToFile(b *testing.B) {
	file, err := ioutil.TempFile("", "")
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
