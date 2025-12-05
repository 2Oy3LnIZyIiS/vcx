package walk

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"vcx/agent/internal/services/filters"
)


type Snapshot struct {
    Files    []string
    Dirs     []string
    Symlinks []string
}

func GetSnapshot(dirPath string, filter ...filters.FilterInterface) Snapshot {
    var snap Snapshot
    var f filters.FilterInterface
    if len(filter) > 0 {
        f = filter[0]
    }

    filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }

        isDir := d.IsDir()
        if f != nil && f.ShouldSkip(path, isDir) {
            if isDir {
                return filepath.SkipDir
            }
            return nil
        }

        if d.Type()&fs.ModeSymlink != 0 {
            snap.Symlinks = append(snap.Symlinks, path)
        } else if isDir {
            snap.Dirs = append(snap.Dirs, path)
        } else {
            snap.Files = append(snap.Files, path)
        }
        return nil
    })

    return snap
}


// Stream streams directory traversal events to the provided channel
func Stream(dirPath string, eventChan chan<- Event, filter filters.FilterInterface, verbose ...bool) {
	defer close(eventChan)

	isVerbose := len(verbose) > 0 && verbose[0]

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			eventChan <- Error(fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		isDir := d.IsDir()
		if filter != nil && filter.ShouldSkip(path, isDir) {
			if isVerbose {
				eventChan <- Skip(path)
			}
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type()&fs.ModeSymlink != 0 {
			eventChan <- Sym(path)
		} else if isDir {
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
