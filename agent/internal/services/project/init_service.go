// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	"vcx/agent/internal/domains/project"
	"vcx/agent/internal/services/fs/walk"
	"vcx/agent/internal/services/message"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


func InitializeProject( ctx context.Context, projectPath string,
                      ) (*project.Project, <-chan message.Event) {
    // Create project record
    // No changeID as this is the first record
    msgChan     := make(chan message.Event)
    projectName := filepath.Base(projectPath)
    proj, err   := project.New(ctx, projectName, "")
    if err != nil {
        msgChan <- message.Error(fmt.Errorf("failed to create project: %w", err).Error())
        return nil, msgChan
    }

    // Walk filesystem and ingest files
    eventChan := make(chan walk.WalkEvent, 100)
    filter    := walk.DefaultFilter()

    go func() {
        defer close(msgChan)
        go walk.Walk(projectPath, eventChan, filter)

        // Process walk events
        // "file", "dir", "skip", "error", "complete"
        for event := range eventChan {
            log.Info("Walk Event", "type", event.Type, "path", event.Path)
            msgChan <- message.Log(event.Path)
            switch event.Type {
            case "file":
                if err := ingestFile(ctx, proj.ID, event.Path); err != nil {
                    log.Error("Failed to ingest file", "path", event.Path, "error", err)
                }
            case "error":
                log.Error("Walk error", "path", event.Path, "error", event.Error)
            case "complete":
                log.Info("Project initialization completed", "project", proj.Name)
            }
        }
    }()

    return proj, msgChan
}

func ingestFile(ctx context.Context, projectID, filePath string) error {
    // TODO: Create file records, calculate hashes, etc.
    log.Debug("Ingesting file", "path", filePath)
    return nil
}
