package output

import "github.com/fabienbellanger/goCodeAnalyser/cloc"

// Writer is an interface for writting on console, JSON, CSV, etc.
type Writer interface {
	Write(*cloc.Result, *cloc.Options) error
}
