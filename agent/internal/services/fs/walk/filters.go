package walk

import (
	"path/filepath"
	"strings"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// SimpleFilter implements basic pattern matching
type SimpleFilter struct {
	patterns []string
}

func NewSimpleFilter(patterns []string) *SimpleFilter {
	return &SimpleFilter{patterns: patterns}
}

func (f *SimpleFilter) ShouldSkip(path string, isDir bool) bool {
	name := filepath.Base(path)

	for _, pattern := range f.patterns {
		// Simple pattern matching
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}

		// Check if path contains pattern (for nested paths)
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

// DefaultFilter with common ignore patterns
func DefaultFilter() *SimpleFilter {
	return NewSimpleFilter([]string{
		"node_modules",
		".git",
		".DS_Store",
		"*.tmp",
		"*.log",
		".vcx",  // VCX internal directory
	})
}
