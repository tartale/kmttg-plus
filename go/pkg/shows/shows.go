package shows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	liberrorz "github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/mathx"
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
		panic(fmt.Errorf("unexpected collection type for recording '%s': '%v': %w",
			recording.Title, recording.CollectionType, liberrorz.ErrFatal))
	}

	return result
}

// UnmarshalShowFromJSON unmarshals a Show from JSON bytes, handling both
// wrapper types (with Details) and plain model types.
func UnmarshalShowFromJSON(data []byte, tivo *model.Tivo) (model.Show, error) {
	var kind struct {
		Kind model.ShowKind `json:"kind"`
	}
	if err := json.Unmarshal(data, &kind); err != nil {
		return nil, fmt.Errorf("unmarshal show kind: %w", err)
	}

	// Check if Details field exists
	var hasDetails struct {
		Details json.RawMessage `json:"details,omitempty"`
	}
	json.Unmarshal(data, &hasDetails)

	if len(hasDetails.Details) > 0 {
		// Unmarshal into wrapper type with Details
		switch kind.Kind {
		case model.ShowKindMovie:
			m := movie{Movie: &model.Movie{}}
			if err := json.Unmarshal(data, &m); err != nil {
				return nil, fmt.Errorf("unmarshal movie with details: %w", err)
			}
			if m.Details.Tivo == nil {
				m.Details.Tivo = tivo
			}
			return &m, nil
		case model.ShowKindSeries:
			s := series{Series: &model.Series{}}
			if err := json.Unmarshal(data, &s); err != nil {
				return nil, fmt.Errorf("unmarshal series with details: %w", err)
			}
			if s.Details.Tivo == nil {
				s.Details.Tivo = tivo
			}
			return &s, nil
		case model.ShowKindEpisode:
			e := episode{Episode: &model.Episode{}}
			if err := json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("unmarshal episode with details: %w", err)
			}
			if e.Details.Tivo == nil {
				e.Details.Tivo = tivo
			}
			return &e, nil
		default:
			return nil, fmt.Errorf("unknown show kind: %s", kind.Kind)
		}
	}

	// No Details field, unmarshal into plain model type
	switch kind.Kind {
	case model.ShowKindMovie:
		var m model.Movie
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("unmarshal movie: %w", err)
		}
		return &m, nil
	case model.ShowKindSeries:
		var s model.Series
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("unmarshal series: %w", err)
		}
		return &s, nil
	case model.ShowKindEpisode:
		var e model.Episode
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("unmarshal episode: %w", err)
		}
		return &e, nil
	default:
		return nil, fmt.Errorf("unknown show kind: %s", kind.Kind)
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
		panic(fmt.Errorf("unexpected show kind: %v: %w", show.GetKind(), liberrorz.ErrFatal))
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
	Tivo       *model.Tivo            `json:"-"`
	ObjectID   string                 `json:"objectID,omitempty"`
	Recording  message.RecordingItem  `json:"recording"`
	Collection message.CollectionItem `json:"collection"`
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

func GetCanonicalName(show model.Show) string {
	switch s := show.(type) {
	case *movie:
		return s.CanonicalName()

	case *series:
		return s.CanonicalName()

	case *episode:
		return s.CanonicalName()

	default:
		return s.GetTitle()
	}
}

type movie struct {
	*model.Movie
	Details Details `json:"details"`
}

func newMovie(tivo *model.Tivo, objectID string, recording *message.RecordingItem, collection *message.CollectionItem) *movie {
	if recording.CollectionType != message.CollectionTypeMovie &&
		recording.CollectionType != message.CollectionTypeSpecial {

		panic(fmt.Errorf("unexpected collection type for recording '%s': '%v': %w",
			recording.Title, recording.CollectionType, liberrorz.ErrFatal))
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

func (m *movie) CanonicalName() string {
	return fmt.Sprintf("%s (%d)", m.Title, m.MovieYear)
}

type series struct {
	*model.Series
	Details Details `json:"details,omitempty"`
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

func (s *series) CanonicalName() string {
	return s.Title
}

type episode struct {
	*model.Episode
	Details Details `json:"details,omitempty"`
}

func newEpisode(tivo *model.Tivo, objectID string, recording *message.RecordingItem, collection *message.CollectionItem) *episode {
	if recording.CollectionType != message.CollectionTypeSeries {
		panic(fmt.Errorf("unexpected collection type for recording '%s': '%v': %w",
			recording.Title, recording.CollectionType, liberrorz.ErrFatal))
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

func (e *episode) CanonicalName() string {
	return fmt.Sprintf("%s - %s - E%02dS%02d", e.Title, e.EpisodeTitle, e.SeasonNumber, e.EpisodeNumber)
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
