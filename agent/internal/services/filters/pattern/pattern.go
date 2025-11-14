package pattern

import (
	"regexp"
	"strings"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// type Buckets struct{
//     // Buckets of rules for filter speed optimization (fastest to slowest)

//     // Ultra-fast: String operations only (~1-10ns)
//     RootFileOnly    []Rule  // /build.log - exact string match
//     RootDirOnly     []Rule  // /src/ - prefix string match
//     RootEither      []Rule  // /path/file - exact string match
//     FileOnly        []Rule  // file.log - filename string match
//     DirOnly         []Rule  // dirname/ - directory string match
//     Either          []Rule  // path/file - string contains/suffix

//     // Fast: Segmented wildcards with string pre-filtering (~20-100ns)
//     RootFileOnlyWildcard []Rule  // /src/*.log - check "src/" first
//     RootDirOnlyWildcard  []Rule  // /build/*/ - check "build/" first
//     RootEitherWildcard   []Rule  // /src/*/file - check "src/" and "file"
//     FileOnlyWildcard     []Rule  // *.log - direct filename pattern
//     DirOnlyWildcard      []Rule  // temp*/ - directory pattern
//     EitherWildcard       []Rule  // src/*.log - check "src/" first

//     // Slowest: Full regex (~1000-5000ns)
//     Regex           []Rule  // ~pattern - full regex matching
// }


type Pattern struct {
	Text         string
    Segments     []string
    Tokens       []string
	IsRegex      bool
	IsOverride   bool
    IsFileOnly   bool
	IsDirOnly    bool
	IsRoot       bool
    HasWildcard  bool
    HasExpansion bool
    IsValid      bool

    // Common pattern types
    IsSingular  bool
    IsExtension bool
    IsStem      bool

    // Compiled regex for performance
    CompiledRegex *regexp.Regexp
}


// type Segment struct {
// 	pattern      string
//     isLiteral    bool
// }


func Parse(pattern string) *Pattern {
    log = logging.GetLogger()
    p := &Pattern{ Text:    pattern,
                   IsValid: true,
    }
    if len(pattern) == 0 {
        log.Debug(pattern)
        return p
    }

    if strings.HasPrefix(pattern, "r:"){
        p.IsRegex = true
        // Compile regex during preprocessing for performance
        regexPattern := pattern[2:] // Remove "r:" prefix
        compiled, err := regexp.Compile(regexPattern)
        if err != nil {
            log.Debug("Invalid regex pattern: " + pattern + " - " + err.Error())
            p.IsValid = false
        } else {
            p.CompiledRegex = compiled
        }
        return p
    }

    p.IsOverride   = pattern[0] == '!'
    p.IsFileOnly   = pattern[0] == '@'
    p.IsDirOnly    = pattern[len(pattern)-1] == '/'
    p.IsRoot       = pattern[0] == '/' || strings.HasPrefix(pattern, "@/")
    p.HasWildcard  = strings.ContainsAny(pattern, "*?")
    p.HasExpansion = strings.ContainsAny(pattern, "[]{}")

    if !p.validate() {
        return p
    }

    p.tokenize()

    p.IsSingular  = len(p.Segments) == 1
    p.IsExtension = p.IsSingular && strings.HasPrefix(p.Text, "*.")
    p.IsStem      = p.IsSingular && strings.HasSuffix(p.Text, ".*")
    if p.IsExtension && !strings.ContainsAny(pattern[2:], "*?") {
        p.HasWildcard = false
    }
    if p.IsStem && !strings.ContainsAny(pattern[:len(pattern)-2], "*?") {
        p.HasWildcard = false
    }
    return p
}


func (p *Pattern) validate() bool{
    if p.IsFileOnly && p.IsDirOnly {
        p.IsValid = false
    }
    return p.IsValid
}


func (p *Pattern) tokenize() {
    for segment := range strings.SplitSeq(p.Text, "/") {
        if segment == "" || segment == "@" {
            continue
        }
        p.Segments = append(p.Segments, segment)
        if strings.ContainsAny(segment, "*[]?") {
           continue
        }
        p.Tokens = append(p.Tokens, segment)
    }
}



func (p *Pattern) String() string {
    return p.Text
}


func (p *Pattern) Normal() string {
    var normalText string = p.Text
    if p.IsOverride {
        return normalText[1:]
    }
    if p.IsRegex {
        return normalText[2:]
    }
    if p.IsFileOnly{
        normalText = normalText[1:]
    }
    if p.IsRoot{
        normalText = normalText[1:]
    }
    if p.IsDirOnly {
        normalText = normalText[:len(normalText)-1]
    }
    if p.IsStem && p.IsExtension {
        return normalText
    }
    if p.HasExpansion{
        return normalText
    }
    if p.HasWildcard{
        return normalText
    }
    if p.IsStem {
        return normalText[:len(normalText)-2]
    }
    if p.IsExtension {
        return normalText[1:]
    }
    return normalText
}


// func extractDirConstraints(pattern string) []string {
//     if strings.HasPrefix(pattern, "/") {
//         // Root-relative: /dir1/**/dir2 constrains to exact path from root
//         return []string{pattern} // full path constraint
//     } else {
//         // Non-root: dir1/**/dir2 can match anywhere
//         // Extract all directory components as potential constraints
//         parts := strings.Split(pattern, "/")
//         var constraints []string
//         for _, part := range parts {
//             if part != "**" && !strings.Contains(part, "*") {
//                 constraints = append(constraints, part)
//             }
//         }
//         return constraints
//     }
// }
