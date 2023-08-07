package model

import (
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
)

func NewShow(recordingItem *message.RecordingItem, collectionDetails *message.CollectionItem) (Show, error) {

	switch recordingItem.CollectionType {
	case message.CollectionTypeSeries:
		return NewSeries(recordingItem, collectionDetails)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return NewMovie(recordingItem)

	default:
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))
	}

}

func NewMovie(recordingItem *message.RecordingItem) (*Movie, error) {
	if recordingItem.CollectionType != message.CollectionTypeMovie && recordingItem.CollectionType != message.CollectionTypeSpecial {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))
	}

	return &Movie{
		Kind:        ShowKindMovie,
		RecordingID: recordingItem.RecordingID,
		Title:       recordingItem.Title,
		RecordedOn:  recordingItem.StartTime.Time,
		Description: recordingItem.Description,
		MovieYear:   recordingItem.MovieYear,
	}, nil
}

func NewSeries(recordingItem *message.RecordingItem, collectionItem *message.CollectionItem) (*Series, error) {
	if recordingItem.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))
	}

	return &Series{
		Kind:         ShowKindSeries,
		CollectionID: recordingItem.CollectionID,
		Title:        recordingItem.Title,
		RecordedOn:   recordingItem.StartTime.Time,
		Description:  collectionItem.Description,
	}, nil
}

func NewEpisode(recordingItem *message.RecordingItem) (*Episode, error) {
	if recordingItem.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingItem.CollectionType)))
	}

	var episodeNumber int
	if len(recordingItem.EpisodeNum) > 0 {
		episodeNumber = recordingItem.EpisodeNum[0]
	}
	return &Episode{
		Kind:               ShowKindEpisode,
		RecordingID:        recordingItem.RecordingID,
		Title:              recordingItem.Title,
		RecordedOn:         recordingItem.StartTime.Time,
		Description:        "",
		OriginalAirDate:    recordingItem.OriginalAirDate,
		SeasonNumber:       recordingItem.SeasonNumber,
		EpisodeNumber:      episodeNumber,
		EpisodeTitle:       recordingItem.Subtitle,
		EpisodeDescription: recordingItem.Description,
	}, nil
}

func MergeEpisodes(shows []Show) []Show {

	combinedShowsMap := make(map[string]Show)

	for _, show := range shows {
		if show.GetKind() == ShowKindEpisode {
			episode := show.(*Episode)
			if existingSeries, exists := combinedShowsMap[episode.Title]; exists {
				series := existingSeries.(*Series)
				series.Episodes = append(series.Episodes, episode)
				if series.RecordedOn.Before(episode.RecordedOn) {
					series.RecordedOn = episode.RecordedOn
				}
			} else {
				combinedShowsMap[episode.Title] = &Series{
					Kind:        ShowKindSeries,
					Title:       show.GetTitle(),
					RecordedOn:  show.GetRecordedOn(),
					Description: show.GetDescription(),
					Episodes:    []*Episode{episode},
				}
			}
		} else if show.GetKind() == ShowKindMovie {
			combinedShowsMap[show.GetTitle()] = show
		}
	}

	combinedShows := make([]Show, 0, len(combinedShowsMap))
	for _, show := range combinedShowsMap {
		combinedShows = append(combinedShows, show)
	}

	return combinedShows
}
