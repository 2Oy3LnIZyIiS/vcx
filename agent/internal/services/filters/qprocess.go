package filters

import (
	"path/filepath"
	"strings"
	patternlib "vcx/agent/internal/services/filters/pattern"
	"vcx/pkg/set"
	"vcx/pkg/toolkit/pathkit"

	"github.com/bmatcuk/doublestar/v4"
)



func (f *QuickFilter) hasFilter(relativePath string, isDir bool) bool {
    if f.RootPatterns.Contains(relativePath) { return true }

    if f.shouldSkipByExtension(relativePath) { return true }
    if f.shouldSkipByStem(relativePath) { return true }

    segments := pathkit.QuickSplit(relativePath)

    if f.shouldSkipBySegment(segments, isDir){ return true }
    if f.ShouldSkipByKeyword(relativePath, segments, isDir){ return true }
    if f.shouldSkipByRegex(relativePath) { return true }

    return false
}


func (f *QuickFilter) shouldSkipByExtension(relativePath string) bool {
    ext := filepath.Ext(relativePath)
    return len(ext) > 0 && f.ExtensionAny.Contains(ext)
}


func (f *QuickFilter) shouldSkipByStem(relativePath string) bool {
    stem := getFileStem(relativePath)
    return len(stem) > 0 && f.StemAny.Contains(stem)
}


func (f *QuickFilter) shouldSkipBySegment(segments []string, isDir bool) bool {
    lastSegment := len(segments) - 1
    for index, segment := range segments {
        if index == 0 && f.Segment1Any.Contains(segment){ return true }

        isLastSegment := index == lastSegment
        if (isLastSegment && isDir) || !isLastSegment {
            if f.SegmentNDir.Contains(segment) { return true }
        }
        if f.SegmentNAny.Contains(segment) { return true }
    }
    return false
}


func (f *QuickFilter) ShouldSkipByKeyword(path string, segments []string, isDir bool) bool {
    testedPatterns := set.New[string]()

    if isDir {
        path = path + "/"
    }

    for _, segment := range segments {
        if patterns, exists := f.KeywordPatterns[segment]; exists {
            for _, pattern := range patterns {
                if !testedPatterns.Contains(pattern.String()) {
                    if match(path, pattern) { return true }
                    testedPatterns.Add(pattern.String())
                }
            }
        }
    }

    return f.shouldSkipByNormal(path, testedPatterns)
}


func (f *QuickFilter) shouldSkipByNormal( path string, testedPatterns *set.Set[string]) bool {
    for _, pattern := range f.NormalFilePatterns {
        if testedPatterns.Contains(pattern.String()) { continue }
        if match(path, pattern) { return true }
    }
    for _, pattern := range f.NormalDirPatterns {
        if testedPatterns.Contains(pattern.String()) { continue }
        if match(path, pattern) { return true }
    }
    for _, pattern := range f.NormalAnyPatterns {
        if testedPatterns.Contains(pattern.String()) { continue }
        if match(path, pattern) { return true }
    }

    return false
}


func (f *QuickFilter) shouldSkipByRegex(relativePath string) bool {
    if len(f.RegexPatterns) == 0 {
        return false
    }

    for _, pattern := range f.RegexPatterns {
        if pattern.CompiledRegex != nil && pattern.CompiledRegex.MatchString(relativePath) {
            return true
        }
    }

    return false
}



func (f *QuickFilter) hasOverride(relativePath string) bool {
    for _, pattern := range f.Overrides {
        if match(relativePath, pattern) { return true }
    }
    return false
}


func match(path string, pattern *patternlib.Pattern) bool {
    if !pattern.IsDirOnly && strings.HasSuffix(path, "/") {
        return false
    }
    p := pattern.String()
    if pattern.IsOverride {
        p = p[1:]
    }
    if !pattern.IsRoot {
        p = "**/" + p
    }
    matched, err := doublestar.Match(p, path)
    if err == nil && matched { return true }

    if !pattern.IsOverride {
        if pattern.IsDirOnly {
            p = p + "**/*"
        } else {
            p = p + "/**/*"
        }
        matched, err = doublestar.Match(p, path)
        if err == nil && matched { return true }
    }

    return false
}


func getFileStem(path string) string {
    base := filepath.Base(path)
    ext  := filepath.Ext(base)
    if ext == "" {
        return base
    }
    return base[:len(base)-len(ext)]
}
