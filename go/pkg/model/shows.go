package model

import (
	"fmt"
	"sort"

	"github.com/PaesslerAG/gval"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"go.uber.org/zap"
)

type ShowFilterFn = func(s Show) bool

func NewShowFilter(sf *ShowFilter) ShowFilterFn {

	return func(s Show) bool {

		expression := filter.GetExpression(s)
		values := filter.GetValues(sf, s)
		eval, err := gval.Evaluate(expression, values)
		if err != nil {
			logz.Logger.Warn("error attempting to filter shows", zap.Error(err))
			return true
		}

		return eval.(bool)
	}
}

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
	sort.Slice(combinedShows, func(i, j int) bool {
		return combinedShows[i].GetTitle() < combinedShows[j].GetTitle()
	})

	return combinedShows
}
