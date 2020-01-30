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

// Result returns the analyse results
type Result struct {
	// Total         *Language
	// Files         map[string]*ClocFile
	Languages     map[string]*Language
	MaxPathLength int
}

// NewProcessor returns a processor.
func NewProcessor(langs *DefinedLanguages, options *Options, paths []string) *Processor {
	return &Processor{
		langs: langs,
		opts:  options,
		paths: paths,
	}
}

// Analyse starts files analyse.
func (p *Processor) Analyse() (*Result, error) {
	// List all files and init languages
	// ---------------------------------
	languages, err := p.initLanguages()
	if err != nil {
		return nil, err
	}
	fmt.Printf("languages=%v\n", languages)

	return nil, nil
}

// initLanguages lists all files form paths and inits languages.
func (p *Processor) initLanguages() (result map[string]*Language, err error) {
	result = make(map[string]*Language, 0)
	// filesCache := make(map[string]struct{})

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
			dir := filepath.Dir(path)
			if p.opts.NotMatchDir != nil && p.opts.NotMatchDir.MatchString(dir) {
				return nil
			}
			if p.opts.MatchDir != nil && !p.opts.MatchDir.MatchString(dir) {
				return nil
			}

			// Check file extension
			// --------------------
			// TODO: To do!
			if ext, ok := getExtension(path, p.opts); ok {
				if lang, ok := Extensions[ext]; ok {
					fmt.Printf("ext=%v,\tlang=%v\n", ext, lang)

					// TODO: Check options

					// TODO: Fill result
				}
			}

			return nil
		})
	}

	return result, err
}
