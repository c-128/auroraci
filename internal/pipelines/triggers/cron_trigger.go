package triggers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/c-128/auroraci/internal/cron"
	"github.com/c-128/auroraci/internal/pipelines"
	"github.com/c-128/auroraci/internal/pipelines/artifacts"
	"github.com/c-128/auroraci/internal/pipelines/logger"
	"github.com/c-128/auroraci/internal/pipelines/runner"
	"github.com/docker/docker/client"
)

func CronTriggerJob(
	logger logger.Logger,
	uploaderProvider artifacts.UploaderProvider,
	docker *client.Client,
) {
	nextTime := time.Now().Truncate(time.Minute)
	for {
		nextTime = nextTime.Add(time.Minute)
		go cronTrigger(logger, uploaderProvider, docker)

		time.Sleep(time.Until(nextTime))
	}
}

func cronTrigger(
	logger logger.Logger,
	uploaderProvider artifacts.UploaderProvider,
	docker *client.Client,
) {
	startTime := time.Now()

	for pipelineID, pipeline := range pipelines.GetPipelines() {
		if pipeline.Build == nil {
			continue
		}

		shouldTrigger := shouldTriggerCron(pipeline.Build, startTime)
		if !shouldTrigger {
			continue
		}

		runID := strconv.FormatInt(time.Now().Unix(), 10)
		uploader, err := uploaderProvider(pipelineID, runID)
		if err != nil {
			log.Printf("Failed to create uploader for pipeline \"%s\": %s", pipelineID, err)
			continue
		}

		log.Printf("Triggering cron job for pipeline \"%s\"", pipelineID)
		go runner.RunPipeline(
			context.Background(),
			logger,
			uploader,
			docker,
			pipeline,
		)
	}
}

func shouldTriggerCron(build *pipelines.Build, time time.Time) bool {
	for _, trigger := range build.Triggers {
		if trigger.Cron == "" {
			continue
		}

		isTime, err := cron.IsExpressionTime(trigger.Cron, time)
		if err != nil {
			continue
		}

		if isTime {
			return true
		}
	}

	return false
}
