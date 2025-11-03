package pattern

import (
	"strings"
	"vcx/pkg/set"
)

// Patterns requiring preprocessing:
// patterns := []string{
//     "",                // Empty lines (ignored) - fastest check
//     "# comment",       // Comments (ignored) - simple prefix check
//     "~.*\\.log$",      // Regex patterns - detected and skipped from preprocessing
//     "*.log # comment", // Inline comments - strip after #
//     "*.{jpg,png}",     // Brace expansion - multiple expansions
// }
//
// Patterns NOT requiring preprocessing:
// patterns := []string{
//     "filterme.log",    // Literal match - exact filename
//     "file\\#name.txt", // Escaped # - literal # character
//     "\\*.log",         // Escaped characters - literal match
//     "/build",          // Root-relative (leading slash) - exact match
//     "node_modules/",   // Directory-only (trailing slash) - exact match
//     "*.log",           // Simple wildcards - basic pattern
//     "**/*.tmp",        // Recursive wildcards - same pattern complexity
//     "dir/**/file.log", // Complex recursive wildcards - nested path matching
//     "dir/**/dir2/",    // Recursive directory matching at any depth
//     "file?.txt",       // Question mark wildcards - single char
//     "temp[0-9]*",      // Character classes - simple range
//     "temp[a-z]*",      // Range character classes - simple range
//     "file[!0-9]*",     // Negated character classes - negated range
//     "@*.log",          // Files-only patterns - matches files but not directories
//     "!important.log",  // Negation patterns - just prefix detection
// }


func Normalize(patterns []string) []string {
    var result []string
    var seen = set.New[string]()
    for _, pattern := range patterns {
        pattern = strings.TrimSpace(pattern)

        if isEmpty(pattern) || isComment(pattern) || isDangerous(pattern) {
            continue
        }
        if seen.Contains(pattern) {
            continue
        }
        seen.Add(pattern)

        if isRegex(pattern) {
            result = append(result, pattern)
            continue
        }

        if strings.Contains(pattern, "#") {
            pattern = removeInlineComment(pattern)
        }

        if strings.HasPrefix(pattern, "/**/") || strings.HasPrefix(pattern, "**/") {
            pattern = removeDoubleStars(pattern)
            if isEmpty(pattern){
                continue
            }
        }

        if strings.Contains(pattern, "{") {
            result = append(result, expandBrace(pattern)...)
        } else {
            result = append(result, pattern)
        }
    }
    return result
}


func isEmpty(s string) bool {
    return s == ""
}


func isComment(s string) bool {
    return s[0] == '#'
}


func isRegex(s string) bool {
    return s[0] == '~'
}


func removeInlineComment(s string) string {
    // i cannot be 0 since that would be a normal comment
    for i := 1; i < len(s); i++ {
        if s[i] == '#' && s[i-1] != '\\' {
            return strings.TrimSpace(s[:i])
        }
    }
    return s
}


func removeDoubleStars(s string) string {
    if strings.HasPrefix(s, "/**/") {
        return s[4:]
    }
    if strings.HasPrefix(s, "**/") {
        return s[3:]
    }
    return s
}


func unescape(s string) string {
    var result strings.Builder
    result.Grow(len(s))

    for i := 0; i < len(s); i++ {
        if s[i] == '\\' && i+1 < len(s) {
            next := s[i+1]
            // NOTE: if statement chain is likely more efficient than a map lookup
            if next == '#' || next == '*' || next == '?' || next == '[' || next == '\\' || next == ' ' {
                result.WriteByte(next)
                i++ // skip the escaped character
            } else {
                result.WriteByte('\\')
            }
        } else {
            result.WriteByte(s[i])
        }
    }

    return result.String()
}


func expandBrace(pattern string) []string {
    start := strings.Index(pattern, "{")
    if start == -1 {
        return []string{pattern}
    }

    end := strings.Index(pattern[start:], "}")
    if end == -1 {
        return []string{pattern}
    }
    end += start

    prefix  := pattern[:start]
    suffix  := pattern[end+1:]
    options := strings.Split(pattern[start+1:end], ",")

    var results []string
    for _, option := range options {
        expanded := prefix + strings.TrimSpace(option) + suffix
        // Recursively expand in case there are nested braces
        results = append(results, expandBrace(expanded)...)
    }

    return results
}


// func isFileOnly(s string) bool {
//     return s[0] == '@'
// }
func isDangerous(s string) bool {
    switch s {
    case "*", "**", "/", ".", "..":
        return true
    default:
        return false
    }
}
