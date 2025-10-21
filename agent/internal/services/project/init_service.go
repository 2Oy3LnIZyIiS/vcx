// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	"vcx/agent/internal/domains/project"
	"vcx/agent/internal/services/fs/walk"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()


type InitService struct {}

func NewInitService() *InitService {
    return &InitService{}
}

func (s *InitService) InitializeProject(ctx context.Context, projectPath string) (*project.Project, error) {
    // Create project record
    // No changeID as this is the first record
    projectName := filepath.Base(projectPath)
    proj, err   := project.New(ctx, projectName, "")
    if err != nil {
        return nil, fmt.Errorf("failed to create project: %w", err)
    }

    // Walk filesystem and ingest files
    eventChan := make(chan walk.WalkEvent, 100)
    filter    := walk.DefaultFilter()

    go walk.Walk(projectPath, eventChan, filter)

    // Process walk events
	// "file", "dir", "skip", "error", "complete"
    for event := range eventChan {
        log.Info("Walk Event", "type", event.Type, "path", event.Path)
        switch event.Type {
        case "file":
            if err := s.ingestFile(ctx, proj.ID, event.Path); err != nil {
                log.Error("Failed to ingest file", "path", event.Path, "error", err)
            }
        case "error":
            log.Error("Walk error", "path", event.Path, "error", event.Error)
        case "complete":
            log.Info("Project initialization completed", "project", proj.Name)
        }
    }

    return proj, nil
}

func (s *InitService) ingestFile(ctx context.Context, projectID, filePath string) error {
    // TODO: Create file records, calculate hashes, etc.
    log.Debug("Ingesting file", "path", filePath)
    return nil
}
