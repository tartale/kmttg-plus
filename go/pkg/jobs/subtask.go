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

func (st *Subtask) Activate(ctx context.Context) (taskWasActivated bool) {

	existingSubtask, loaded := activeSubtasks.LoadOrStore(st.ID, st)
	if !loaded {
		st.ctx = ctx
		st.Status.State = model.JobStateRunning
	} else {
		st.ctx = existingSubtask.ctx
		st.Status = existingSubtask.Status
	}

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
	}

	if err != nil {
		st.Status.State = model.JobStateFailed
	} else {
		st.Status.State = model.JobStateComplete
	}

	return err
}
