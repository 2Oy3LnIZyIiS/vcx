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
	patternlib "vcx/agent/internal/services/filters/pattern"
	"vcx/pkg/logging"
	"vcx/pkg/set"
)


var log = logging.GetLogger()


func NewQuickFilter(instancePath string, patterns []string) *QuickFilter {
	return &QuickFilter{
        InstancePath:       instancePath,
        AllPatterns:        patterns,
        QuickMatchPatterns: set.New[string](),
        NormalFilePatterns: make([]*patternlib.Pattern, 0),
        NormalDirPatterns:  make([]*patternlib.Pattern, 0),
        NormalAnyPatterns:  make([]*patternlib.Pattern, 0),
        HasError:           set.New[string](),

        SegmentNAny:        set.New[string](),
        SegmentNDir:        set.New[string](),
        Segment1Any:        set.New[string](),
        ExtensionAny:       set.New[string](),
        StemAny:            set.New[string](),

        RootPatterns:       set.New[string](),
        KeywordPatterns:    make(map[string][]*patternlib.Pattern),

        RegexPatterns:      make([]*patternlib.Pattern, 0),
        Overrides:          make([]*patternlib.Pattern, 0),
	}
}


func (f *QuickFilter) Init() {
    log = logging.GetLogger()
    log.Debug("Initilizing Filter")

    // TEMP TEMP TEMP
    startTime := time.Now()
    // TEMP TEMP TEMP

    f.clearAllData()

    log.Debug(fmt.Sprintf("Patterns before normalization: %d", len(f.AllPatterns)))
    f.AllPatterns = patternlib.Normalize(f.AllPatterns)
    log.Debug(fmt.Sprintf("Patterns after normalization: %d", len(f.AllPatterns)))

    for _, pattern := range f.AllPatterns {
        f.ProcessPattern(pattern)
    }

    // TEMP TEMP TEMP
    f.logStats(time.Since(startTime))
    // TEMP TEMP TEMP
}


func (f *QuickFilter) ShouldSkip(path string, isDir bool) bool {
    relativePath := path[len(f.InstancePath):]
    if relativePath == "" { return false }
    // e.g.: /projectdir/logs/log.log

    if f.hasFilter(relativePath, isDir) {
        if strings.Contains(relativePath, "t32.exe"){
            log.Debug("hit")
        }
        if strings.Contains(relativePath, "icon.ico"){
            log.Debug("hit")
        }
        return !f.hasOverride(relativePath)
    }


    // * Test spaces in directories

    return false
}
