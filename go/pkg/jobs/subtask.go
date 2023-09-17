package jobs

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
	"go.uber.org/zap"
)

var activeSubtasks = xsync.NewMapOf[*Subtask]()

type Subtask struct {
	*model.JobSubtask
	id        string
	show      model.Show
	tmpdir    string
	outputdir string
	activated bool
	ctx       context.Context
}

func MakeSubtaskID(action model.JobAction, showID string) string {

	return fmt.Sprintf("%s/%s", strings.ToLower(string(action)), showID)
}

func NewSubtask(action model.JobAction, show model.Show) *Subtask {

	tmpdir := path.Join(config.Values.TempDir, strings.ToLower(string(action)), shows.GetPath(show))
	outputdir := path.Join(config.Values.OutputDir, strings.ToLower(string(action)), shows.GetPath(show))

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
		id:        uuid.NewString(),
		show:      show,
		tmpdir:    tmpdir,
		outputdir: outputdir,
	}
}

func (st *Subtask) GetID() string {

	return MakeSubtaskID(st.Action, st.ShowID)
}

func (st *Subtask) Run(ctx context.Context) error {

	ctx, cancel := context.WithCancelCause(ctx)
	taskAlreadyStarted := st.activate(ctx)
	if taskAlreadyStarted {
		// wait for the already-running task
		<-st.ctx.Done()
		return st.ctx.Err()
	}

	var err error
	defer func() { cancel(err) }()
	defer activeSubtasks.Delete(st.GetID())

	logz.Logger.Info("starting background task for show",
		zap.String("task", st.Action.String()), zap.String("title", st.show.GetTitle()))
	switch st.Action {

	case model.JobActionDownload:
		err = Download(ctx, st)

	case model.JobActionComskip:
		err = Comskip(ctx, st)

	case model.JobActionEncode:
		err = Encode(ctx, st)

	case model.JobActionPlay:
		err = Play(ctx, st)

	default:
		err = fmt.Errorf("%w: invalid action '%s'", errorz.ErrInvalidArgument, st.Action)
	}
	logz.Logger.Info("finished background task for show",
		zap.String("task", st.Action.String()), zap.String("title", st.show.GetTitle()), zap.Error(err))

	if err != nil {
		st.Status.State = model.JobStateFailed
		st.Status.Error = primitives.Ref(err.Error())
	} else {
		st.Status.State = model.JobStateComplete
	}

	return err
}

func (st *Subtask) activate(ctx context.Context) (taskAlreadyStarted bool) {

	if st.activated {
		panic(fmt.Errorf("%w: subtask was activated multiple times", errorz.ErrFatal))
	}
	existingSubtask, taskAlreadyStarted := activeSubtasks.LoadOrStore(st.GetID(), st)
	if !taskAlreadyStarted {
		st.ctx = ctx
		st.Status.State = model.JobStateRunning
	} else {
		st.ctx = existingSubtask.ctx
		st.Status = existingSubtask.Status
	}
	st.activated = true

	return
}

type ProgressWriter struct {
	subtask      *Subtask
	currentBytes int64
	totalBytes   int64
}

func NewProgressWriter(subtask *Subtask, totalBytes int64) *ProgressWriter {

	return &ProgressWriter{
		subtask:    subtask,
		totalBytes: totalBytes,
	}
}

func (p *ProgressWriter) Write(b []byte) (n int, err error) {

	p.currentBytes += int64(len(b))
	progress := mathx.DivideAndRound(p.currentBytes*100, p.totalBytes)
	p.subtask.Status.Progress = int(mathx.Min(progress, int64(100)))

	return len(b), nil
}
