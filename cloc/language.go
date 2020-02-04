package cloc

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"unicode"
)

// Language represents a language with its properties.
type Language struct {
	Name         string
	lineComments []string
	multiLines   [][]string
	Files        []string
	Code         int32
	Comments     int32
	Blanks       int32
	Total        int32
	Lines        int32
	Size         int64
}

// DefinedLanguages represents a map of available Language.
type DefinedLanguages struct {
	Langs map[string]*Language
}

// Languages is a slice of Language.
type Languages []Language

var (
	// shebang regex
	shebangEnvRegex  = regexp.MustCompile(`^#! *(\\S+/env) ([a-zA-Z]+)`)
	shebangLangRegex = regexp.MustCompile(`^#! *[.a-zA-Z/]+/([a-zA-Z]+)`)

	// shebangToExtension converts shebang to extension.
	shebangToExtension = map[string]string{
		"gosh":    "scm",
		"make":    "make",
		"perl":    "pl",
		"rc":      "plan9sh",
		"python":  "py",
		"ruby":    "rb",
		"escript": "erl",
	}
)

// NewLanguage returns a pointer to Language.
func NewLanguage(name string, lineComments []string, multiLines [][]string) *Language {
	return &Language{
		Name:         name,
		lineComments: lineComments,
		multiLines:   multiLines,
		Files:        []string{},
	}
}

// NewDefinedLanguages returns the list of all available languages with their properties.
func NewDefinedLanguages() *DefinedLanguages {
	return &DefinedLanguages{
		Langs: map[string]*Language{
			"ActionScript":        NewLanguage("ActionScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Ada":                 NewLanguage("Ada", []string{"--"}, [][]string{{"", ""}}),
			"Ant":                 NewLanguage("Ant", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"AsciiDoc":            NewLanguage("AsciiDoc", []string{}, [][]string{{"", ""}}),
			"Assembly":            NewLanguage("Assembly", []string{"//", ";", "#", "@", "|", "!"}, [][]string{{"/*", "*/"}}),
			"ATS":                 NewLanguage("ATS", []string{"//"}, [][]string{{"/*", "*/"}, {"(*", "*)"}}),
			"Awk":                 NewLanguage("Awk", []string{"#"}, [][]string{{"", ""}}),
			"Arduino Sketch":      NewLanguage("Arduino Sketch", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Batch":               NewLanguage("Batch", []string{"REM", "rem"}, [][]string{{"", ""}}),
			"BASH":                NewLanguage("BASH", []string{"#"}, [][]string{{"", ""}}),
			"BitBake":             NewLanguage("BitBake", []string{"#"}, [][]string{{"", ""}}),
			"C":                   NewLanguage("C", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C Header":            NewLanguage("C Header", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C Shell":             NewLanguage("C Shell", []string{"#"}, [][]string{{"", ""}}),
			"Cap'n Proto":         NewLanguage("Cap'n Proto", []string{"#"}, [][]string{{"", ""}}),
			"Carp":                NewLanguage("Carp", []string{";"}, [][]string{{"", ""}}),
			"C#":                  NewLanguage("C#", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Chapel":              NewLanguage("Chapel", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Clojure":             NewLanguage("Clojure", []string{"#", "#_"}, [][]string{{"", ""}}),
			"COBOL":               NewLanguage("COBOL", []string{"*", "/"}, [][]string{{"", ""}}),
			"CoffeeScript":        NewLanguage("CoffeeScript", []string{"#"}, [][]string{{"###", "###"}}),
			"Coq":                 NewLanguage("Coq", []string{"(*"}, [][]string{{"(*", "*)"}}),
			"ColdFusion":          NewLanguage("ColdFusion", []string{}, [][]string{{"<!---", "--->"}}),
			"ColdFusion CFScript": NewLanguage("ColdFusion CFScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"CMake":               NewLanguage("CMake", []string{"#"}, [][]string{{"", ""}}),
			"C++":                 NewLanguage("C++", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C++ Header":          NewLanguage("C++ Header", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Crystal":             NewLanguage("Crystal", []string{"#"}, [][]string{{"", ""}}),
			"CSS":                 NewLanguage("CSS", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Cython":              NewLanguage("Cython", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}}),
			"CUDA":                NewLanguage("CUDA", []string{"//"}, [][]string{{"/*", "*/"}}),
			"D":                   NewLanguage("D", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Dart":                NewLanguage("Dart", []string{"//", "///"}, [][]string{{"/*", "*/"}}),
			"Dhall":               NewLanguage("Dhall", []string{"--"}, [][]string{{"{-", "-}"}}),
			"DTrace":              NewLanguage("DTrace", []string{}, [][]string{{"/*", "*/"}}),
			"Device Tree":         NewLanguage("Device Tree", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Eiffel":              NewLanguage("Eiffel", []string{"--"}, [][]string{{"", ""}}),
			"Elm":                 NewLanguage("Elm", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Elixir":              NewLanguage("Elixir", []string{"#"}, [][]string{{"", ""}}),
			"Erlang":              NewLanguage("Erlang", []string{"%"}, [][]string{{"", ""}}),
			"Expect":              NewLanguage("Expect", []string{"#"}, [][]string{{"", ""}}),
			"Fish":                NewLanguage("Fish", []string{"#"}, [][]string{{"", ""}}),
			"Frege":               NewLanguage("Frege", []string{"--"}, [][]string{{"{-", "-}"}}),
			"F*":                  NewLanguage("F*", []string{"(*", "//"}, [][]string{{"(*", "*)"}}),
			"F#":                  NewLanguage("F#", []string{"(*"}, [][]string{{"(*", "*)"}}),
			"Lean":                NewLanguage("Lean", []string{"--"}, [][]string{{"/-", "-/"}}),
			"Logtalk":             NewLanguage("Logtalk", []string{"%"}, [][]string{{"", ""}}),
			"Lua":                 NewLanguage("Lua", []string{"--"}, [][]string{{"--[[", "]]"}}),
			"LISP":                NewLanguage("LISP", []string{";;"}, [][]string{{"#|", "|#"}}),
			"LiveScript":          NewLanguage("LiveScript", []string{"#"}, [][]string{{"/*", "*/"}}),
			"FORTRAN Legacy":      NewLanguage("FORTRAN Legacy", []string{"c", "C", "!", "*"}, [][]string{{"", ""}}),
			"FORTRAN Modern":      NewLanguage("FORTRAN Modern", []string{"!"}, [][]string{{"", ""}}),
			"Gherkin":             NewLanguage("Gherkin", []string{"#"}, [][]string{{"", ""}}),
			"GLSL":                NewLanguage("GLSL", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Go":                  NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Groovy":              NewLanguage("Groovy", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Haskell":             NewLanguage("Haskell", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Haxe":                NewLanguage("Haxe", []string{"//"}, [][]string{{"/*", "*/"}}),
			"HLSL":                NewLanguage("HLSL", []string{"//"}, [][]string{{"/*", "*/"}}),
			"HTML":                NewLanguage("HTML", []string{"//", "<!--"}, [][]string{{"<!--", "-->"}}),
			"Idris":               NewLanguage("Idris", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Io":                  NewLanguage("Io", []string{"//", "#"}, [][]string{{"/*", "*/"}}),
			"SKILL":               NewLanguage("SKILL", []string{";"}, [][]string{{"/*", "*/"}}),
			"JAI":                 NewLanguage("JAI", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Janet":               NewLanguage("Janet", []string{"#"}, [][]string{{"", ""}}),
			"Java":                NewLanguage("Java", []string{"//"}, [][]string{{"/*", "*/"}}),
			"JavaScript":          NewLanguage("JavaScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Julia":               NewLanguage("Julia", []string{"#"}, [][]string{{"#:=", ":=#"}}),
			"Jupyter Notebook":    NewLanguage("Jupyter Notebook", []string{"#"}, [][]string{{"", ""}}),
			"JSON":                NewLanguage("JSON", []string{}, [][]string{{"", ""}}),
			"JSX":                 NewLanguage("JSX", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Kotlin":              NewLanguage("Kotlin", []string{"//"}, [][]string{{"/*", "*/"}}),
			"LD Script":           NewLanguage("LD Script", []string{"//"}, [][]string{{"/*", "*/"}}),
			"LESS":                NewLanguage("LESS", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Objective-C":         NewLanguage("Objective-C", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Markdown":            NewLanguage("Markdown", []string{}, [][]string{{"", ""}}),
			"Nix":                 NewLanguage("Nix", []string{"#"}, [][]string{{"/*", "*/"}}),
			"NSIS":                NewLanguage("NSIS", []string{"#", ";"}, [][]string{{"/*", "*/"}}),
			"Nu":                  NewLanguage("Nu", []string{";", "#"}, [][]string{{"", ""}}),
			"OCaml":               NewLanguage("OCaml", []string{}, [][]string{{"(*", "*)"}}),
			"Objective-C++":       NewLanguage("Objective-C++", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Makefile":            NewLanguage("Makefile", []string{"#"}, [][]string{{"", ""}}),
			"MATLAB":              NewLanguage("MATLAB", []string{"%"}, [][]string{{"%{", "}%"}}),
			"Mercury":             NewLanguage("Mercury", []string{"%"}, [][]string{{"/*", "*/"}}),
			"Maven":               NewLanguage("Maven", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Meson":               NewLanguage("Meson", []string{"#"}, [][]string{{"", ""}}),
			"Mustache":            NewLanguage("Mustache", []string{}, [][]string{{"{{!", "}}"}}),
			"M4":                  NewLanguage("M4", []string{"#"}, [][]string{{"", ""}}),
			"Nim":                 NewLanguage("Nim", []string{"#"}, [][]string{{"#[", "]#"}}),
			"lex":                 NewLanguage("lex", []string{}, [][]string{{"/*", "*/"}}),
			"PHP":                 NewLanguage("PHP", []string{"#", "//"}, [][]string{{"/*", "*/"}}),
			"Pascal":              NewLanguage("Pascal", []string{"//"}, [][]string{{"{", ")"}}),
			"Perl":                NewLanguage("Perl", []string{"#"}, [][]string{{":=", ":=cut"}}),
			"Plain Text":          NewLanguage("Plain Text", []string{}, [][]string{{"", ""}}),
			"Plan9 Shell":         NewLanguage("Plan9 Shell", []string{"#"}, [][]string{{"", ""}}),
			"Pony":                NewLanguage("Pony", []string{"//"}, [][]string{{"/*", "*/"}}),
			"PowerShell":          NewLanguage("PowerShell", []string{"#"}, [][]string{{"<#", "#>"}}),
			"Polly":               NewLanguage("Polly", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Protocol Buffers":    NewLanguage("Protocol Buffers", []string{"//"}, [][]string{{"", ""}}),
			"Python":              NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}}),
			"Q":                   NewLanguage("Q", []string{"/ "}, [][]string{{"\\", "/"}, {"/", "\\"}}),
			"QML":                 NewLanguage("QML", []string{"//"}, [][]string{{"/*", "*/"}}),
			"R":                   NewLanguage("R", []string{"#"}, [][]string{{"", ""}}),
			"Rebol":               NewLanguage("Rebol", []string{";"}, [][]string{{"", ""}}),
			"Red":                 NewLanguage("Red", []string{";"}, [][]string{{"", ""}}),
			"RMarkdown":           NewLanguage("RMarkdown", []string{}, [][]string{{"", ""}}),
			"RAML":                NewLanguage("RAML", []string{"#"}, [][]string{{"", ""}}),
			"Racket":              NewLanguage("Racket", []string{";"}, [][]string{{"#|", "|#"}}),
			"ReStructuredText":    NewLanguage("ReStructuredText", []string{}, [][]string{{"", ""}}),
			"Ruby":                NewLanguage("Ruby", []string{"#"}, [][]string{{":=begin", ":=end"}}),
			"Ruby HTML":           NewLanguage("Ruby HTML", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Rust":                NewLanguage("Rust", []string{"//", "///", "//!"}, [][]string{{"/*", "*/"}}),
			"Scala":               NewLanguage("Scala", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Sass":                NewLanguage("Sass", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Scheme":              NewLanguage("Scheme", []string{";"}, [][]string{{"#|", "|#"}}),
			"sed":                 NewLanguage("sed", []string{"#"}, [][]string{{"", ""}}),
			"Stan":                NewLanguage("Stan", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Solidity":            NewLanguage("Solidity", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Bourne Shell":        NewLanguage("Bourne Shell", []string{"#"}, [][]string{{"", ""}}),
			"Standard ML":         NewLanguage("Standard ML", []string{}, [][]string{{"(*", "*)"}}),
			"SQL":                 NewLanguage("SQL", []string{"--"}, [][]string{{"/*", "*/"}}),
			"Swift":               NewLanguage("Swift", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Terra":               NewLanguage("Terra", []string{"--"}, [][]string{{"--[[", "]]"}}),
			"TeX":                 NewLanguage("TeX", []string{"%"}, [][]string{{"", ""}}),
			"Isabelle":            NewLanguage("Isabelle", []string{}, [][]string{{"(*", "*)"}}),
			"TLA":                 NewLanguage("TLA", []string{"/*"}, [][]string{{"(*", "*)"}}),
			"Tcl/Tk":              NewLanguage("Tcl/Tk", []string{"#"}, [][]string{{"", ""}}),
			"TOML":                NewLanguage("TOML", []string{"#"}, [][]string{{"", ""}}),
			"TypeScript":          NewLanguage("TypeScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"HCL":                 NewLanguage("HCL", []string{"#", "//"}, [][]string{{"/*", "*/"}}),
			"Unity-Prefab":        NewLanguage("Unity-Prefab", []string{}, [][]string{{"", ""}}),
			"MSBuild script":      NewLanguage("MSBuild script", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Vala":                NewLanguage("Vala", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Verilog":             NewLanguage("Verilog", []string{"//"}, [][]string{{"/*", "*/"}}),
			"VimL":                NewLanguage("VimL", []string{`"`}, [][]string{{"", ""}}),
			"Vue":                 NewLanguage("Vue", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"WiX":                 NewLanguage("WiX", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XML":                 NewLanguage("XML", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XSLT":                NewLanguage("XSLT", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XSD":                 NewLanguage("XSD", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"YAML":                NewLanguage("YAML", []string{"#"}, [][]string{{"", ""}}),
			"Yacc":                NewLanguage("Yacc", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Zephir":              NewLanguage("Zephir", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Zig":                 NewLanguage("Zig", []string{"//", "///"}, [][]string{{"", ""}}),
			"Zsh":                 NewLanguage("Zsh", []string{"#"}, [][]string{{"", ""}}),
		},
	}
}

// getShebang returns shebang.
func getShebang(line string) (shebangLang string, ok bool) {
	ret := shebangEnvRegex.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) == 3 {
		shebangLang = ret[0][2]
		if sl, ok := shebangToExtension[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	ret = shebangLangRegex.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) >= 2 {
		shebangLang = ret[0][1]
		if sl, ok := shebangToExtension[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	return "", false
}

// getExtensionByShebang returns extension from shebang.
func getExtensionByShebang(path string) (shebangLang string, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		return shebangLang, false
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return shebangLang, ok
	}
	line = bytes.TrimLeftFunc(line, unicode.IsSpace)

	if len(line) > 2 && line[0] == '#' && line[1] == '!' {
		return getShebang(string(line))
	}
	return shebangLang, ok
}

// getTotalFiles returns the number of files to analyze.
func getTotalFiles(langs map[string]*Language) (t int) {
	for _, lang := range langs {
		t += len(lang.Files)
	}
	return t
}

// isLanguageAnalysable checks if a language must be analyze
// (VCS, match and not-match directory).
// The function returns true if it can be analyzed.
func isLanguageAnalysable(path string, vcsInRoot bool, opts *Options) bool {
	// Check VCS
	// ---------
	if !vcsInRoot && isVCSDir(path) {
		return false
	}

	// Check match and not-match directory options
	// -------------------------------------------
	dir := filepath.Dir(path)
	if opts.NotMatchDir != nil && opts.NotMatchDir.MatchString(dir) {
		return false
	}
	if opts.MatchDir != nil && !opts.MatchDir.MatchString(dir) {
		return false
	}
	return true
}
