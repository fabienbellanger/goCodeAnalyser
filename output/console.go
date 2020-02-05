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
	body(opts.ByFile, opts.Sort, maxTitle, result)
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
func body(byFile bool, sortType string, maxLength int, r *cloc.Result) {
	if byFile {
		// Sort
		// ----
		filesSlice := make([]*cloc.File, 0, len(r.Files))
		for k := range r.Files {
			filesSlice = append(filesSlice, r.Files[k])
		}
		switch sortType {
		case "size":
			sort.Sort(cloc.FilesSort{Files: filesSlice, LessCmp: cloc.FileBySize})
		case "lines":
			sort.Sort(cloc.FilesSort{Files: filesSlice, LessCmp: cloc.FileByLines})
		case "comments":
			sort.Sort(cloc.FilesSort{Files: filesSlice, LessCmp: cloc.FileByComments})
		case "blanks":
			sort.Sort(cloc.FilesSort{Files: filesSlice, LessCmp: cloc.FileByBlanks})
		default:
			sort.Sort(cloc.FilesSort{Files: filesSlice, LessCmp: cloc.FileByCode})
		}

		for k := range filesSlice {
			fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
				maxLength+4,
				filesSlice[k].Name,
				"",
				goutils.HumanSizeWithPrecision(float64(filesSlice[k].Size), 0),
				filesSlice[k].Lines,
				filesSlice[k].Blanks,
				filesSlice[k].Comments,
				filesSlice[k].Code)
		}
	} else {
		// Sort by Code
		// ------------
		languagesSlice := make([]*cloc.Language, 0, len(r.Languages))
		for k := range r.Languages {
			languagesSlice = append(languagesSlice, r.Languages[k])
		}
		switch sortType {
		case "files":
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesByFiles})
		case "size":
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesBySize})
		case "lines":
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesByLines})
		case "comments":
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesByComments})
		case "blanks":
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesByBlanks})
		default:
			sort.Sort(cloc.LanguagesSort{Langs: languagesSlice, LessCmp: cloc.LanguagesByCode})
		}

		for k := range languagesSlice {
			fmt.Printf("| %-[1]*[2]v | %9v | %9v | %9v | %9v | %9v | %9v |\n",
				maxLength+4,
				languagesSlice[k].Name,
				languagesSlice[k].Total,
				goutils.HumanSizeWithPrecision(float64(languagesSlice[k].Size), 0),
				languagesSlice[k].Lines,
				languagesSlice[k].Blanks,
				languagesSlice[k].Comments,
				languagesSlice[k].Code)
		}
	}
}
