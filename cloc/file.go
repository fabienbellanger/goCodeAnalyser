package cloc

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// File represents a file with its properties.
type File struct {
	Name     string `xml:"name,attr" json:"name"`
	Language string `xml:"language,attr" json:"language"`
	Code     int32  `xml:"code,attr" json:"code"`
	Comments int32  `xml:"comment,attr" json:"comment"`
	Blanks   int32  `xml:"blank,attr" json:"blank"`
}

// Files is a slice of File.
type Files []File

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
