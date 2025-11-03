package filters

import (
	"path/filepath"
	"strings"
	"vcx/pkg/set"
)

// NOTE: Am choosing NOT to deal with escaped characters at this time

type Filter struct {
	patterns          []string
    RootPatterns      *set.Set[string]


    RootUnspecified   *set.Set[string]
    RootFile          *set.Set[string]
    RootDir           *set.Set[string]
    SingleFile        *set.Set[string]
    SingleDir         *set.Set[string]
    SingleUnspecified *set.Set[string]
}

func NewFilter(patterns []string) *Filter {
	return &Filter{
        patterns:          patterns,
        RootPatterns: set.New[string](),

        RootUnspecified:   set.New[string](),
        RootFile:          set.New[string](),
        RootDir:           set.New[string](),
        SingleFile:        set.New[string](),
        SingleDir:         set.New[string](),
        SingleUnspecified: set.New[string](),
    }
}


func (f *Filter) Init() {}


func (f *Filter) Process() {
    for _, pattern := range f.patterns {
        if !strings.ContainsAny(pattern, "*[]?") {
            if f.isRootPattern(pattern) {
                continue
            }
            if f.isSinglePattern(pattern){
                continue
            }
        }
        // process the rest
    }
}

func (f *Filter) isRootPattern(pattern string) bool {
    if strings.HasPrefix(pattern, "@/"){
        f.RootPatterns.Add(pattern[1:])
        return true
    }
    if strings.HasPrefix(pattern, "/") {
        f.RootPatterns.Add(pattern)
        return true
    }

    // if strings.HasPrefix(pattern, "@/") {
    //     f.RootFile.Add(pattern[1:])
    //     return true
    // }
    // if strings.HasPrefix(pattern, "/") {
    //     if strings.HasSuffix(pattern, "/") {
    //         f.RootDir.Add(pattern[:len(pattern) - 1])
    //     } else {
    //         f.RootUnspecified.Add(pattern)
    //     }
    //     return true
    // }
    return false
}


func (f *Filter) isSinglePattern(pattern string) bool {
    if strings.Contains(pattern[:len(pattern) - 1], "/") {
        return false
    }
    if strings.HasPrefix(pattern, "@") {
        f.SingleFile.Add(pattern[1:])
        return true
    }
    if strings.HasSuffix(pattern, "/") {
        f.SingleDir.Add(pattern[:len(pattern)-1])
        return true
    }
    f.SingleUnspecified.Add(pattern)
    return true
}



func (f *Filter) ShouldSkip(path string, isDir bool) bool {
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

// DefaultSimpleFilter with common ignore patterns
func DefaultSimpleFilter() *Filter {
	return NewFilter([]string{
		"node_modules",
		".git",
		".DS_Store",
		"*.tmp",
		"*.log",
		".vcx",  // VCX internal directory
	})
}
