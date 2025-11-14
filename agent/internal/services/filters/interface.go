package filters

import (
	"bufio"
	"os"
	"strings"
	"vcx/pkg/logging"
)

// FilterInterface defines the interface for path filtering
type FilterInterface interface {
	ShouldSkip(path string, isDir bool) bool
    Init()
}

// DefaultFilter returns the recommended filter for VCX
func DefaultFilter(instancePath string) FilterInterface {
	// return DefaultSimpleFilter()
    return DefaultQuickFilter(instancePath)
}

func FromFile(instancePath string, path string) FilterInterface {
	file, err := os.Open(path)
	if err != nil {
        log := logging.GetLogger()
        log.Error("Could not open file")
		return DefaultQuickFilter(instancePath)
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return NewQuickFilter(instancePath, patterns)
}
