package model

import (
	"path"
	"strings"

	"github.com/tartale/kmttg-plus/go/pkg/config"
)

func (s *JobSubtask) Tmpdir() string {

	return path.Join(string(config.Values.OutputDir), "tmp", strings.ToLower(string(s.Action)), s.ShowID)
}

func (s *JobSubtask) OutputDir() string {

	return path.Join(string(config.Values.OutputDir), strings.ToLower(string(s.Action)), s.ShowID)
}
