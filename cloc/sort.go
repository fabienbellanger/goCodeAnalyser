package cloc

// Languages sorts
// ===============

// LanguagesByCode represents an array of Language for sorting by code.
type LanguagesByCode []Language

// Sort by Code
func (l LanguagesByCode) Len() int {
	return len(l)
}
func (l LanguagesByCode) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l LanguagesByCode) Less(i, j int) bool {
	if l[i].Code == l[j].Code {
		return l[i].Name < l[j].Name
	}
	return l[i].Code > l[j].Code
}

// Files sorts
// ===========

// FilesByCode represents an array of File for sorting by code.
type FilesByCode []File

// Sort by Code
func (f FilesByCode) Len() int {
	return len(f)
}
func (f FilesByCode) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f FilesByCode) Less(i, j int) bool {
	if f[i].Code == f[j].Code {
		return f[i].Name < f[j].Name
	}
	return f[i].Code > f[j].Code
}
