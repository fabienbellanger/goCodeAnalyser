package cloc

// Languages sorts
// ===============

// LanguagesSort is used to sort languages.
type LanguagesSort struct {
	Langs   []*Language
	LessCmp func(i, j *Language) bool
}

func (l LanguagesSort) Len() int           { return len(l.Langs) }
func (l LanguagesSort) Swap(i, j int)      { l.Langs[i], l.Langs[j] = l.Langs[j], l.Langs[i] }
func (l LanguagesSort) Less(i, j int) bool { return l.LessCmp(l.Langs[i], l.Langs[j]) }

// LanguagesByCode sorts languages by code.
func LanguagesByCode(i, j *Language) bool {
	if i.Code == j.Code {
		return i.Name < j.Name
	}
	return i.Code > j.Code
}

// LanguagesByFiles sorts languages by files.
func LanguagesByFiles(i, j *Language) bool {
	if len(i.Files) == len(j.Files) {
		return i.Name < j.Name
	}
	return len(i.Files) > len(j.Files)
}

// LanguagesBySize sorts languages by size.
func LanguagesBySize(i, j *Language) bool {
	if i.Size == j.Size {
		return i.Name < j.Name
	}
	return i.Size > j.Size
}

// LanguagesByLines sorts languages by lines.
func LanguagesByLines(i, j *Language) bool {
	if i.Lines == j.Lines {
		return i.Name < j.Name
	}
	return i.Lines > j.Lines
}

// LanguagesByComments sorts languages by comments.
func LanguagesByComments(i, j *Language) bool {
	if i.Comments == j.Comments {
		return i.Name < j.Name
	}
	return i.Comments > j.Comments
}

// LanguagesByBlanks sorts languages by blanks.
func LanguagesByBlanks(i, j *Language) bool {
	if i.Blanks == j.Blanks {
		return i.Name < j.Name
	}
	return i.Blanks > j.Blanks
}

// Files sorts
// ===========

// FilesSort is used to sort files.
type FilesSort struct {
	Files   []*File
	LessCmp func(i, j *File) bool
}

func (f FilesSort) Len() int           { return len(f.Files) }
func (f FilesSort) Swap(i, j int)      { f.Files[i], f.Files[j] = f.Files[j], f.Files[i] }
func (f FilesSort) Less(i, j int) bool { return f.LessCmp(f.Files[i], f.Files[j]) }

// FileByCode sorts files by code.
func FileByCode(i, j *File) bool {
	if i.Code == j.Code {
		return i.Name < j.Name
	}
	return i.Code > j.Code
}

// FileBySize sorts files by size.
func FileBySize(i, j *File) bool {
	if i.Size == j.Size {
		return i.Name < j.Name
	}
	return i.Size > j.Size
}

// FileByLines sorts files by lines.
func FileByLines(i, j *File) bool {
	if i.Lines == j.Lines {
		return i.Name < j.Name
	}
	return i.Lines > j.Lines
}

// FileByComments sorts files by comments.
func FileByComments(i, j *File) bool {
	if i.Comments == j.Comments {
		return i.Name < j.Name
	}
	return i.Comments > j.Comments
}

// FileByBlanks sorts files by blanks.
func FileByBlanks(i, j *File) bool {
	if i.Blanks == j.Blanks {
		return i.Name < j.Name
	}
	return i.Blanks > j.Blanks
}
