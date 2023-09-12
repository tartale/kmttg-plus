package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"golang.org/x/exp/slices"
)

func TestNewPipeline(t *testing.T) {
	testJob := &model.Job{
		ID:     primitives.Ref("12345"),
		Action: model.JobActionDownload,
		ShowID: "foo",
	}

	p := NewPipeline(testJob)
	assert.Equal(t, "12345", p.jobID)
	assert.Len(t, p.subtasks, 1)
	assert.Equal(t, model.JobActionDownload, p.subtasks[0].Action)

	testJob.Action = model.JobActionEncode
	p = NewPipeline(testJob)

	assert.Len(t, p.subtasks, len(model.AllJobAction))
	assert.Equal(t, model.JobActionDownload, p.subtasks[0].Action)
	assert.Equal(t, model.JobActionEncode, p.subtasks[3].Action)
}

func TestStartPipeline(t *testing.T) {

	testJob := &model.Job{
		ID:     primitives.Ref("12345"),
		Action: model.JobActionDownload,
		ShowID: "foo",
	}
	p := NewPipeline(testJob)
	err := p.Start()
	assert.Nil(t, err)
	assert.Len(t, pipelineQueue, 1)

	task := <-pipelineQueue
	assert.Equal(t, "12345", task.jobID)
}

func TestStartPipeline_MaxInProgress(t *testing.T) {

	pipelineQueue = make(chan *Pipeline, 1)
	testJob := &model.Job{
		ID:     primitives.Ref("12345"),
		Action: model.JobActionDownload,
		ShowID: "foo",
	}
	p := NewPipeline(testJob)
	err := p.Start()
	assert.Nil(t, err)
	assert.Len(t, pipelineQueue, 1)

	testJob.ID = primitives.Ref("67890")
	p = NewPipeline(testJob)
	err = p.Start()
	assert.ErrorIs(t, err, ErrTooManyTasks)
	assert.Len(t, pipelineQueue, 1)

	task := <-pipelineQueue
	assert.Equal(t, "12345", task.jobID)

	err = p.Start()
	assert.Nil(t, err)
	task = <-pipelineQueue
	assert.Equal(t, "67890", task.jobID)
}

func TestRunPipeline(t *testing.T) {
	assert.True(t, true)
}

func TestPipelineStatus_SingleSubtask(t *testing.T) {

	testJob := &model.Job{
		ID:     primitives.Ref("12345"),
		Action: model.JobActionDownload,
		ShowID: "foo",
	}
	p := NewPipeline(testJob)

	status := p.Status()

	assert.Len(t, status.Subtasks, 1)
	assert.Equal(t, model.JobStateQueued, status.Subtasks[0].State)
	assert.Equal(t, 0, status.Subtasks[0].Progress)
	assert.Equal(t, model.JobStateQueued, status.State)
	assert.Equal(t, 0, status.Progress)

	p.subtasks[0].Status.State = model.JobStateRunning
	p.subtasks[0].Status.Progress = 50
	status = p.Status()

	assert.Equal(t, model.JobStateRunning, status.State)
	assert.Equal(t, 50, status.Progress)
}

func TestPipelineStatus_MultipleSubtask(t *testing.T) {

	testAction := model.JobActionDecrypt
	testJob := &model.Job{
		ID:     primitives.Ref("12345"),
		Action: testAction,
		ShowID: "foo",
	}
	testActionNumber := slices.Index(model.AllJobAction, testAction)
	assert.GreaterOrEqual(t, testActionNumber, 0)
	p := NewPipeline(testJob)

	status := p.Status()

	assert.Len(t, status.Subtasks, testActionNumber+1)
	for _, subtask := range status.Subtasks {
		assert.Equal(t, model.JobStateQueued, subtask.State)
	}

	p.subtasks[0].Status.State = model.JobStateRunning
	p.subtasks[0].Status.Progress = 50
	status = p.Status()

	assert.Equal(t, model.JobStateRunning, status.State)
	assert.Equal(t, 25, status.Progress)

	p.subtasks[0].Status.State = model.JobStateComplete
	p.subtasks[0].Status.Progress = 100
	p.subtasks[1].Status.State = model.JobStateRunning
	p.subtasks[1].Status.Progress = 0
	status = p.Status()

	assert.Equal(t, model.JobStateRunning, status.State)
	assert.Equal(t, 50, status.Progress)

	p.subtasks[1].Status.Progress = 50
	status = p.Status()

	assert.Equal(t, model.JobStateRunning, status.State)
	assert.Equal(t, 75, status.Progress)

	p.subtasks[1].Status.State = model.JobStateComplete
	p.subtasks[1].Status.Progress = 100
	status = p.Status()

	assert.Equal(t, model.JobStateComplete, status.State)
	assert.Equal(t, 100, status.Progress)
}
