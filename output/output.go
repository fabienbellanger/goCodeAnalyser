package output

import "github.com/fabienbellanger/goCodeAnalyser/cloc"

// Writer is an interface for writting on console, JSON, etc.
type Writer interface {
	Write(*cloc.Result) error
}
