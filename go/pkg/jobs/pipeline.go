package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var (
	pipelineQueue   = make(chan *Pipeline, config.Values.MaxBackgroundTasks)
	jobDependencies = map[model.JobAction][]model.JobAction{
		model.JobActionDownload: {},
		model.JobActionComskip:  {model.JobActionDownload},
		model.JobActionEncode:   {model.JobActionDownload},
		model.JobActionPlay:     {model.JobActionDownload},
	}

	ErrTooManyTasks = errors.New("too many tasks; try again later")
)

type Pipeline struct {
	jobID    string
	action   model.JobAction
	show     model.Show
	subtasks []*Subtask
}

func NewPipeline(job *model.Job, show model.Show) *Pipeline {

	pipeline := Pipeline{
		jobID:    *job.ID,
		action:   job.Action,
		show:     show,
		subtasks: []*Subtask{},
	}

	dependencies := jobDependencies[job.Action]
	for _, action := range dependencies {
		subtask := NewSubtask(action, show)
		pipeline.subtasks = append(pipeline.subtasks, subtask)
	}
	subtask := NewSubtask(job.Action, show)
	pipeline.subtasks = append(pipeline.subtasks, subtask)

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

func (p *Pipeline) Run(ctx context.Context) {

	for _, subtask := range p.subtasks {
		err := subtask.Run(ctx)
		if err != nil {
			return
		}
	}
}

func (p *Pipeline) Status() *model.JobStatus {

	numSubtasks := len(p.subtasks)
	jobStatus := &model.JobStatus{
		JobID:    p.jobID,
		Action:   p.action,
		ShowID:   p.show.GetID(),
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
