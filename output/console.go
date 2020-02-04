package output

import (
	"fmt"
	"strings"

	"github.com/fabienbellanger/goCodeAnalyser/cloc"
)

const maxLanguagesLength = 20

// Console type.
type Console struct{}

// NewConsole return a pointer to a Console.
func NewConsole() *Console {
	return &Console{}
}

// Write displays result in the console.
func (c *Console) Write(result *cloc.Result, opts *cloc.Options) error {
	fmt.Printf("files=%v\nmax files length=%d\nlanguages=%v\ntotal=%v\n",
		result.Files,
		maxFilesLength(result.Files),
		result.Languages,
		result.Total)

	// Max length for title
	// --------------------
	maxTitle := maxLanguagesLength
	if opts.ByFile {
		maxTitle = maxFilesLength(result.Files)
	}

	header(opts.ByFile, maxTitle)
	footer(opts.ByFile, maxTitle)

	// fmt.Printf("<%-[1]*[2]v>\n", 10, "toto")
	// fmt.Printf("<%[1]*[2]v>\n", 10, 888)
	// fmt.Printf("<%-[1]*[2]v>\n", 10, 888)

	return nil
}

// maxFilesLength returns the max length of files.
func maxFilesLength(files map[string]*cloc.File) int {
	max := 0
	for k := range files {
		l := len(files[k].Name)
		if l > max {
			max = l
		}
	}
	return max
}

func header(byFile bool, maxLength int) {
	// 2*2 + 6*(3 + 9) + (maxLength + 4)
	title := "Language"
	if byFile {
		title = "File"
	}
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
	fmt.Printf("| %-[1]*[2]v | %-9v | %-9v | %-9v | %-9v | %-9v | %-9v |\n",
		maxLength+4, title, "Files", "Size", "Total", "Blanks", "Comments", "Code")
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
}

func footer(byFile bool, maxLength int) {
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
	fmt.Printf("| %-[1]*[2]v | %-9v | %-9v | %-9v | %-9v | %-9v | %-9v |\n",
		maxLength+4, "Total", "-", "-", "-", "-", "-", "-")
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
}
