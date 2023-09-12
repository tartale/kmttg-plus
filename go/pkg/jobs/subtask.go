package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var activeSubtasks = xsync.NewMapOf[*Subtask]()

type Subtask struct {
	*model.JobSubtask
	ctx context.Context
}

func MakeSubtaskID(action model.JobAction, showID string) string {

	return fmt.Sprintf("%s/%s", strings.ToLower(string(action)), showID)
}

func NewSubtask(job *model.Job) *Subtask {

	return &Subtask{
		JobSubtask: &model.JobSubtask{
			ID:     MakeSubtaskID(job.Action, job.ShowID),
			Action: job.Action,
			ShowID: job.ShowID,
			Status: &model.JobSubtaskStatus{
				Action:   job.Action,
				ShowID:   job.ShowID,
				State:    model.JobStateQueued,
				Progress: 0,
			},
		},
	}
}

func (st *Subtask) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancelCause(ctx)
	isMine := st.Activate(ctx)
	if !isMine {
		<-st.ctx.Done()
		return st.ctx.Err()
	}

	var err error
	defer func() { cancel(err) }()

	switch st.Action {
	case model.JobActionDownload:
		err = Download(ctx, st)
	}

	return err
}

func (st *Subtask) Activate(ctx context.Context) (activated bool) {

	_, loaded := activeSubtasks.LoadOrCompute(st.ID, func() *Subtask {
		st.Status.State = model.JobStateRunning
		st.ctx = ctx
		return st
	})

	return loaded
}

func (st *Subtask) Complete(ctx context.Context) error {

	st.Status.State = model.JobStateComplete
	st.Status.Progress = 100
	activeSubtasks.Delete(st.ID)

	return nil
}
