package cloc

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/src-d/enry/v2"
)

var (
	// Extensions lists all available extensions.
	Extensions = map[string]string{
		"as":          "ActionScript",
		"ada":         "Ada",
		"adb":         "Ada",
		"ads":         "Ada",
		"Ant":         "Ant",
		"adoc":        "AsciiDoc",
		"asciidoc":    "AsciiDoc",
		"asm":         "Assembly",
		"S":           "Assembly",
		"s":           "Assembly",
		"dats":        "ATS",
		"sats":        "ATS",
		"hats":        "ATS",
		"awk":         "Awk",
		"bat":         "Batch",
		"btm":         "Batch",
		"bb":          "BitBake",
		"cbl":         "COBOL",
		"cmd":         "Batch",
		"bash":        "BASH",
		"sh":          "Bourne Shell",
		"c":           "C",
		"carp":        "Carp",
		"csh":         "C Shell",
		"ec":          "C",
		"erl":         "Erlang",
		"hrl":         "Erlang",
		"pgc":         "C",
		"capnp":       "Cap'n Proto",
		"chpl":        "Chapel",
		"cs":          "C#",
		"clj":         "Clojure",
		"coffee":      "CoffeeScript",
		"cfm":         "ColdFusion",
		"cfc":         "ColdFusion CFScript",
		"cmake":       "CMake",
		"cc":          "C++",
		"cpp":         "C++",
		"cxx":         "C++",
		"pcc":         "C++",
		"c++":         "C++",
		"cr":          "Crystal",
		"css":         "CSS",
		"cu":          "CUDA",
		"d":           "D",
		"dart":        "Dart",
		"dhall":       "Dhall",
		"dtrace":      "DTrace",
		"dts":         "Device Tree",
		"dtsi":        "Device Tree",
		"e":           "Eiffel",
		"elm":         "Elm",
		"el":          "LISP",
		"exp":         "Expect",
		"ex":          "Elixir",
		"exs":         "Elixir",
		"feature":     "Gherkin",
		"fish":        "Fish",
		"fr":          "Frege",
		"fst":         "F*",
		"F#":          "F#",   // deplicated F#/GLSL
		"GLSL":        "GLSL", // both use ext '.fs'
		"vs":          "GLSL",
		"shader":      "HLSL",
		"cg":          "HLSL",
		"cginc":       "HLSL",
		"hlsl":        "HLSL",
		"lean":        "Lean",
		"hlean":       "Lean",
		"lgt":         "Logtalk",
		"lisp":        "LISP",
		"lsp":         "LISP",
		"lua":         "Lua",
		"ls":          "LiveScript",
		"sc":          "LISP",
		"f":           "FORTRAN Legacy",
		"f77":         "FORTRAN Legacy",
		"for":         "FORTRAN Legacy",
		"ftn":         "FORTRAN Legacy",
		"pfo":         "FORTRAN Legacy",
		"f90":         "FORTRAN Modern",
		"f95":         "FORTRAN Modern",
		"f03":         "FORTRAN Modern",
		"f08":         "FORTRAN Modern",
		"go":          "Go",
		"groovy":      "Groovy",
		"gradle":      "Groovy",
		"h":           "C Header",
		"hs":          "Haskell",
		"hpp":         "C++ Header",
		"hh":          "C++ Header",
		"html":        "HTML",
		"hx":          "Haxe",
		"hxx":         "C++ Header",
		"idr":         "Idris",
		"il":          "SKILL",
		"ino":         "Arduino Sketch",
		"io":          "Io",
		"ipynb":       "Jupyter Notebook",
		"jai":         "JAI",
		"java":        "Java",
		"js":          "JavaScript",
		"jl":          "Julia",
		"janet":       "Janet",
		"json":        "JSON",
		"jsx":         "JSX",
		"kt":          "Kotlin",
		"lds":         "LD Script",
		"less":        "LESS",
		"Objective-C": "Objective-C", // deplicated Obj-C/Matlab/Mercury
		"Matlab":      "MATLAB",      // both use ext '.m'
		"Mercury":     "Mercury",     // use ext '.m'
		"md":          "Markdown",
		"markdown":    "Markdown",
		"nix":         "Nix",
		"nsi":         "NSIS",
		"nsh":         "NSIS",
		"nu":          "Nu",
		"ML":          "OCaml",
		"ml":          "OCaml",
		"mli":         "OCaml",
		"mll":         "OCaml",
		"mly":         "OCaml",
		"mm":          "Objective-C++",
		"maven":       "Maven",
		"makefile":    "Makefile",
		"meson":       "Meson",
		"mustache":    "Mustache",
		"m4":          "M4",
		"l":           "lex",
		"nim":         "Nim",
		"php":         "PHP",
		"pas":         "Pascal",
		"PL":          "Perl",
		"pl":          "Perl",
		"pm":          "Perl",
		"plan9sh":     "Plan9 Shell",
		"pony":        "Pony",
		"ps1":         "PowerShell",
		"text":        "Plain Text",
		"txt":         "Plain Text",
		"polly":       "Polly",
		"proto":       "Protocol Buffers",
		"py":          "Python",
		"pxd":         "Cython",
		"pyx":         "Cython",
		"q":           "Q",
		"qml":         "QML",
		"r":           "R",
		"R":           "R",
		"raml":        "RAML",
		"Rebol":       "Rebol",
		"red":         "Red",
		"Rmd":         "RMarkdown",
		"rake":        "Ruby",
		"rb":          "Ruby",
		"rkt":         "Racket",
		"rhtml":       "Ruby HTML",
		"rs":          "Rust",
		"rst":         "ReStructuredText",
		"sass":        "Sass",
		"scala":       "Scala",
		"scss":        "Sass",
		"scm":         "Scheme",
		"sed":         "sed",
		"stan":        "Stan",
		"sml":         "Standard ML",
		"sol":         "Solidity",
		"sql":         "SQL",
		"swift":       "Swift",
		"t":           "Terra",
		"tex":         "TeX",
		"thy":         "Isabelle",
		"tla":         "TLA",
		"sty":         "TeX",
		"tcl":         "Tcl/Tk",
		"toml":        "TOML",
		"TypeScript":  "TypeScript",
		"tsx":         "TypeScript",
		"tf":          "HCL",
		"mat":         "Unity-Prefab",
		"prefab":      "Unity-Prefab",
		"Coq":         "Coq",
		"vala":        "Vala",
		"Verilog":     "Verilog",
		"csproj":      "MSBuild script",
		"vcproj":      "MSBuild script",
		"vim":         "VimL",
		"vue":         "Vue",
		"xml":         "XML",
		"XML":         "XML",
		"xsd":         "XSD",
		"xsl":         "XSLT",
		"xslt":        "XSLT",
		"wxs":         "WiX",
		"yaml":        "YAML",
		"yml":         "YAML",
		"y":           "Yacc",
		"zep":         "Zephir",
		"zig":         "Zig",
		"zsh":         "Zsh",
	}
)

// getExtension returns file extension from a path.
func getExtension(path string, opts *Options) (ext string, ok bool) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)

	switch ext {
	case ".m", ".v", ".fs", ".r", ".ts":
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return "", false
		}

		lang := enry.GetLanguage(path, content)
		if opts.Debug {
			fmt.Printf("path=%v, lang=%v\n", path, lang)
		}
		return lang, true
	}

	switch base {
	case "meson.build", "meson_options.txt":
		return "meson", true
	case "CMakeLists.txt":
		return "cmake", true
	case "configure.ac":
		return "m4", true
	case "Makefile.am":
		return "makefile", true
	case "build.xml":
		return "Ant", true
	case "pom.xml":
		return "maven", true
	}

	switch strings.ToLower(base) {
	case "makefile":
		return "makefile", true
	case "nukefile":
		return "nu", true
	case "rebar": // skip
		return "", false
	}

	shebangLang, ok := getExtensionByShebang(path)
	if ok {
		return shebangLang, true
	}

	if len(ext) >= 2 {
		return ext[1:], true
	}

	return ext, ok
}
