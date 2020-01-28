package cli

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// CmdOptions lists all command options.
type CmdOptions struct {
	ByFile bool
}

const (
	author  = "Fabien Bellanger"
	version = "0.0.1"
)

var (
	// color enables colors in console.
	color aurora.Aurora = aurora.NewAurora(true)

	// opts
	opts = CmdOptions{}

	rootCommand = &cobra.Command{
		Use:     "Go Code Analyser",
		Short:   "Go Code Analyser",
		Long:    "Go Code Analyser",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			tStart := time.Now()

			fmt.Printf("Args: %v\n", args)
			fmt.Printf("Options: %+v\n", opts)

			displayDuration(time.Since(tStart))
		},
	}
)

// Execute starts Cobra.
func Execute() error {
	// Version
	// -------
	// rootCommand.SetVersionTemplate("Vers,kdsfklds")

	// Flags
	// -----
	rootCommand.Flags().BoolVarP(&opts.ByFile, "by-file", "", false, "Display by file")

	// Launch root command
	// -------------------
	if err := rootCommand.Execute(); err != nil {
		return err
	}
	return nil
}

// displayDuration displays commands execution duration.
func displayDuration(d time.Duration) {
	fmt.Printf(color.Sprintf(color.Italic("\nCommand execution time: %v\n\n"), d))
}
