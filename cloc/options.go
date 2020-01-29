package cloc

import "regexp"

// Options lists CLOC application options.
type Options struct {
	Paths          []string
	Debug          bool
	SkipDuplicated bool
	ExcludeExts    map[string]struct{}
	IncludeLangs   map[string]struct{}
	MatchDir       *regexp.Regexp
	NotMatchDir    *regexp.Regexp
}

// NewOptions returns application options.
func NewOptions() *Options {
	return &Options{
		Paths:          make([]string, 0),
		Debug:          false,
		SkipDuplicated: false,
		ExcludeExts:    make(map[string]struct{}),
		IncludeLangs:   make(map[string]struct{}),
	}
}
