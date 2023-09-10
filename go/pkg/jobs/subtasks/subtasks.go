package subtasks

import (
	"fmt"
	"path"
	"strings"

	"github.com/puzpuzpuz/xsync"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

type Subtask interface {
	GetID() string
	GetAction() model.JobAction
}

var (
	subtasks = xsync.NewMapOf[model.JobSubtask]()
)

func New(action model.JobAction, showID string) Subtask {

	id := createID(action, showID)
	subtask, _ := subtasks.LoadOrStore(id, model.JobSubtask{
		ID:     id,
		Action: action,
		ShowID: showID,
		Progress: &model.JobSubtaskProgress{
			Status:   model.JobStatusQueued,
			Progress: 0,
		},
	})

	return subtask
}

func Download(subtaskID string) {

	downloadDir := path.Join(string(config.Values.OutputDir), subtaskID)
	fmt.Println(downloadDir)

}

func Run(subtaskID string) error {

	action, _ := parseID(subtaskID)
	switch action {
	case model.JobActionDownload:
		Download(subtaskID)

	}

	return nil
}

func createID(action model.JobAction, showID string) string {

	return fmt.Sprintf("%s/%s", strings.ToLower(string(action)), showID)
}

func parseID(subtaskID string) (action model.JobAction, showID string) {

	split := strings.SplitN(subtaskID, "/", 1)
	action = model.JobAction(strings.ToUpper(split[0]))
	showID = split[1]

	return
}
