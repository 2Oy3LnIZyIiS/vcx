package filters

import (
	patternlib "vcx/agent/internal/services/filters/pattern"
	"vcx/pkg/set"
)


type QuickFilter struct {
    //Pattern Storage
    InstancePath       string
	AllPatterns        []string          // original normalized pattern list
	QuickMatchPatterns *set.Set[string]  // Patterns in any of the style based sets

    NormalFilePatterns []*patternlib.Pattern
    NormalDirPatterns  []*patternlib.Pattern
    NormalAnyPatterns  []*patternlib.Pattern
    HasError           *set.Set[string]

    // Style based Sets
    // These patterns will also be copied to QuickMatchPatterns
    SegmentNAny        *set.Set[string]
    SegmentNDir        *set.Set[string]
    Segment1Any        *set.Set[string]
    ExtensionAny       *set.Set[string]
    StemAny            *set.Set[string]

    // Additional Sets
    // These patterns will also be copied to one of the Normal*Patterns Sets
    RootPatterns       *set.Set[string]
    KeywordPatterns    map[string][]*patternlib.Pattern

    // Special Sets
    RegexPatterns      []*patternlib.Pattern  // Store compiled regex patterns
    Overrides          []*patternlib.Pattern
}
