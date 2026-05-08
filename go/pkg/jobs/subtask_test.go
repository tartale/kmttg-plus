package jobs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func TestActivateSubtask(t *testing.T) {
	_, testShow := NewTestData()
	testSubtask := NewSubtask("no-op", testShow)
	ctx := context.Background()
	taskAlreadyStarted := testSubtask.activate(ctx)
	assert.False(t, taskAlreadyStarted)
	assert.Panics(t, func() { testSubtask.activate(ctx) })
}

func TestRunSubtask(t *testing.T) {
	_, testShow := NewTestData()
	testSubtask := NewSubtask("foobar", testShow)
	ctx := context.Background()
	err := testSubtask.Run(ctx)

	// the action on the subtask is intentionally invalid
	assert.ErrorIs(t, err, errorz.ErrInvalidArgument)
	assert.Equal(t, model.JobStateFailed, testSubtask.Status.State)
}
