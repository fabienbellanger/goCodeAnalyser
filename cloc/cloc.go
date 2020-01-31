package cloc

import (
	"fmt"
	"os"
	"path/filepath"
)

// Processor represents a process instance
type Processor struct {
	langs *DefinedLanguages
	opts  *Options
	paths []string
}

// Result returns the analysis results
type Result struct {
	// Total         *Language
	Files     map[string]*File
	Languages map[string]*Language
}

// NewProcessor returns a processor.
func NewProcessor(langs *DefinedLanguages, options *Options, paths []string) *Processor {
	return &Processor{
		langs: langs,
		opts:  options,
		paths: paths,
	}
}

// Analyze starts files analysis.
func (p *Processor) Analyze() (*Result, error) {
	// List all files and init languages
	// ---------------------------------
	languages, err := p.initLanguages()
	if err != nil {
		return nil, err
	}

	// Analyze of each file
	// --------------------
	files := make(map[string]*File, getTotalFiles(languages))
	fmt.Printf("files=%+v\n", files)
	for _, language := range languages {
		for _, file := range language.Files {
			// File analysis
			// -------------
			f := NewFile(file, language.Name)
			f.analyze(language, p.opts)

			// Update language
			// ---------------
			language.Size += f.Size

			fmt.Printf("file=%+v\n", f)
		}
	}
	fmt.Printf("\nlanguages=%v\n", languages["Go"])

	return nil, nil
}

// initLanguages lists all files form paths and inits languages.
func (p *Processor) initLanguages() (result map[string]*Language, err error) {
	result = make(map[string]*Language, 0)
	filesCache := make(map[string]struct{})

	for _, root := range p.paths {
		vcsInRoot := isVCSDir(root)
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			// Check VCS
			// ---------
			if !vcsInRoot && isVCSDir(path) {
				return nil
			}

			// Check match and not-match directory options
			// -------------------------------------------
			// TODO: make a function to reduce cycomatic complexity
			dir := filepath.Dir(path)
			if p.opts.NotMatchDir != nil && p.opts.NotMatchDir.MatchString(dir) {
				return nil
			}
			if p.opts.MatchDir != nil && !p.opts.MatchDir.MatchString(dir) {
				return nil
			}

			// Check file extension
			// --------------------
			if ext, ok := getExtension(path, p.opts); ok {
				// Get Language
				// ------------
				if lang, ok := Extensions[ext]; ok {
					// Check Options
					// -------------
					if ok := checkFileOptions(path, lang, p.opts, filesCache); ok {
						// Add to languages list
						// ---------------------
						if _, ok := result[lang]; !ok {
							result[lang] = NewLanguage(
								p.langs.Langs[lang].Name,
								p.langs.Langs[lang].lineComments,
								p.langs.Langs[lang].multiLines)
						}
						result[lang].Files = append(result[lang].Files, path)
					}
				}
			}

			return nil
		})
	}

	return result, err
}
