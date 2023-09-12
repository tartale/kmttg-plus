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
		jobID:    *job.ID,
		action:   job.Action,
		showID:   job.ShowID,
		subtasks: []*Subtask{},
	}

	// the list of job actions is in order of dependency.
	// add all the actions needed up until the requested
	// action; for example, if the job is to "decrypt X",
	// then the subtasks are "download X, decrypt X".
	for _, action := range model.AllJobAction {

		if action == job.Action {
			subtask := NewSubtask(action, job.ShowID)
			pipeline.subtasks = append(pipeline.subtasks, subtask)
			break
		}

		subtask := NewSubtask(action, job.ShowID)
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

	numSubtasks := len(p.subtasks)
	jobStatus := &model.JobStatus{
		JobID:    p.jobID,
		Action:   p.action,
		ShowID:   p.showID,
		Subtasks: []*model.JobSubtaskStatus{},
	}
	for _, subtask := range p.subtasks {
		jobStatus.Subtasks = append(jobStatus.Subtasks, subtask.Status)
	}

	finishedSubtasks := 0
	done := false
	for _, subtask := range p.subtasks {

		switch subtask.Status.State {

		case model.JobStateQueued:
			jobStatus.State = model.JobStateQueued
			done = true

		case model.JobStateRunning:
			jobStatus.State = model.JobStateRunning
			jobStatus.Progress += subtask.Status.Progress / numSubtasks
			done = true

		case model.JobStateComplete:
			finishedSubtasks++
			progress := int(float32(finishedSubtasks) / float32(numSubtasks) * 100)
			jobStatus.Progress = progress
			if progress == 100 {
				jobStatus.State = model.JobStateComplete
				done = true
			}

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
