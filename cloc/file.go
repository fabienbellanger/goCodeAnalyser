package cloc

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unicode"
)

// File represents a file with its properties.
type File struct {
	Name     string `xml:"name,attr" json:"name"`
	Size     int64  `xml:"size,attr" json:"size"`
	Language string `xml:"language,attr" json:"language"`
	Code     int32  `xml:"code,attr" json:"code"`
	Comments int32  `xml:"comment,attr" json:"comment"`
	Blanks   int32  `xml:"blank,attr" json:"blank"`
	Lines    int32  `xml:"lines,attr" json:"lines"`
}

var bsPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 128*1024)
		return &b
	},
}

// getByteSlice returns an array of bytes.
func getByteSlice() *[]byte {
	return bsPool.Get().(*[]byte)
}

// putByteSlice puts an array of bytes in the pool.
func putByteSlice(bs *[]byte) {
	bsPool.Put(bs)
}

// NewFile returns a pointer to File.
func NewFile(name, language string) *File {
	return &File{
		Name:     name,
		Language: language,
	}
}

// analyze analyze a file.
func (f *File) analyze(language *Language, opts *Options) {
	// Open file
	// ---------
	file, err := os.Open(f.Name)
	if err != nil {
		return
	}
	defer file.Close()

	// File size
	// ---------
	fileInfo, _ := file.Stat()
	if err == nil {
		f.Size = fileInfo.Size()
	}

	// Debug mode
	// ----------
	if opts.Debug {
		fmt.Printf("\n> %s\n%s\n", f.Name, strings.Repeat("-", len(f.Name)+2))
	}

	// File analysis
	// -------------
	f.read(file, language, opts)
}

// read reads file to analyze.
func (f *File) read(file *os.File, language *Language, opts *Options) {
	// Buffer creation
	// ---------------
	buf := getByteSlice()
	defer putByteSlice(buf)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(*buf, 1024*1024)

	isFirstLine := true
	inComments := [][2]string{}

	// Lines
	// -----
	for scanner.Scan() {
		f.Lines++
		lineOrg := scanner.Text()
		line := strings.TrimSpace(lineOrg)

		if len(strings.TrimSpace(line)) == 0 {
			f.onBlank(opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		// shebang line is 'code'
		// ----------------------
		if isFirstLine && strings.HasPrefix(line, "#!") {
			f.onCode(opts, len(inComments) > 0, line, lineOrg)
			isFirstLine = false
			continue
		}

		// Single comments
		// ---------------
		if len(inComments) == 0 {
			if isFirstLine {
				line = trimBOM(line)
			}

			stop := false
			inComments, stop = f.analyzeSingleComments(line, lineOrg, language, opts, inComments)
			if stop {
				continue
			}
		}

		if len(inComments) == 0 && !containsComment(line, language.multiLines) {
			f.onCode(opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		// Multi comments
		// --------------
		if len(language.multiLines) == 1 && len(language.multiLines[0]) == 2 && language.multiLines[0][0] == "" {
			f.onCode(opts, len(inComments) > 0, line, lineOrg)
			continue
		}
		inComments = f.analyzeMultiComments(line, lineOrg, language, opts, inComments)
	}
}

// analyzeSingleComments analyzes single comments in file line.
// It returns inComments array and a bool. If it equals to true, the program goes to the next line.
func (f *File) analyzeSingleComments(line, lineOrg string, lang *Language, opts *Options, inComments [][2]string) ([][2]string, bool) {
singleloop:
	for _, singleComment := range lang.lineComments {
		if strings.HasPrefix(line, singleComment) {
			// Check if single comment is a prefix of multi comment
			for _, ml := range lang.multiLines {
				if ml[0] != "" && strings.HasPrefix(line, ml[0]) {
					break singleloop
				}
			}
			f.onComment(opts, len(inComments) > 0, line, lineOrg)
			return inComments, true
		}
	}

	if len(lang.multiLines) == 0 {
		f.onCode(opts, len(inComments) > 0, line, lineOrg)
		return inComments, true
	}

	return inComments, false
}

// analyzeMultiComments analyzes multi comments in file line and return inComments array with updates.
func (f *File) analyzeMultiComments(line, lineOrg string, lang *Language, opts *Options, inComments [][2]string) [][2]string {
	isCode := false
	lenLine := len(line)
	for pos := 0; pos < lenLine; pos++ {
		for _, ml := range lang.multiLines {
			begin, end := ml[0], ml[1]
			lenBegin := len(begin)

			if pos+lenBegin <= lenLine && strings.HasPrefix(line[pos:], begin) && (begin != end || len(inComments) == 0) {
				pos += lenBegin
				inComments = append(inComments, [2]string{begin, end})
				continue
			}

			if n := len(inComments); n > 0 {
				last := inComments[n-1]
				if pos+len(last[1]) <= lenLine && strings.HasPrefix(line[pos:], last[1]) {
					inComments = inComments[:n-1]
					pos += len(last[1])
				}
			} else if pos < lenLine && !unicode.IsSpace(nextRune(line[pos:])) {
				isCode = true
			}
		}
	}

	if isCode {
		f.onCode(opts, len(inComments) > 0, line, lineOrg)
	} else {
		f.onComment(opts, len(inComments) > 0, line, lineOrg)
	}

	return inComments
}

// onBlank update File blanks informations.
func (f *File) onBlank(opts *Options, isInComments bool, line, lineOrg string) {
	f.Blanks++
	if opts.Debug {
		fmt.Printf("[BLNK, cd:%d, cm:%d, bk:%d, iscm:%v] %s\n",
			f.Code, f.Comments, f.Blanks, isInComments, lineOrg)
	}
}

// onComment update File comments informations.
func (f *File) onComment(opts *Options, isInComments bool, line, lineOrg string) {
	f.Comments++
	if opts.Debug {
		fmt.Printf("[COMM, cd:%d, cm:%d, bk:%d, iscm:%v] %s\n",
			f.Code, f.Comments, f.Blanks, isInComments, lineOrg)
	}
}

// onCode update File code informations.
func (f *File) onCode(opts *Options, isInComments bool, line, lineOrg string) {
	f.Code++
	if opts.Debug {
		fmt.Printf("[CODE, cd:%d, cm:%d, bk:%d, iscm:%v] %s\n",
			f.Code, f.Comments, f.Blanks, isInComments, lineOrg)
	}
}

// isVCSDir checks if directory is a version control system.
func isVCSDir(path string) bool {
	if len(path) > 1 && path[0] == os.PathSeparator {
		path = path[1:]
	}
	vcsDirs := []string{".bzr", ".cvs", ".hg", ".git", ".svn"}
	for _, dir := range vcsDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

// trimBOM trims BOM (UTF-8) of a line.
// TODO: In goutils
func trimBOM(line string) string {
	l := len(line)
	if l >= 3 {
		if line[0] == 0xef && line[1] == 0xbb && line[2] == 0xbf {
			trimLine := line[3:]
			return trimLine
		}
	}
	return line
}

// checkMD5Sum checks md5sum for a path and returns true if a file file
// has ready been added.
func checkMD5Sum(path string, fileCache map[string]struct{}) (ignore bool) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return true
	}

	// Calculate md5sum
	// ----------------
	hash := md5.Sum(content)
	c := fmt.Sprintf("%x", hash)
	if _, ok := fileCache[c]; ok {
		return true
	}

	fileCache[c] = struct{}{}

	return false
}

// containsComment checks if a line contains a comment.
func containsComment(line string, multiLines [][]string) bool {
	for _, ml := range multiLines {
		for _, c := range ml {
			if strings.Contains(line, c) {
				return true
			}
		}
	}
	return false
}

// nextRune returns the next rune.
// TODO: In goutils
func nextRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}
