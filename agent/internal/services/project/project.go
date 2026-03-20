// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	projectDomain "vcx/agent/internal/domains/project"
	branchService "vcx/agent/internal/services/branch"
	changeService "vcx/agent/internal/services/change"
	fileService "vcx/agent/internal/services/file"
	"vcx/agent/internal/services/filters"
	"vcx/agent/internal/services/fs/walk"
	instanceService "vcx/agent/internal/services/instance"
	tagService "vcx/agent/internal/services/tag"
	"vcx/agent/internal/session"
	"vcx/pkg/logging"
	"vcx/pkg/message"
)


var log = logging.GetLogger()


// POC
func NewProject(ctx context.Context, projectPath string) (*projectDomain.Project, <-chan message.Event) {
	msgChan := make(chan message.Event)

	project, ctx, err := initializeProjectEntities(ctx, projectPath)
	if err != nil {
		msgChan <- message.Error(fmt.Errorf("failed to create project: %w", err).Error())
		close(msgChan)
		return nil, msgChan
	}

	go ingestProjectFiles(ctx, projectPath, project, msgChan)

	return project, msgChan
}


func initializeProjectEntities(ctx context.Context, projectPath string) (*projectDomain.Project, context.Context, error) {
	// new change
	change, err := changeService.CreateProjectChange(ctx)
	if err != nil {
		return nil, ctx, err
	}
	ctx = session.WithChangeID(ctx, change.ID)

	// new project
	project, err := projectDomain.New(ctx, filepath.Base(projectPath))
	if err != nil {
		return nil, ctx, err
	}
	ctx = session.WithProjectID(ctx, project.ID)

	// New branch
	branch, err := branchService.Create(ctx, "main")
	if err != nil {
		return nil, ctx, err
	}
	ctx = session.WithBranchID(ctx, branch.ID)

	// New instance
	_, err = instanceService.Create(ctx, projectPath)
	if err != nil {
		return nil, ctx, err
	}

	// New Tag
	_, err = tagService.CreateSystemProjectTag(ctx)
	if err != nil {
		return nil, ctx, err
	}

	return project, ctx, nil
}


func ingestProjectFiles(ctx context.Context, projectPath string, project *projectDomain.Project, msgChan chan message.Event) {
	defer close(msgChan)

	quickFilter := filters.FromFile(projectPath, filepath.Join(filters.FILTERSAMPLEPATH, filters.FILTERFILE))
	quickFilter.Init()

	log.Info("Walking Starting")
	log.Info(fmt.Sprintf("Project Path: %s", projectPath))

	numProcess := 0
	eventChan  := make(chan walk.Event, 100)
	go walk.Stream(projectPath, eventChan, quickFilter)

	// Process walk events
	for event := range eventChan {
		msgChan <- message.Log(event.Data)
		switch event.Type {
		case walk.ERROR:
			log.Error("Walk error", "error", event.Data)
		case walk.FILE:
			numProcess++
			if _, err := fileService.Ingest(ctx, projectPath, event.Data); err != nil {
				logIngestionFailure("file", event, err)
			}
		case walk.SYM:
			numProcess++
			if err := ingestSymlink(ctx, projectPath, event.Data); err != nil {
				logIngestionFailure("symlink", event, err)
			}
		}
	}

	log.Info(fmt.Sprintf("Walk processed: %d", numProcess))
	log.Info("Project initialization completed", "project", project.Name)
	log.Info(fmt.Sprintf("Project Path: %s", projectPath))
}


func logIngestionFailure(pathType string, event walk.Event, err error) {
    log.Error( fmt.Sprintf("Failed to ingest %s", pathType),
               "path",  event.Data,
               "error", err )
}


// func ingestDir(ctx context.Context, filePath string) error {
//     // TODO: Create file records, calculate hashes, etc.
//     log.Debug("Ingesting dir", "path", filePath)
//     return nil
// }


func ingestSymlink(ctx context.Context, projectPath, linkPath string) error {
    if _, err := fileService.IngestSymlink(ctx, projectPath, linkPath); err != nil {
        return err
    }
    return nil
}
