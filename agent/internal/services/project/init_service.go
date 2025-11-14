// internal/services/project/init_service.go
package project

import (
	"context"
	"fmt"
	"path/filepath"

	// "strings"
	"time"

	"vcx/agent/internal/domains/project"
	"vcx/agent/internal/services/filters"
	"vcx/agent/internal/services/fs/walk"
	"vcx/agent/internal/services/message"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


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

    // Walk filesystem and ingest files
    eventChan := make(chan walk.Event, 100)
    log.Info("Loading Filter")

    // Create both filters for comparison with same patterns
    quickFilter := filters.FromFile(projectPath, filepath.Join(filters.FILTERSAMPLEPATH, filters.FILTERFILE))
    basicFilter := filters.FromFileBasic(projectPath, filepath.Join(filters.FILTERSAMPLEPATH, filters.FILTERFILE))
    quickFilter.Init()
    basicFilter.Init()
    log.Info("Filters completed")

    // Benchmark comparison
    runFilterBenchmark(projectPath, quickFilter, basicFilter)

    walkStart := time.Now()
    go func() {
        defer close(msgChan)
        log.Info("Walking Starting")
        log.Info(fmt.Sprintf("Project Path: %s", projectPath))
        numProcess := 0
        numSkipped := 0
        numDir     := 0
        go walk.Walk(projectPath, eventChan, quickFilter)

        // Process walk events
        for event := range eventChan {
            // event.Log() // TEMP for debugging
            msgChan <- message.Log(event.Data)
            switch event.Type {
            case walk.FILE:
                numProcess++
                if err := ingestFile(ctx, proj.ID, event.Data); err != nil {
                    log.Error("Failed to ingest file", "path", event.Data, "error", err)
                }
            case walk.ERROR:
                log.Error("Walk error", "error", event.Data)
            case walk.DIR:
                numDir++
                // log.Info("Processing directory", "path", event.Data)
            case walk.SKIP:
                numSkipped++
                // log.Info("Skipped", "path", event.Data)
            }
        }
        log.Info("Walking Completed")
        duration := time.Since(walkStart)
        log.Info(fmt.Sprintf("Walk duration: %s", duration))
        log.Info(fmt.Sprintf("Walk processed: %d", numProcess))
        log.Info(fmt.Sprintf("Walk skipped: %d", numSkipped))
        log.Info(fmt.Sprintf("Walk dirs: %d", numDir))
        log.Info("Project initialization completed", "project", proj.Name)
        log.Info(fmt.Sprintf("Project Path: %s", projectPath))

    }()

    return proj, msgChan
}

func runFilterBenchmark(projectPath string, quickFilter filters.FilterInterface, basicFilter filters.FilterInterface) {
    // Get test paths by walking without filter
    testPaths := []string{}
    testDirs := []bool{}
    eventChan := make(chan walk.Event, 100)

    go walk.Walk(projectPath, eventChan, nil) // No filter to get all paths

    for event := range eventChan {
        if event.Type == walk.FILE || event.Type == walk.DIR {
            testPaths = append(testPaths, event.Data)
            testDirs = append(testDirs, event.Type == walk.DIR)
        }
    }

    if len(testPaths) == 0 {
        log.Info("No paths found for benchmark")
        return
    }

    // Compare filters and collect results
    quickResults := make([]bool, len(testPaths))
    basicResults := make([]bool, len(testPaths))

    // Benchmark QuickFilter
    quickStart := time.Now()
    quickSkipped := 0
    for i, path := range testPaths {
        shouldSkip := quickFilter.ShouldSkip(path, testDirs[i])
        quickResults[i] = shouldSkip
        if shouldSkip {
            quickSkipped++
        }
    }
    quickDuration := time.Since(quickStart)
    quickProcessed := len(testPaths) - quickSkipped

    // Benchmark BasicFilter
    basicStart := time.Now()
    basicSkipped := 0
    for i, path := range testPaths {
        shouldSkip := basicFilter.ShouldSkip(path, testDirs[i])
        basicResults[i] = shouldSkip
        if shouldSkip {
            basicSkipped++
        }
    }
    basicDuration := time.Since(basicStart)
    basicProcessed := len(testPaths) - basicSkipped

    // Check for discrepancies
    discrepancies := 0
    for i, path := range testPaths {
        if quickResults[i] != basicResults[i] {
            discrepancies++
            quickAction := "PROCESS"
            basicAction := "PROCESS"
            if quickResults[i] {
                quickAction = "SKIP"
            }
            if basicResults[i] {
                basicAction = "SKIP"
            }
            log.Info(fmt.Sprintf("DISCREPANCY: %s (isDir: %v) - Quick: %s, Basic: %s",
                path, testDirs[i], quickAction, basicAction))
        }
    }

    // Calculate throughput
    quickThroughput := float64(len(testPaths)) / quickDuration.Seconds()
    basicThroughput := float64(len(testPaths)) / basicDuration.Seconds()

    log.Info("=== FILTER BENCHMARK RESULTS ===")
    log.Info(fmt.Sprintf("QuickFilter: %d skipped, %d processed, %s duration",
        quickSkipped, quickProcessed, quickDuration))
    log.Info(fmt.Sprintf("BasicFilter: %d skipped, %d processed, %s duration",
        basicSkipped, basicProcessed, basicDuration))

    // Correctness check
    if discrepancies > 0 {
        log.Info(fmt.Sprintf("CORRECTNESS ISSUE: %d path discrepancies found", discrepancies))
        log.Info(fmt.Sprintf("QuickFilter skipped %d, BasicFilter skipped %d (diff: %d)",
            quickSkipped, basicSkipped, quickSkipped-basicSkipped))
    } else {
        log.Info("CORRECTNESS: QuickFilter matches BasicFilter exactly")
    }

    if quickDuration < basicDuration {
        speedup := float64(basicDuration) / float64(quickDuration)
        log.Info(fmt.Sprintf("Performance: QuickFilter is %.2fx faster", speedup))
    } else {
        slowdown := float64(quickDuration) / float64(basicDuration)
        log.Info(fmt.Sprintf("Performance: QuickFilter is %.2fx slower", slowdown))
    }

    log.Info(fmt.Sprintf("Throughput: Quick=%.0f paths/sec, Basic=%.0f paths/sec",
        quickThroughput, basicThroughput))
}

func ingestFile(ctx context.Context, projectID, filePath string) error {
    // TODO: Create file records, calculate hashes, etc.
    log.Debug("Ingesting file", "path", filePath)
    return nil
}
