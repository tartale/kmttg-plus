package jobs

import (
	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/go/pkg/primitives"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

func NewTestData() (testJob *model.Job, testShow model.Show) {

	jobs = xsync.NewMapOf[*Job]()
	activeSubtasks = xsync.NewMapOf[*Subtask]()
	return &model.Job{
			ID:     primitives.Ref("12345"),
			Action: model.AllJobAction[0],
			ShowID: "foo",
		},
		&model.Movie{
			ID:    "foo",
			Kind:  model.ShowKindMovie,
			Title: "Back to the Future",
		}

}
