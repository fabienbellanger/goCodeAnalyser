package cloc

import "github.com/fabienbellanger/goCodeAnalyser/language"

// Processor represents a process instance
type Processor struct {
	langs *language.DefinedLanguages
	opts  *Options
}

// Result returns the analyse results
type Result struct {
	Total         *language.Language
	// Files         map[string]*ClocFile
	Languages     map[string]*language.Language
	MaxPathLength int
}

// NewProcessor returns a processor.
func NewProcessor(langs *language.DefinedLanguages, options *Options) *Processor {
	return &Processor{
		langs: langs,
		opts:  options,
	}
}

// Analyse starts files analyse.
func (p *Processor) Analyse() (*Result, error) {
	return nil, nil
}
