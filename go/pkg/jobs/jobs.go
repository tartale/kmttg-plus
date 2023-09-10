package jobs

import (
	"github.com/tartale/kmttg-plus/go/pkg/jobs/subtasks"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

type pipeline map[model.JobAction]subtasks.Subtask

func (p pipeline) Add(st subtasks.Subtask) {

	p[st.GetAction()] = st
}

func (p pipeline) Run() error {

	for _, action := range model.AllJobAction {
		if st, ok := p[action]; ok {
			err := subtasks.Run(st.GetID())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Start(job *model.Job) *model.JobStatus {

	var pipe pipeline

	for _, action := range model.AllJobAction {

		if action == job.Action {
			pipe.Add(subtasks.New(action, job.ShowID))
			break
		}

		pipe.Add(subtasks.New(action, job.ShowID))
	}

	pipe.Run()

	return nil
}
