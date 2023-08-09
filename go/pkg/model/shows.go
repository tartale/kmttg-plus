package model

import (
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
)

func NewShow(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (Show, error) {

	switch recordingDetails.CollectionType {
	case message.CollectionTypeSeries:
		return NewEpisode(recordingDetails, collectionDetails)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return NewMovie(recordingDetails)

	default:
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

}

func NewMovie(recordingDetails *message.RecordingItem) (*Movie, error) {
	if recordingDetails.CollectionType != message.CollectionTypeMovie && recordingDetails.CollectionType != message.CollectionTypeSpecial {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	return &Movie{
		Kind:        ShowKindMovie,
		RecordingID: recordingDetails.RecordingID,
		Title:       recordingDetails.Title,
		RecordedOn:  recordingDetails.StartTime.Time,
		Description: recordingDetails.Description,
		MovieYear:   recordingDetails.MovieYear,
	}, nil
}

func NewSeries(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*Series, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	return &Series{
		Kind:         ShowKindEpisode,
		CollectionID: recordingDetails.CollectionID,
		Title:        recordingDetails.Title,
		RecordedOn:   recordingDetails.StartTime.Time,
		Description:  collectionDetails.Description,
	}, nil
}

func NewEpisode(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*Episode, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected collection type: %s", string(recordingDetails.CollectionType)))
	}

	var episodeNumber int
	if len(recordingDetails.EpisodeNum) > 0 {
		episodeNumber = recordingDetails.EpisodeNum[0]
	}
	return &Episode{
		Kind:               ShowKindEpisode,
		RecordingID:        recordingDetails.RecordingID,
		CollectionID:       collectionDetails.CollectionID,
		Title:              recordingDetails.Title,
		RecordedOn:         recordingDetails.StartTime.Time,
		Description:        collectionDetails.Description,
		OriginalAirDate:    recordingDetails.OriginalAirDate,
		SeasonNumber:       recordingDetails.SeasonNumber,
		EpisodeNumber:      episodeNumber,
		EpisodeTitle:       recordingDetails.Subtitle,
		EpisodeDescription: recordingDetails.Description,
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
					Kind:         ShowKindSeries,
					CollectionID: episode.CollectionID,
					Title:        episode.Title,
					RecordedOn:   episode.RecordedOn,
					Description:  episode.Description,
					Episodes:     []*Episode{episode},
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
