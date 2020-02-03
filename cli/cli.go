package cli

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fabienbellanger/goCodeAnalyser/cloc"
	"github.com/fabienbellanger/goutils"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// CmdOptions lists all command options.
type CmdOptions struct {
	ByFile         bool
	Debug          bool
	SkipDuplicated bool
	OutputType     string
	ExcludeExt     string
	IncludeLang    string
	MatchDir       string
	NotMatchDir    string
}

const (
	appName = "Go Code Analyser"
	// author  = "Fabien Bellanger"
	version = "0.1.0"
)

var (
	// color enables colors in console.
	color aurora.Aurora = aurora.NewAurora(true)

	// cmdOpts stores command options.
	cmdOpts = CmdOptions{}

	rootCommand = &cobra.Command{
		Use:     "goCodeAnalyser [paths]",
		Short:   "goCodeAnalyser [paths]",
		Long:    "goCodeAnalyser [paths]",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			tStart := time.Now()

			// Manage paths
			// ------------
			if len(args) == 0 {
				if err := cmd.Usage(); err != nil {
					goutils.CheckError(err, 1)
				}
				return
			}

			// List of all available languages
			// -------------------------------
			languages := cloc.NewDefinedLanguages()

			// Fill application options
			// ------------------------
			appOpts := fillOptions(cmdOpts, languages)

			fmt.Printf("Paths:       %+v\n", args)
			fmt.Printf("Cmd Options: %+v\n", cmdOpts)
			fmt.Printf("App Options: %+v\n", appOpts)
			fmt.Println("")

			// Launch process
			// --------------
			// TODO: To implement
			processor := cloc.NewProcessor(languages, appOpts, args)
			result, err := processor.Analyze()
			fmt.Printf("result=%v, err=%v\n", result, err)

			displayDuration(time.Since(tStart))
		},
	}
)

// Execute starts Cobra.
func Execute() error {
	// Version
	// -------
	rootCommand.SetVersionTemplate(appName + " version " + version + "\n")

	// Flags
	// -----
	rootCommand.Flags().BoolVar(&cmdOpts.ByFile, "by-file", false, "Display by file")
	rootCommand.Flags().BoolVar(&cmdOpts.Debug, "debug", false, "Display debug log")
	rootCommand.Flags().BoolVar(&cmdOpts.SkipDuplicated, "skip-duplicated", false, "Skip duplicated files")
	rootCommand.Flags().StringVar(&cmdOpts.OutputType, "output-type", "", "Output type [values: default,json,html]")
	rootCommand.Flags().StringVar(&cmdOpts.ExcludeExt, "exclude-ext", "", "Exclude file name extensions (separated commas)")
	rootCommand.Flags().StringVar(&cmdOpts.IncludeLang, "include-lang", "", "Include language name (separated commas)")
	rootCommand.Flags().StringVar(&cmdOpts.MatchDir, "match-dir", "", "Include dir name (regex)")
	rootCommand.Flags().StringVar(&cmdOpts.NotMatchDir, "not-match-dir", "", "Exclude dir name (regex)")

	// Launch root command
	// -------------------
	if err := rootCommand.Execute(); err != nil {
		return err
	}
	return nil
}

// fillOptions fills applications options from command options.
// TODO: Test
func fillOptions(cmdOpts CmdOptions, languages *cloc.DefinedLanguages) *cloc.Options {
	opts := cloc.NewOptions()
	opts.Debug = cmdOpts.Debug
	opts.SkipDuplicated = cmdOpts.SkipDuplicated

	// Excluded extensions
	// -------------------
	for _, ext := range strings.Split(cmdOpts.ExcludeExt, ",") {
		e, ok := cloc.Extensions[ext]
		if ok {
			opts.ExcludeExts[e] = struct{}{}
		}
	}

	// Match or not directory
	// ----------------------
	if cmdOpts.NotMatchDir != "" {
		opts.NotMatchDir = regexp.MustCompile(cmdOpts.NotMatchDir)
	}
	if cmdOpts.MatchDir != "" {
		opts.MatchDir = regexp.MustCompile(cmdOpts.MatchDir)
	}

	// Included languages
	// ------------------
	for _, lang := range strings.Split(cmdOpts.IncludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			opts.IncludeLangs[lang] = struct{}{}
		}
	}

	return opts
}

// displayDuration displays commands execution duration.
func displayDuration(d time.Duration) {
	fmt.Println(color.Sprintf(color.Italic("\nCommand execution time: %v\n"), d))
}
