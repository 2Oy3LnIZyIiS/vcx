/*
quickfilter.go

Quick Filter is a first pass style filter that can reliable provide a DEFINITE MATCH
or a POTENTIAL MATCH based on common ignore pattern styles.

While there are many ignore pattern styles, the quickfilter focuses on:
- build  - Single Segment pattern that can match to any path segment
- build/ - Single Segment pattern that can match to any directory segment
- /build - Single Segment pattern that matches to any root level segment
- temp.* - Broad Stem pattern
- *.tmp  - Broad Extension pattern
The above pattern styles are consolidated into a set that a path can be quickly
matched against.  This provides a DEFINITE MATCH and the path can be ignored

In addition, two additional checks will be made:
- If there is no above match, then the path will compared to the set of all root patterns.
This is quick string comparison and will provide a DEFINITE MATCH or NOT

- An additional map is created that contains the set of all possible matches and
their corresponding patterns.  If any path segment matches against this set then
we have a POTENTIAL Match and the potential match patterns will be tested first.
If there is a match, the path is ignore.  If there is no match then the path will
be tested agains all remaining filters minus the ones used to create the quick
matches

*/

package filters

import (
	"fmt"
	"strings"
	"time"
	"vcx/agent/internal/services/filters/pattern"
	patternlib "vcx/agent/internal/services/filters/pattern"
	"vcx/pkg/logging"
	"vcx/pkg/set"
)


var log = logging.GetLogger()


type QuickFilter struct {
    //Pattern Storage
	AllPatterns        []string          // original normalized pattern list
	QuickMatchPatterns *set.Set[string]  // Patterns in any of the style based sets

    NormalFilePatterns *set.Set[string]
    NormalDirPatterns  *set.Set[string]
    NormalAnyPatterns  *set.Set[string]
    HasError           *set.Set[string]

    // Style based Sets
    // These patterns will also be copied to QuickMatchPatterns
    SegmentNAny        *set.Set[string]
    SegmentNDir        *set.Set[string]
    Segment1Any        *set.Set[string]
    FileExtensions         *set.Set[string]
    FileStem              *set.Set[string]

    // Additional Sets
    // These patterns will also be copied to one of the Normal*Patterns Sets
    RootPatterns       *set.Set[string]
    KeywordPatterns    map[string][]*pattern.Pattern

    // Special Sets
    RegexPatterns      *set.Set[string]
    Overrides          *set.Set[string]
}

func NewQuickFilter(patterns []string) *QuickFilter {
	return &QuickFilter{
        AllPatterns:        patterns,
        QuickMatchPatterns: set.New[string](),
        NormalFilePatterns: set.New[string](),
        NormalDirPatterns:  set.New[string](),
        NormalAnyPatterns:  set.New[string](),
        HasError:           set.New[string](),

        SegmentNAny:        set.New[string](),
        SegmentNDir:        set.New[string](),
        Segment1Any:        set.New[string](),
        FileExtensions:         set.New[string](),
        FileStem:              set.New[string](),

        RootPatterns:       set.New[string](),
        KeywordPatterns:    make(map[string][]*patternlib.Pattern),

        RegexPatterns:      set.New[string](),
        Overrides:          set.New[string](),
	}
}


func (f *QuickFilter) Init() {
    log = logging.GetLogger()
    log.Debug("Initilizing Filter")
    startTime := time.Now()
    f.ResetSets()
    log.Debug(fmt.Sprintf("Patterns before normalization: %d", len(f.AllPatterns)))
    f.AllPatterns = pattern.Normalize(f.AllPatterns)
    log.Debug(fmt.Sprintf("Patterns after normalization: %d", len(f.AllPatterns)))
    for _, pattern := range f.AllPatterns {
        f.ProcessPattern(pattern)
    }
    duration := time.Since(startTime)
    f.logStats(duration)
}


func (f *QuickFilter) logStats(duration time.Duration) {
    total := len(f.AllPatterns)
    quickCount := f.QuickMatchPatterns.Size()
    normalCount := f.NormalFilePatterns.Size() + f.NormalDirPatterns.Size() + f.NormalAnyPatterns.Size()
    regexCount := f.RegexPatterns.Size()
    overrideCount := f.Overrides.Size()
    errorCount := f.HasError.Size()

    accounted := quickCount + normalCount + regexCount + overrideCount + errorCount
    missing := total - accounted

    log.Debug(fmt.Sprintf("Filter initialized in %s - %d patterns processed", duration, total))
    log.Debug(fmt.Sprintf("Quick: %d (%.1f%%) | Normal: %d (%.1f%%) | Regex: %d | Override: %d | Error: %d",
        quickCount, float64(quickCount)/float64(total)*100,
        normalCount, float64(normalCount)/float64(total)*100,
        regexCount, overrideCount, errorCount))

    if missing != 0 {
        log.Debug(fmt.Sprintf("WARNING: %d patterns unaccounted for!", missing))
        f.checkForDuplicates()
    }

    // Quick match buckets (only show non-empty)
    if f.FileExtensions.Size() > 0 {
        log.Debug(fmt.Sprintf("  Extensions: %d [%s]", f.FileExtensions.Size(), strings.Join(f.FileExtensions.Items(), " | ")))
    }
    if f.FileStem.Size() > 0 {
        log.Debug(fmt.Sprintf("  Stems: %d [%s]", f.FileStem.Size(), strings.Join(f.FileStem.Items(), " | ")))
    }
    if f.SegmentNAny.Size() > 0 {
        log.Debug(fmt.Sprintf("  Any segments: %d [%s]", f.SegmentNAny.Size(), strings.Join(f.SegmentNAny.Items(), " | ")))
    }
    if f.SegmentNDir.Size() > 0 {
        log.Debug(fmt.Sprintf("  Dir segments: %d [%s]", f.SegmentNDir.Size(), strings.Join(f.SegmentNDir.Items(), " | ")))
    }
    if f.Segment1Any.Size() > 0 {
        log.Debug(fmt.Sprintf("  Root segments: %d [%s]", f.Segment1Any.Size(), strings.Join(f.Segment1Any.Items(), " | ")))
    }

    // Special patterns (always show counts, content if present)
    if regexCount > 0 {
        log.Debug(fmt.Sprintf("Regex patterns: %d [%s]", regexCount, strings.Join(f.RegexPatterns.Items(), " | ")))
    }
    if overrideCount > 0 {
        log.Debug(fmt.Sprintf("Override patterns: %d [%s]", overrideCount, strings.Join(f.Overrides.Items(), " | ")))
    }
    if errorCount > 0 {
        log.Debug(fmt.Sprintf("ERROR patterns: %d [%s]", errorCount, strings.Join(f.HasError.Items(), " | ")))
    }

    // Normal pattern sets (only show if non-empty)
    if f.NormalFilePatterns.Size() > 0 {
        log.Debug(fmt.Sprintf("  Normal File: %d [%s]", f.NormalFilePatterns.Size(), strings.Join(f.NormalFilePatterns.Items(), " | ")))
    }
    if f.NormalDirPatterns.Size() > 0 {
        log.Debug(fmt.Sprintf("  Normal Dir: %d [%s]", f.NormalDirPatterns.Size(), strings.Join(f.NormalDirPatterns.Items(), " | ")))
    }
    if f.NormalAnyPatterns.Size() > 0 {
        log.Debug(fmt.Sprintf("  Normal Any: %d [%s]", f.NormalAnyPatterns.Size(), strings.Join(f.NormalAnyPatterns.Items(), " | ")))
    }
    
    // Keywords (only show count if many)
    if len(f.KeywordPatterns) > 0 {
        log.Debug(fmt.Sprintf("Keyword index: %d unique tokens", len(f.KeywordPatterns)))
    }
}

func (f *QuickFilter) checkForDuplicates() {
    seen := make(map[string]int)
    duplicates := 0

    for _, pattern := range f.AllPatterns {
        seen[pattern]++
        if seen[pattern] == 2 {
            duplicates++
            log.Debug(fmt.Sprintf("  Duplicate pattern: %s", pattern))
        }
    }

    if duplicates > 0 {
        log.Debug(fmt.Sprintf("Found %d duplicate patterns (may explain missing count)", duplicates))
    } else {
        log.Debug("No duplicates found - missing patterns due to processing logic")
    }
}

func (f *QuickFilter) ResetSets() {
    f.QuickMatchPatterns.Clear()
    f.NormalFilePatterns.Clear()
    f.NormalDirPatterns.Clear()
    f.NormalAnyPatterns.Clear()
    f.HasError.Clear()
    f.SegmentNAny.Clear()
    f.SegmentNDir.Clear()
    f.Segment1Any.Clear()
    f.FileExtensions.Clear()
    f.FileStem.Clear()
    f.RootPatterns.Clear()
    f.KeywordPatterns = make(map[string][]*patternlib.Pattern)
    f.RegexPatterns.Clear()
    f.Overrides.Clear()
}


func (f *QuickFilter) ProcessPattern(p string) {
    pattern := patternlib.Parse(p)
    if pattern.IsRegex {
        f.RegexPatterns.Add(pattern.Normal())
        return
    }
    if pattern.IsOverride {
        f.Overrides.Add(pattern.Normal())
        return
    }

    if !pattern.IsValid {
        f.HasError.Add(pattern.String())
        return
    }

    if f.isQuickMatch(pattern) {
        f.QuickMatchPatterns.Add(pattern.String())
        return
    }

    // pattern will be added to one of the normal sets
    // pattern may also be added to RootPatterns and KeywordPatterns
    if pattern.IsRoot && !pattern.HasWildcard && !pattern.HasExpansion{
        f.RootPatterns.Add(pattern.String())
    }

    for _, token := range pattern.Tokens {
        f.KeywordPatterns[token] = append(f.KeywordPatterns[token], pattern)
    }

    if pattern.IsFileOnly {
        f.NormalFilePatterns.Add(pattern.Normal())
    } else if pattern.IsDirOnly {
        f.NormalDirPatterns.Add(pattern.Normal())
    } else {
        f.NormalAnyPatterns.Add(pattern.Normal())
    }
}


func (f *QuickFilter) isQuickMatch(pattern *patternlib.Pattern) bool {
    if pattern.IsRegex || pattern.IsOverride {
        return false
    }
    if pattern.IsStem && pattern.IsExtension {
        return false
    }
    if pattern.IsStem  && !pattern.IsDirOnly && !pattern.HasExpansion {
        f.FileStem.Add(pattern.Normal())
        return true
    }

    if pattern.IsExtension && !pattern.IsDirOnly && !pattern.HasExpansion {
        f.FileExtensions.Add(pattern.Normal())
        return true
    }

    if pattern.IsSingular {
        // Single Segment pattern
        if pattern.IsRoot {
            if !pattern.IsFileOnly && !pattern.IsDirOnly {
                // NOTE: this eliminates patterns like /build/ from quick matches
                // but it is my opinion that this pattern isn't typical
                // can always add another quickmatch later
                f.Segment1Any.Add(pattern.String())
                return true
            }
        } else {
            if pattern.IsDirOnly && !pattern.HasWildcard && !pattern.HasExpansion {
                f.SegmentNDir.Add(pattern.Normal())
                return true
            } else {
                if !pattern.IsFileOnly && !pattern.HasWildcard && !pattern.HasExpansion {
                    f.SegmentNAny.Add(pattern.Normal())
                    return true
                }
            }
        }
    }

    return false
}


func (f *QuickFilter) ShouldSkip(path string, isDir bool) bool {
    // * Test spaces in directories
    return false
}
