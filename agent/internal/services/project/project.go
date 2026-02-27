// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	// "strings"

	branchDomain "vcx/agent/internal/domains/branch"
	changeDomain "vcx/agent/internal/domains/change"
	instanceDomain "vcx/agent/internal/domains/instance"
	projectDomain "vcx/agent/internal/domains/project"
	tagDomain "vcx/agent/internal/domains/tag"
	fileService "vcx/agent/internal/services/file"
	"vcx/agent/internal/services/filters"
	"vcx/agent/internal/services/fs/walk"
	"vcx/agent/internal/services/message"
	"vcx/agent/internal/session"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// POC
func NewProject( ctx context.Context, projectPath string,
                      ) (*projectDomain.Project, <-chan message.Event) {
    // Create project record
    // called from agent/internal/infra/http/api/project/project.go::initProject()
    msgChan      := make(chan message.Event)
    // accountID := session.GetAccountID(ctx)
    // new change
    change, err  := changeDomain.NewProject(ctx)
    ctx           = session.WithChangeID(ctx, change.ID)
    // new project
    project, err := projectDomain.New(ctx, filepath.Base(projectPath))
    ctx           = session.WithProjectID(ctx, project.ID)
    // New branch
    branch, err  := branchDomain.New(ctx, "main")
    ctx           = session.WithBranchID(ctx, branch.ID)
    // New instance
    _, err        = instanceDomain.New(ctx, projectPath)
    // New Tag
    _, err        = tagDomain.NewSystemProjectTag(ctx)
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
                    if _, err := fileService.Ingest(ctx, event.Data); err != nil {
                        logIngestionFailure("file", event, err)
                    }
                // case walk.DIR:
                //     // ? Ingest empty directories ?
                //     // if err := ingestDir(ctx, event.Data); err != nil {
                //     //     logIngestionFailure("directory", event, err)
                //     }
                case walk.SYM:
                    numProcess++
                    if err := ingestSymlink(ctx, event.Data); err != nil {
                        logIngestionFailure("symlink", event, err)
                    }
            // case walk.SKIP:
            //     // log.Info("Skipped", "path", event.Data)
            }
        }

        log.Info(fmt.Sprintf("Walk processed: %d", numProcess))
        log.Info("Project initialization completed", "project", project.Name)
        log.Info(fmt.Sprintf("Project Path: %s", projectPath))
    }()

    return project, msgChan
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


func ingestSymlink(ctx context.Context, linkPath string) error {
    // TODO: Store symlink metadata (path and target), don't hash target content
    log.Debug("Ingesting symlink", "path", linkPath)
    return nil
}
