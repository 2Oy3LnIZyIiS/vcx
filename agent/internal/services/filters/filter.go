package filters

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// BasicFilter provides correct gitignore matching using go-git gitignore
// This filter prioritizes correctness over performance and serves as validation
// for QuickFilter optimizations
// NOTE: The original point of this filter was to test quickfilter for correctness
// .. BUT there is actually a bug in the go-git gitignore code that does not
// .. treat patterns like /doc/_build/ correctly.  Some VERY brief investigation
// .. leads me to guess that that go-git's usage of the glob match is what's
// .. causing this.  glob match works in many cases but not all cases and I
// .. didn't see the edge case code to handle this situation
type BasicFilter struct {
    InstancePath string
	patterns     []gitignore.Pattern
	matcher      gitignore.Matcher
}

func NewBasicFilter(instancePath string, patternStrings []string) *BasicFilter {
	var patterns []gitignore.Pattern
	for _, p := range patternStrings {
		if p != "" && !strings.HasPrefix(p, "#") {
			patterns = append(patterns, gitignore.ParsePattern(p, nil))
		}
	}
	return &BasicFilter{
        InstancePath: instancePath,
		patterns: patterns,
		matcher: gitignore.NewMatcher(patterns),
	}
}

func (f *BasicFilter) Init() {
	// Patterns are already processed in NewBasicFilter
	// Matcher is already created
}

func (f *BasicFilter) ShouldSkip(path string, isDir bool) bool {
    relativePath := strings.TrimPrefix(path[len(f.InstancePath):], "/")
    if relativePath == "" {
        return false
    }

    // Split path into segments for go-git matcher
    pathSegments := strings.Split(relativePath, "/")

    // go-git matcher expects the path segments and isDir flag
    result := f.matcher.Match(pathSegments, isDir)
    return result
}


// FromFileBasic creates a BasicFilter from a gitignore file
func FromFileBasic(instancePath string, gitignorePath string) FilterInterface {
	file, err := os.Open(gitignorePath)
	if err != nil {
		return DefaultBasicFilter(instancePath)
	}
	defer file.Close()

	// Read patterns manually since ReadPatterns expects fs.FS
	var patternStrings []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		patternStrings = append(patternStrings, line)
	}

	return NewBasicFilter(instancePath, patternStrings)
}

// DefaultBasicFilter with common ignore patterns
func DefaultBasicFilter(instancePath string) *BasicFilter {
	return NewBasicFilter(instancePath, []string{
		"node_modules",
		".git",
		".DS_Store",
		"*.tmp",
		"*.log",
		".vcx",  // VCX internal directory
	})
}
