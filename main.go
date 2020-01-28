package main

import (
	"os"

	"github.com/fabienbellanger/goCodeAnalyser/cli"
)

func main() {
	// Lancement du CLI
	// ----------------
	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
