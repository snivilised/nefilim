package nef

import (
	"path/filepath"
	"strings"
)

// Path
type PathCalc interface {
	Base(path string) string
	Clean(path string) string
	Dir(name string) string
	Elements(path string) []string
	Join(elements ...string) string
	Split(path string) (dir, file string)
	Truncate(path string) string
}

type AbsoluteCalc struct {
}

// Base returns the last element of the path
func (c *AbsoluteCalc) Base(path string) string {
	return filepath.Base(path)
}

// Clean returns the shortest path name equivalent to path
// by purely lexical processing.
func (c *AbsoluteCalc) Clean(path string) string {
	return filepath.Clean(path)
}

// Dir returns all but the last element of the path
func (c *AbsoluteCalc) Dir(name string) string {
	return filepath.Dir(name)
}

// Join joins any number of path elements into a single path
func (c *AbsoluteCalc) Join(elements ...string) string {
	return filepath.Join(elements...)
}

// Split splits the path immediately following the final separator
func (c *AbsoluteCalc) Split(path string) (dir, file string) {
	return filepath.Split(path)
}

func (c *AbsoluteCalc) Truncate(path string) string {
	if path == "" {
		return "."
	}

	var (
		separator = string(filepath.Separator)
	)

	if !strings.HasSuffix(path, separator) {
		return path
	}

	return path[:strings.LastIndex(path, separator)]
}

func (c *AbsoluteCalc) Elements(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, string(filepath.Separator))
}

const (
	separator = '/'
)

var (
	separatorStr = string(separator)
)

type RelativeCalc struct {
	Root string
}

// Base returns the last element of the path
func (c *RelativeCalc) Base(path string) string {
	if path == "" {
		return "."
	}

	if !strings.Contains(path, separatorStr) {
		return path
	}

	return path[strings.LastIndex(path, separatorStr)+1:]
}

// Clean returns the shortest path name equivalent to path
// by purely lexical processing.
func (c *RelativeCalc) Clean(path string) string {
	clean := filepath.Clean(path)

	if clean == separatorStr {
		return "."
	}

	if strings.HasPrefix(clean, separatorStr) {
		clean = clean[1:]
	}

	return clean
}

func (c *RelativeCalc) Elements(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, separatorStr)
}

// Dir returns all but the last element of the path
func (c *RelativeCalc) Dir(path string) string {
	if path == "" {
		return "."
	}

	if !strings.Contains(path, separatorStr) {
		return "."
	}

	return path[:strings.LastIndex(path, separatorStr)]
}

// Join joins any number of path elements into a single path
func (c *RelativeCalc) Join(elements ...string) string {
	return strings.Join(elements, separatorStr)
}

// Split splits the path immediately following the final separator
func (c *RelativeCalc) Split(path string) (dir, file string) {
	if path == "" {
		return "", ""
	}

	if !strings.Contains(path, separatorStr) {
		return "", path
	}

	return c.Dir(path), c.Base(path)
}

func (c *RelativeCalc) Truncate(path string) string {
	if path == "" {
		return "."
	}

	if !strings.HasSuffix(path, separatorStr) {
		return path
	}

	return path[:strings.LastIndex(path, separatorStr)]
}
