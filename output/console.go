package output

import (
	"fmt"
	"github.com/fabienbellanger/goCodeAnalyser/cloc"
)

// Console type.
type Console struct{}

// NewConsole return a pointer to a Console.
func NewConsole() *Console {
	return &Console{}
}

// Write displays result in the console.
func (c *Console) Write(result *cloc.Result) error {
	fmt.Println("Display result in console...")
	fmt.Printf("files=%v\nlanguages=%v\n", result.Files, result.Languages)
	return nil
}
