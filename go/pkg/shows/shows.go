package shows

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"

	liberrorz "github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/go/pkg/stringz"
	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"golang.org/x/exp/slices"
)

func New(tivo *model.Tivo, objectID string, recording *message.RecordingItem, collection *message.CollectionItem) model.Show {

	var result model.Show
	switch recording.CollectionType {

	case message.CollectionTypeSeries:
		result = newEpisode(tivo, objectID, recording, collection)

	case message.CollectionTypeMovie, message.CollectionTypeSpecial:
		result = newMovie(tivo, objectID, recording, collection)

	default:
		panic(fmt.Errorf("%w: unexpected collection type for recording '%s': '%v'",
			liberrorz.ErrFatal, recording.Title, recording.CollectionType))
	}

	return result
}

func WithImageURL(show model.Show, targetDimensions *apicontext.ImageDimensions) model.Show {

	if targetDimensions == nil {
		return show
	}

	var result model.Show

	switch show.GetKind() {
	case model.ShowKindMovie:
		movie := *(show.(*movie))
		movie.ImageURL = findBestImageURL(movie.Details.Collection.Images, targetDimensions)
		result = &movie

	case model.ShowKindSeries:
		series := *(show.(*series))
		series.ImageURL = findBestImageURL(series.Details.Collection.Images, targetDimensions)
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

	default:
		panic(fmt.Errorf("%w: unexpected show kind: %v", liberrorz.ErrFatal, show.GetKind()))
	}
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

func ParseIDNumber(id string) string {
	// example tivo ID: tivo:rc.20479
	split := strings.Split(id, ".")
	return split[len(split)-1]
}

type Details struct {
	Tivo       *model.Tivo
	ObjectID   string
	Recording  message.RecordingItem
	Collection message.CollectionItem
}

func GetDetails(show model.Show) *Details {

	switch s := show.(type) {

	case *movie:
		return &s.Details

	case *series:
		return &s.Details

	case *episode:
		return &s.Details

	default:
		return nil
	}
}

func GetPath(show model.Show) string {

	details := GetDetails(show)
	if details == nil {
		return stringz.ToAlphaNumeric(show.GetTitle())
	}

	switch show.GetKind() {

	case model.ShowKindMovie:
		return stringz.ToAlphaNumeric(details.Recording.Title)

	case model.ShowKindSeries:
		return stringz.ToAlphaNumeric(details.Collection.Title)

	case model.ShowKindEpisode:
		parentDir := stringz.ToAlphaNumeric(details.Collection.Title)
		subDir := stringz.ToAlphaNumeric(details.Recording.Title)
		return path.Join(parentDir, subDir)

	default:
		panic(fmt.Errorf("%w: unexpected show kind '%s'", liberrorz.ErrFatal, show.GetKind()))
	}
}

type movie struct {
	*model.Movie
	Details Details
}

func newMovie(tivo *model.Tivo, objectID string, recording *message.RecordingItem, collection *message.CollectionItem) *movie {

	if recording.CollectionType != message.CollectionTypeMovie &&
		recording.CollectionType != message.CollectionTypeSpecial {

		panic(fmt.Errorf("%w: unexpected collection type for recording '%s': '%v'",
			liberrorz.ErrFatal, recording.Title, recording.CollectionType))
	}

	return &movie{
		Movie: &model.Movie{
			ID:          recording.RecordingID,
			Kind:        model.ShowKindMovie,
			Title:       recording.Title,
			RecordedOn:  recording.StartTime.Time,
			Description: recording.Description,
			MovieYear:   recording.MovieYear,
		},
		Details: Details{
			Tivo:       tivo,
			ObjectID:   objectID,
			Recording:  *recording,
			Collection: *collection,
		},
	}
}

type series struct {
	*model.Series
	Details Details
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
		Details: episode.Details,
	}
}

type episode struct {
	*model.Episode
	Details Details
}

func newEpisode(tivo *model.Tivo, objectID string, recording *message.RecordingItem, collection *message.CollectionItem) *episode {

	if recording.CollectionType != message.CollectionTypeSeries {
		panic(fmt.Errorf("%w: unexpected collection type for recording '%s': '%v'",
			liberrorz.ErrFatal, recording.Title, recording.CollectionType))
	}

	var episodeNumber int
	if len(recording.EpisodeNum) > 0 {
		episodeNumber = recording.EpisodeNum[0]
	}

	return &episode{
		Episode: &model.Episode{
			ID:                 recording.RecordingID,
			SeriesID:           recording.CollectionID,
			Kind:               model.ShowKindEpisode,
			Title:              recording.Title,
			RecordedOn:         recording.StartTime.Time,
			Description:        collection.Description,
			OriginalAirDate:    recording.OriginalAirDate,
			SeasonNumber:       recording.SeasonNumber,
			EpisodeNumber:      episodeNumber,
			EpisodeTitle:       recording.Subtitle,
			EpisodeDescription: recording.Description,
		},
		Details: Details{
			Tivo:       tivo,
			ObjectID:   objectID,
			Recording:  *recording,
			Collection: *collection,
		},
	}
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
