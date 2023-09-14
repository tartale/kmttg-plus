package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var activeSubtasks = xsync.NewMapOf[*Subtask]()

type Subtask struct {
	*model.JobSubtask
	id        string
	show      model.Show
	activated bool
	ctx       context.Context
}

func MakeSubtaskID(action model.JobAction, showID string) string {

	return fmt.Sprintf("%s/%s", strings.ToLower(string(action)), showID)
}

func NewSubtask(action model.JobAction, show model.Show) *Subtask {

	return &Subtask{
		JobSubtask: &model.JobSubtask{
			Action: action,
			ShowID: show.GetID(),
			Status: &model.JobSubtaskStatus{
				Action:   action,
				ShowID:   show.GetID(),
				State:    model.JobStateQueued,
				Progress: 0,
			},
		},
		id:   uuid.NewString(),
		show: show,
	}
}

func (st *Subtask) GetID() string {

	return MakeSubtaskID(st.Action, st.ShowID)
}

func (st *Subtask) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancelCause(ctx)
	taskWasStarted := st.activate(ctx)
	if !taskWasStarted {
		// wait for the already-running task
		<-st.ctx.Done()
		return st.ctx.Err()
	}

	var err error
	defer func() { cancel(err) }()
	defer activeSubtasks.Delete(st.GetID())

	switch st.Action {
	case model.JobActionDownload:
		err = Download(ctx, st)
	default:
		err = fmt.Errorf("%w: invalid action '%s'", errorz.ErrInvalidArgument, st.Action)
	}

	if err != nil {
		st.Status.State = model.JobStateFailed
		st.Status.Error = primitives.Ref(err.Error())
	} else {
		st.Status.State = model.JobStateComplete
	}

	return err
}

func (st *Subtask) activate(ctx context.Context) (taskWasStarted bool) {

	if st.activated {
		panic(fmt.Errorf("%w: subtask was activated multiple times", errorz.ErrFatal))
	}
	existingSubtask, loaded := activeSubtasks.LoadOrStore(st.GetID(), st)
	if !loaded {
		st.ctx = ctx
		st.Status.State = model.JobStateRunning
	} else {
		st.ctx = existingSubtask.ctx
		st.Status = existingSubtask.Status
	}
	st.activated = true

	return !loaded
}
