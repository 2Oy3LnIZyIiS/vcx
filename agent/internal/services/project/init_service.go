// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	// "strings"

	"vcx/agent/internal/domains/project"
	"vcx/agent/internal/services/filters"
	"vcx/agent/internal/services/fs/walk"
	"vcx/agent/internal/services/message"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// POC
func InitializeProject( ctx context.Context, projectPath string,
                      ) (*project.Project, <-chan message.Event) {
    // Create project record
    // called from agent/internal/infra/http/api/project/project.go::initProject()
    // No changeID as this is the first record
    msgChan     := make(chan message.Event)
    projectName := filepath.Base(projectPath)
    proj, err   := project.New(ctx, projectName, "")
    if err != nil {
        msgChan <- message.Error(fmt.Errorf("failed to create project: %w", err).Error())
        return nil, msgChan
    }

    eventChan := make(chan walk.Event, 100)

    // Create both filters for comparison with same patterns
    quickFilter := filters.FromFile(projectPath, filepath.Join(filters.FILTERSAMPLEPATH, filters.FILTERFILE))
    quickFilter.Init()

    go func() {
        defer close(msgChan)
        log.Info("Walking Starting")
        log.Info(fmt.Sprintf("Project Path: %s", projectPath))
        numProcess := 0
        go walk.Stream(projectPath, eventChan, quickFilter)

        // Process walk events
        for event := range eventChan {
            msgChan <- message.Log(event.Data)
            switch event.Type {
                case walk.ERROR:
                    log.Error("Walk error", "error", event.Data)
                case walk.FILE:
                    numProcess++
                    if err := ingestFile(ctx, proj.ID, event.Data); err != nil {
                        logIngestionFailure("file", event, err)
                    }
                case walk.DIR:
                    if err := ingestDir(ctx, proj.ID, event.Data); err != nil {
                        logIngestionFailure("directory", event, err)
                    }
                case walk.SYM:
                    numProcess++
                    if err := ingestSymlink(ctx, proj.ID, event.Data); err != nil {
                        logIngestionFailure("symlink", event, err)
                    }
            // case walk.SKIP:
            //     // log.Info("Skipped", "path", event.Data)
            }
        }

        log.Info(fmt.Sprintf("Walk processed: %d", numProcess))
        log.Info("Project initialization completed", "project", proj.Name)
        log.Info(fmt.Sprintf("Project Path: %s", projectPath))
    }()

    return proj, msgChan
}


func logIngestionFailure(pathType string, event walk.Event, err error) {
    log.Error( fmt.Sprintf("Failed to ingest %s", pathType),
               "path",  event.Data,
               "error", err )
}


func ingestFile(ctx context.Context, projectID, filePath string) error {
    // TODO: Create file records, calculate hashes, etc.
    log.Debug("Ingesting file", "path", filePath)
    return nil
}


func ingestDir(ctx context.Context, projectID, filePath string) error {
    // TODO: Create file records, calculate hashes, etc.
    log.Debug("Ingesting dir", "path", filePath)
    return nil
}


func ingestSymlink(ctx context.Context, projectID, linkPath string) error {
    // TODO: Store symlink metadata (path and target), don't hash target content
    log.Debug("Ingesting symlink", "path", linkPath)
    return nil
}
