package shows

import (
	"fmt"
	"sort"

	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

func New(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (model.Show, error) {

	switch recordingDetails.CollectionType {
	case message.CollectionTypeSeries:
		return newEpisode(recordingDetails, collectionDetails)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		return newMovie(recordingDetails, collectionDetails)

	default:
		return nil, fmt.Errorf("%w: unexpected collection type: %s",
			errorz.ErrResponse, recordingDetails.CollectionType)
	}

}

func WithImageURL(show model.Show, width, height int) model.Show {

	if width == 0 || height == 0 {
		return show
	}

	var result model.Show

	switch show.GetKind() {
	case model.ShowKindMovie:
		movie := *(show.(*movie))
		images := movie.collectionDetails.Images
		bestImage := findBestImage(images, width, height)
		movie.ImageURL = bestImage.ImageURL
		result = &movie

	case model.ShowKindSeries:
		series := *(show.(*series))
		images := series.collectionDetails.Images
		bestImage := findBestImage(images, width, height)
		series.ImageURL = bestImage.ImageURL
		result = &series

	case model.ShowKindEpisode:
		return show
	}

	return result
}

func AsAPIType(show model.Show) model.Show {
	switch show.GetKind() {
	case model.ShowKindMovie:
		return show.(*movie).Movie
	case model.ShowKindSeries:
		return show.(*series).Series
	case model.ShowKindEpisode:
		return show.(*episode).Episode
	}
	logz.Logger.Warn("unable to cast show to API type",
		zap.Any("kind", show.GetKind()), zap.String("showTitle", show.GetTitle()))

	return show
}

func MergeEpisodes(shows []model.Show) []model.Show {

	combinedShowsMap := make(map[string]model.Show)

	for _, show := range shows {
		if show.GetKind() == model.ShowKindEpisode {
			if existingSeries, exists := combinedShowsMap[show.GetTitle()]; exists {
				episode := show.(*episode).Episode
				series := existingSeries.(*series)
				series.Episodes = append(series.Episodes, episode)
				if series.RecordedOn.Before(episode.RecordedOn) {
					series.RecordedOn = episode.RecordedOn
				}
			} else {
				episode := show.(*episode)
				combinedShowsMap[show.GetTitle()] = newSeries(episode)
			}
		} else if show.GetKind() == model.ShowKindMovie {
			combinedShowsMap[show.GetTitle()] = show
		}
	}

	combinedShows := make([]model.Show, 0, len(combinedShowsMap))
	for _, show := range combinedShowsMap {
		combinedShows = append(combinedShows, show)
	}
	sort.Slice(combinedShows, func(i, j int) bool {
		return combinedShows[i].GetTitle() < combinedShows[j].GetTitle()
	})

	return combinedShows
}

type movie struct {
	*model.Movie
	recordingDetails  *message.RecordingItem
	collectionDetails *message.CollectionItem
}

func newMovie(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*movie, error) {
	if recordingDetails.CollectionType != message.CollectionTypeMovie && recordingDetails.CollectionType != message.CollectionTypeSpecial {
		return nil, fmt.Errorf("%w: unexpected collection type: %s",
			errorz.ErrResponse, recordingDetails.CollectionType)
	}

	return &movie{
		Movie: &model.Movie{
			ID:          recordingDetails.RecordingID,
			Kind:        model.ShowKindMovie,
			Title:       recordingDetails.Title,
			RecordedOn:  recordingDetails.StartTime.Time,
			Description: recordingDetails.Description,
			MovieYear:   recordingDetails.MovieYear,
		},
		recordingDetails:  recordingDetails,
		collectionDetails: collectionDetails,
	}, nil
}

type series struct {
	*model.Series
	recordingDetails  *message.RecordingItem
	collectionDetails *message.CollectionItem
}

func newSeries(episode *episode) *series {

	return &series{
		Series: &model.Series{
			ID:          episode.SeriesID,
			Kind:        model.ShowKindSeries,
			Title:       episode.Title,
			RecordedOn:  episode.RecordedOn,
			Description: episode.Description,
			Episodes:    []*model.Episode{episode.Episode},
		},
		recordingDetails:  episode.recordingDetails,
		collectionDetails: episode.collectionDetails,
	}
}

type episode struct {
	*model.Episode
	recordingDetails  *message.RecordingItem
	collectionDetails *message.CollectionItem
}

func newEpisode(recordingDetails *message.RecordingItem, collectionDetails *message.CollectionItem) (*episode, error) {
	if recordingDetails.CollectionType != message.CollectionTypeSeries {
		return nil, fmt.Errorf("%w: unexpected collection type: %s",
			errorz.ErrResponse, recordingDetails.CollectionType)
	}

	var episodeNumber int
	if len(recordingDetails.EpisodeNum) > 0 {
		episodeNumber = recordingDetails.EpisodeNum[0]
	}

	return &episode{
		Episode: &model.Episode{
			ID:                 recordingDetails.RecordingID,
			SeriesID:           recordingDetails.CollectionID,
			Kind:               model.ShowKindEpisode,
			Title:              recordingDetails.Title,
			RecordedOn:         recordingDetails.StartTime.Time,
			Description:        collectionDetails.Description,
			OriginalAirDate:    recordingDetails.OriginalAirDate,
			SeasonNumber:       recordingDetails.SeasonNumber,
			EpisodeNumber:      episodeNumber,
			EpisodeTitle:       recordingDetails.Subtitle,
			EpisodeDescription: recordingDetails.Description,
		},
		recordingDetails:  recordingDetails,
		collectionDetails: collectionDetails,
	}, nil
}

func findBestImage(images []message.CollectionImage, width, height int) message.CollectionImage {

	slices.SortFunc(images, func(a, b message.CollectionImage) int {
		differenceA := mathx.Abs(a.Height-height) + mathx.Abs(a.Width-width)
		differenceB := mathx.Abs(b.Height-height) + mathx.Abs(b.Width-width)

		return differenceA - differenceB
	})

	return images[0]
}
