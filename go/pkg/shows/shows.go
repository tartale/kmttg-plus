package shows

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
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

func WithImageURL(show model.Show, targetDimensions *apicontext.ImageDimensions) model.Show {

	if targetDimensions == nil {
		return show
	}

	var result model.Show

	switch show.GetKind() {
	case model.ShowKindMovie:
		movie := *(show.(*movie))
		movie.ImageURL = findBestImageURL(movie.collectionDetails.Images, targetDimensions)
		result = &movie

	case model.ShowKindSeries:
		series := *(show.(*series))
		series.ImageURL = findBestImageURL(series.collectionDetails.Images, targetDimensions)
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

func imageIsInvalid(image message.CollectionImage) bool {

	resp, err := http.Get(image.ImageURL)
	if err != nil {
		return true
	}
	if resp.StatusCode == http.StatusOK {
		return false
	}

	return true
}

func findBestImageURL(images []message.CollectionImage, target *apicontext.ImageDimensions) string {

	if len(images) == 0 || target == nil {
		return ""
	}
	slices.SortFunc(images, func(a, b message.CollectionImage) int {
		differenceA := mathx.Abs(a.Height-target.Height) + mathx.Abs(a.Width-target.Width)
		differenceB := mathx.Abs(b.Height-target.Height) + mathx.Abs(b.Width-target.Width)

		return differenceA - differenceB
	})
	for _, image := range images {
		if !imageIsInvalid(image) {
			return image.ImageURL
		}
	}

	return ""
}
