// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/tartale/kmttg-plus/go/pkg/filter"
)

type Show interface {
	IsShow()
	GetKind() ShowKind
	GetTitle() string
	GetRecordedOn() time.Time
	GetDescription() string
}

type Episode struct {
	Kind               ShowKind  `json:"kind"`
	RecordingID        string    `json:"recordingId"`
	CollectionID       string    `json:"collectionId"`
	Title              string    `json:"title"`
	RecordedOn         time.Time `json:"recordedOn"`
	Description        string    `json:"description"`
	OriginalAirDate    string    `json:"originalAirDate"`
	SeasonNumber       int       `json:"seasonNumber"`
	EpisodeNumber      int       `json:"episodeNumber"`
	EpisodeTitle       string    `json:"episodeTitle"`
	EpisodeDescription string    `json:"episodeDescription"`
}

func (Episode) IsShow()                       {}
func (this Episode) GetKind() ShowKind        { return this.Kind }
func (this Episode) GetTitle() string         { return this.Title }
func (this Episode) GetRecordedOn() time.Time { return this.RecordedOn }
func (this Episode) GetDescription() string   { return this.Description }

type EpisodeFilter struct {
	Kind               *filter.Operator `json:"kind,omitempty"`
	Title              *filter.Operator `json:"title,omitempty"`
	RecordedOn         *filter.Operator `json:"recordedOn,omitempty"`
	Description        *filter.Operator `json:"description,omitempty"`
	OriginalAirDate    *filter.Operator `json:"originalAirDate,omitempty"`
	SeasonNumber       *filter.Operator `json:"seasonNumber,omitempty"`
	EpisodeNumber      *filter.Operator `json:"episodeNumber,omitempty"`
	EpisodeTitle       *filter.Operator `json:"episodeTitle,omitempty"`
	EpisodeDescription *filter.Operator `json:"episodeDescription,omitempty"`
}

type Movie struct {
	Kind        ShowKind  `json:"kind"`
	RecordingID string    `json:"recordingId"`
	Title       string    `json:"title"`
	RecordedOn  time.Time `json:"recordedOn"`
	Description string    `json:"description"`
	MovieYear   int       `json:"movieYear"`
}

func (Movie) IsShow()                       {}
func (this Movie) GetKind() ShowKind        { return this.Kind }
func (this Movie) GetTitle() string         { return this.Title }
func (this Movie) GetRecordedOn() time.Time { return this.RecordedOn }
func (this Movie) GetDescription() string   { return this.Description }

type MovieFilter struct {
	Kind        *filter.Operator `json:"kind,omitempty"`
	Title       *filter.Operator `json:"title,omitempty"`
	RecordedOn  *filter.Operator `json:"recordedOn,omitempty"`
	Description *filter.Operator `json:"description,omitempty"`
	MovieYear   *filter.Operator `json:"movieYear,omitempty"`
}

type Series struct {
	Kind         ShowKind   `json:"kind"`
	CollectionID string     `json:"collectionId"`
	Title        string     `json:"title"`
	RecordedOn   time.Time  `json:"recordedOn"`
	Description  string     `json:"description"`
	Episodes     []*Episode `json:"episodes"`
}

func (Series) IsShow()                       {}
func (this Series) GetKind() ShowKind        { return this.Kind }
func (this Series) GetTitle() string         { return this.Title }
func (this Series) GetRecordedOn() time.Time { return this.RecordedOn }
func (this Series) GetDescription() string   { return this.Description }

type SeriesFilter struct {
	Kind        *filter.Operator `json:"kind,omitempty"`
	Title       *filter.Operator `json:"title,omitempty"`
	RecordedOn  *filter.Operator `json:"recordedOn,omitempty"`
	Description *filter.Operator `json:"description,omitempty"`
}

type ShowFilter struct {
	Kind        *filter.Operator `json:"kind,omitempty"`
	Title       *filter.Operator `json:"title,omitempty"`
	RecordedOn  *filter.Operator `json:"recordedOn,omitempty"`
	Description *filter.Operator `json:"description,omitempty"`
}

type SortBy struct {
	Field     interface{}   `json:"field"`
	Direction SortDirection `json:"direction"`
}

type Sorter struct {
	Fields []*SortBy `json:"fields"`
}

type Tivo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Tsn     string `json:"tsn"`
	Shows   []Show `json:"shows,omitempty"`
}

type TivoFilter struct {
	Name *filter.Operator `json:"name,omitempty"`
}

type ShowKind string

const (
	ShowKindMovie   ShowKind = "MOVIE"
	ShowKindSeries  ShowKind = "SERIES"
	ShowKindEpisode ShowKind = "EPISODE"
)

var AllShowKind = []ShowKind{
	ShowKindMovie,
	ShowKindSeries,
	ShowKindEpisode,
}

func (e ShowKind) IsValid() bool {
	switch e {
	case ShowKindMovie, ShowKindSeries, ShowKindEpisode:
		return true
	}
	return false
}

func (e ShowKind) String() string {
	return string(e)
}

func (e *ShowKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ShowKind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ShowKind", str)
	}
	return nil
}

func (e ShowKind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

var AllSortDirection = []SortDirection{
	SortDirectionAsc,
	SortDirectionDesc,
}

func (e SortDirection) IsValid() bool {
	switch e {
	case SortDirectionAsc, SortDirectionDesc:
		return true
	}
	return false
}

func (e SortDirection) String() string {
	return string(e)
}

func (e *SortDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirection", str)
	}
	return nil
}

func (e SortDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
