package jobs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func TestActivateSubtask(t *testing.T) {

	testSubtask := NewSubtask("no-op", "foo")
	ctx := context.Background()
	activated := testSubtask.Activate(ctx)
	assert.True(t, activated)

	activated = testSubtask.Activate(ctx)
	assert.False(t, activated)
}

func TestRunSubtask(t *testing.T) {

	testSubtask := NewSubtask("no-op", "bar")
	ctx := context.Background()
	err := testSubtask.Run(ctx)

	// the action on the subtask is intentionally invalid
	assert.ErrorIs(t, err, errorz.ErrInvalidArgument)
	assert.Equal(t, model.JobStateFailed, testSubtask.Status.State)
}
