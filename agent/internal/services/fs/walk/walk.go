package walk

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)


type WalkEvent struct {
	Type    string `json:"type"`    // "file", "dir", "error", "complete"
	Path    string `json:"path"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Walk streams directory traversal events to the provided channel
func Walk(dirPath string, eventChan chan<- WalkEvent) {
	defer close(eventChan)

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			eventChan <- WalkEvent{
				Type:  "error",
				Path:  path,
				Error: err.Error(),
			}
			return nil
		}

		// Skip common directories
		if d.IsDir() {
			if strings.Contains(path, "node_modules") ||
				strings.Contains(path, ".git") {
				eventChan <- WalkEvent{
					Type:    "dir",
					Path:    path,
					Message: "Skipping directory",
				}
				return filepath.SkipDir
			}
			eventChan <- WalkEvent{
				Type:    "dir",
				Path:    path,
				Message: "Scanning directory",
			}
		} else {
			eventChan <- WalkEvent{
				Type:    "file",
				Path:    path,
				Message: "Processing file",
			}
		}

		return nil
	})

	if err != nil {
		eventChan <- WalkEvent{
			Type:  "error",
			Error: fmt.Sprintf("Walk failed: %v", err),
		}
	} else {
		eventChan <- WalkEvent{
			Type:    "complete",
			Message: "Directory scan completed",
		}
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
