package filters

import (
	patternlib "vcx/agent/internal/services/filters/pattern"
)


func (f *QuickFilter) clearAllData() {
    f.QuickMatchPatterns.Clear()
    f.NormalFilePatterns = make([]*patternlib.Pattern, 0)
    f.NormalDirPatterns  = make([]*patternlib.Pattern, 0)
    f.NormalAnyPatterns  = make([]*patternlib.Pattern, 0)
    f.HasError.Clear()
    f.SegmentNAny.Clear()
    f.SegmentNDir.Clear()
    f.Segment1Any.Clear()
    f.ExtensionAny.Clear()
    f.StemAny.Clear()
    f.RootPatterns.Clear()
    f.KeywordPatterns    = make(map[string][]*patternlib.Pattern)
    f.RegexPatterns      = make([]*patternlib.Pattern, 0)
    f.Overrides          = make([]*patternlib.Pattern, 0)
}


func (f *QuickFilter) ProcessPattern(p string) {
    pattern := patternlib.Parse(p)
    if pattern.IsRegex {
        if pattern.IsValid {
            f.RegexPatterns = append(f.RegexPatterns, pattern)
        }
        return
    }
    if pattern.IsOverride {
        f.Overrides = append(f.Overrides, pattern)
        return
    }

    if !pattern.IsValid {
        f.HasError.Add(pattern.String())
        return
    }

    if pattern.IsRoot && !pattern.HasWildcard && !pattern.HasExpansion{
        f.RootPatterns.Add(pattern.String())
    }

    if f.isQuickMatch(pattern) {
        f.QuickMatchPatterns.Add(pattern.String())
        return
    }

    for _, token := range pattern.Tokens {
        f.KeywordPatterns[token] = append(f.KeywordPatterns[token], pattern)
    }

    if pattern.IsFileOnly {
        f.NormalFilePatterns = append(f.NormalFilePatterns, pattern)
    } else if pattern.IsDirOnly {
        f.NormalDirPatterns = append(f.NormalDirPatterns, pattern)
    } else {
        f.NormalAnyPatterns = append(f.NormalAnyPatterns, pattern)
    }
}


func (f *QuickFilter) isQuickMatch(pattern *patternlib.Pattern) bool {
    if pattern.IsRegex || pattern.IsOverride {
        return false
    }
    if pattern.IsStem && pattern.IsExtension {
        return false
    }
    if pattern.IsStem  && !pattern.IsDirOnly && !pattern.HasWildcard && !pattern.HasExpansion{
        f.StemAny.Add(pattern.Normal())
        return true
    }

    if pattern.IsExtension && !pattern.IsDirOnly && !pattern.HasWildcard && !pattern.HasExpansion {
        f.ExtensionAny.Add(pattern.Normal())
        return true
    }

    if pattern.IsSingular {
        // Single Segment pattern
        if pattern.IsRoot {
            if !pattern.IsFileOnly && !pattern.IsDirOnly && !pattern.HasWildcard && !pattern.HasExpansion {
                // NOTE: this eliminates patterns like /build/ from quick matches
                // I don't think this pattern is typical
                // can always add another quickmatch later
                f.Segment1Any.Add(pattern.Normal())
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
