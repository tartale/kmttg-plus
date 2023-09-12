package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var (
	pipelineQueue = make(chan *Pipeline, config.Values.MaxBackgroundTasks)

	ErrTooManyTasks = errors.New("too many tasks; try again later")
)

type Pipeline struct {
	jobID    string
	action   model.JobAction
	showID   string
	subtasks []*Subtask
}

func NewPipeline(job *model.Job) *Pipeline {

	pipeline := Pipeline{
		jobID: *job.ID,
	}

	for _, action := range model.AllJobAction {

		if action == job.Action {
			subtask := NewSubtask(job)
			pipeline.subtasks = append(pipeline.subtasks, subtask)
			break
		}

		subtask := NewSubtask(job)
		pipeline.subtasks = append(pipeline.subtasks, subtask)
	}

	return &pipeline
}

func (p *Pipeline) Start() error {

	select {
	case pipelineQueue <- p:
		break
	default:
		return ErrTooManyTasks
	}

	return nil
}

func (p *Pipeline) Run(ctx context.Context) error {

	for _, subtask := range p.subtasks {
		err := subtask.Run(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Pipeline) Status() *model.JobStatus {

	jobStatus := &model.JobStatus{
		JobID:  p.jobID,
		Action: p.action,
		ShowID: p.showID,
	}
	numSubtasks := len(p.subtasks)
	finishedSubtasks := 0
	done := false
	for _, subtask := range p.subtasks {

		jobStatus.Subtasks = append(jobStatus.Subtasks, subtask.Status)
		switch subtask.Status.State {

		case model.JobStateQueued:
			jobStatus.State = model.JobStateQueued
			done = true

		case model.JobStateRunning:
			jobStatus.State = model.JobStateRunning
			jobStatus.Progress += subtask.Status.Progress / numSubtasks * 100
			done = true

		case model.JobStateComplete:
			finishedSubtasks++
			jobStatus.Progress = finishedSubtasks / numSubtasks * 100
			done = false

		case model.JobStateFailed:
			jobStatus.State = model.JobStateFailed
			done = true

		default:
			panic(fmt.Errorf("unexpected subtask state: %s", subtask.Status.State))
		}
		if done {
			break
		}
	}

	return jobStatus
}
