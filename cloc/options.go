package cloc

import (
	"fmt"
	"regexp"

	"github.com/fabienbellanger/goutils"
)

// Options lists CLOC application options.
type Options struct {
	ByFile         bool
	Debug          bool
	SkipDuplicated bool
	ExcludeExts    map[string]struct{}
	IncludeLangs   map[string]struct{}
	MatchDir       *regexp.Regexp
	NotMatchDir    *regexp.Regexp
	Sort           string
}

// NewOptions returns application options.
func NewOptions() *Options {
	return &Options{
		ByFile:         false,
		Debug:          false,
		SkipDuplicated: false,
		ExcludeExts:    make(map[string]struct{}),
		IncludeLangs:   make(map[string]struct{}),
		Sort:           "code",
	}
}

// CheckSort checks if sort value is correct (code, size, lines, comments, blanks or files).
func CheckSort(s string) bool {
	correctSorts := []string{"code", "size", "lines", "comments", "blanks", "files"}
	return goutils.StringInSlice(s, correctSorts)
}

// checkFileOptions checks if a file respects options.
func checkFileOptions(path, lang string, opts *Options, filesCache map[string]struct{}) bool {
	if _, ok := opts.ExcludeExts[lang]; ok {
		return false
	}

	if len(opts.IncludeLangs) != 0 {
		if _, ok := opts.IncludeLangs[lang]; !ok {
			return false
		}
	}

	if !opts.SkipDuplicated {
		ignore := checkMD5Sum(path, filesCache)
		if ignore {
			if opts.Debug {
				fmt.Printf("[ignore=%v] find same md5\n", path)
			}
			return false
		}
	}

	return true
}
