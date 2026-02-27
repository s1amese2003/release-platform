package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TaskTypeScanRelease = "release:scan"

type ScanTaskPayload struct {
	ReleaseID uint `json:"release_id"`
}

func NewScanTask(releaseID uint, queueName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScanTaskPayload{ReleaseID: releaseID})
	if err != nil {
		return nil, fmt.Errorf("marshal scan task payload: %w", err)
	}
	return asynq.NewTask(TaskTypeScanRelease, payload, asynq.Queue(queueName)), nil
}

type ScanJobRunner interface {
	RunScan(context.Context, uint) error
}

func NewMux(runner ScanJobRunner) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTypeScanRelease, func(ctx context.Context, task *asynq.Task) error {
		var payload ScanTaskPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("decode scan payload: %w", err)
		}
		return runner.RunScan(ctx, payload.ReleaseID)
	})
	return mux
}
