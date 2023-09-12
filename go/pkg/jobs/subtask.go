package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorz"
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

func NewSubtask(action model.JobAction, showID string) *Subtask {

	return &Subtask{
		JobSubtask: &model.JobSubtask{
			ID:     MakeSubtaskID(action, showID),
			Action: action,
			ShowID: showID,
			Status: &model.JobSubtaskStatus{
				Action:   action,
				ShowID:   showID,
				State:    model.JobStateQueued,
				Progress: 0,
			},
		},
	}
}

func (st *Subtask) Activate(ctx context.Context) (activated bool) {

	_, loaded := activeSubtasks.LoadOrCompute(st.ID, func() *Subtask {
		st.Status.State = model.JobStateRunning
		st.ctx = ctx
		return st
	})

	return !loaded
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

	default:
		err = fmt.Errorf("%w: invalid action '%s'", errorz.ErrInvalidArgument, st.Action)
		st.Fail(ctx)
	}

	return err
}

func (st *Subtask) Complete(ctx context.Context) {

	st.Status.State = model.JobStateComplete
	st.Status.Progress = 100
	activeSubtasks.Delete(st.ID)
}

func (st *Subtask) Fail(ctx context.Context) {

	st.Status.State = model.JobStateFailed
	activeSubtasks.Delete(st.ID)
}
