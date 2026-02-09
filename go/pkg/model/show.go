package model

import (
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
)

func NewShow(recordingItem *message.RecordingItem) (Show, error) {

	switch recordingItem.CollectionType {
	case message.CollectionTypeSeries:
		var episodeNumber int
		if len(recordingItem.EpisodeNum) > 0 {
			episodeNumber = recordingItem.EpisodeNum[0]
		}
		return &Episode{
			Kind:               ShowKindEpisode,
			Title:              recordingItem.Title,
			RecordedOn:         recordingItem.StartTime.Time,
			Description:        recordingItem.Description,
			OriginalAirDate:    recordingItem.OriginalAirDate.Time,
			SeasonNumber:       recordingItem.SeasonNumber,
			EpisodeNumber:      episodeNumber,
			EpisodeTitle:       recordingItem.Subtitle,
			EpisodeDescription: recordingItem.Description,
		}, nil
	}

	return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))
}
