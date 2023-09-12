package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/kmttg-plus/go/pkg/model"
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

func TestPipelineStatus(t *testing.T) {
	assert.True(t, true)
}
