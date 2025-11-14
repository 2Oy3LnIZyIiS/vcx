package walk

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"vcx/agent/internal/services/filters"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()

// Walk streams directory traversal events to the provided channel
func Walk(dirPath string, eventChan chan<- Event, filter filters.FilterInterface, verbose ...bool) {
	defer close(eventChan)
    log = logging.GetLogger()
    log.Debug("Walking")

	// Default verbose to false
	isVerbose := false
	if len(verbose) > 0 {
		isVerbose = verbose[0]
	}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			eventChan <- Error(fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		// Check filter
		if filter != nil && filter.ShouldSkip(path, d.IsDir()) {
			if isVerbose {
				eventChan <- Skip(path)
			}
			if d.IsDir() {
				// eventChan <- Skip(path)
				return filepath.SkipDir  // Skip directory and its contents
			} else {
				eventChan <- Skip(path)
				return nil  // Skip file
			}
		}

		if d.IsDir() {
			eventChan <- Dir(path)
		} else {
			eventChan <- File(path)
		}

		return nil
	})

	if err != nil {
		eventChan <- Error(fmt.Sprintf("Walk failed: %v", err))
	}
}

/*
    currentPaths := make(map[string]struct{})

    err := filepath.WalkDir(dm.dirPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            dm.errorsChan <- fmt.Errorf("error accessing %s: %v", path, err)
            return nil
        }
        currentPaths[path] = struct{}{}

        // TODO!: handle filters here...
        if d.IsDir() {
            if strings.Contains(path, "node_modules") ||
               strings.Contains(path, ".git") {
                return filepath.SkipDir
            }
        }

        // Process the file asynchronously
        dm.wg.Add(1)
        go dm.analyzeFile(path, d)

        return nil
    })

    if err != nil {
        dm.errorsChan <- err
    }

    dm.checkDeletedFiles(currentPaths)
    dm.wg.Wait()
*/
