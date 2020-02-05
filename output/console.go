package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fabienbellanger/goCodeAnalyser/cloc"
	"github.com/fabienbellanger/goutils"
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
	// Max length for title
	// --------------------
	maxTitle := maxLanguagesLength
	if opts.ByFile {
		maxTitle = maxFilesLength(result.Files)
	}

	// Display results
	// ---------------
	header(opts.ByFile, maxTitle)
	body(opts.ByFile, maxTitle, result)
	footer(opts.ByFile, maxTitle, result.Total)

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

// header displays array header.
func header(byFile bool, maxLength int) {
	// 2*2 + 6*(3 + 9) + (maxLength + 4)
	title := "Language"
	if byFile {
		title = "File"
	}
	fmt.Printf("\n%v\n", strings.Repeat("-", 80+maxLength))
	fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
		maxLength+4, title, "Files", "Size", "Lines", "Blanks", "Comments", "Code")
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
}

// footer displays array footer.
func footer(byFile bool, maxLength int, t *cloc.Language) {
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
	fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
		maxLength+4, "Total", t.Total, goutils.HumanSizeWithPrecision(float64(t.Size), 0), t.Lines, t.Blanks, t.Comments, t.Code)
	fmt.Printf("%v\n", strings.Repeat("-", 80+maxLength))
}

// body displays languages or files information.
func body(byFile bool, maxLength int, r *cloc.Result) {
	if byFile {
		// Sort by Code
		// ------------
		sortedFiles := make(cloc.FilesByCode, 0, len(r.Files))
		for k := range r.Files {
			sortedFiles = append(sortedFiles, *r.Files[k])
		}
		sort.Sort(sortedFiles)

		for k := range sortedFiles {
			fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
				maxLength+4,
				sortedFiles[k].Name,
				"",
				goutils.HumanSizeWithPrecision(float64(sortedFiles[k].Size), 0),
				sortedFiles[k].Lines,
				sortedFiles[k].Blanks,
				sortedFiles[k].Comments,
				sortedFiles[k].Code)
		}
	} else {
		// Sort by Code
		// ------------
		sortedLanguages := make(cloc.LanguagesByCode, 0, len(r.Languages))
		for k := range r.Languages {
			sortedLanguages = append(sortedLanguages, *r.Languages[k])
		}
		sort.Sort(sortedLanguages)

		for k := range sortedLanguages {
			fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
				maxLength+4,
				sortedLanguages[k].Name,
				sortedLanguages[k].Total,
				goutils.HumanSizeWithPrecision(float64(sortedLanguages[k].Size), 0),
				sortedLanguages[k].Lines,
				sortedLanguages[k].Blanks,
				sortedLanguages[k].Comments,
				sortedLanguages[k].Code)
		}
	}
}
