package message

import (
	"github.com/tartale/go/pkg/jsontime"
)

type RecordingFolderItemSearchRequestBody struct {
	Type         Type    `json:"type,omitempty"`
	BodyID       string  `json:"bodyId,omitempty"`
	CollectionID *string `json:"collectionId,omitempty"`
	Offset       *int    `json:"offset,omitempty"`
	Count        *int    `json:"count,omitempty"`
	Flatten      *bool   `json:"flatten,omitempty"`
}

type RecordingFolderItemSearchResponseBody struct {
	Type                Type                  `json:"type,omitempty"`
	Status              StatusType            `json:"status,omitempty"`
	Message             string                `json:"message,omitempty"`
	RecordingFolderItem []RecordingFolderItem `json:"recordingFolderItem,omitempty"`
}

type RecordingFolderItem struct {
	ChildRecordingID      string         `json:"childRecordingId,omitempty"`
	RecordingFolderItemID string         `json:"recordingFolderItemId,omitempty"`
	CollectionID          string         `json:"collectionId,omitempty"`
	StartTime             jsontime.Time  `json:"startTime,omitempty" format:"2006-01-02 15:04:05"`
	Title                 string         `json:"title,omitempty"`
	CollectionType        CollectionType `json:"collectionType,omitempty"`
	PercentWatched        int            `json:"percentWatched,omitempty"`
}

type RecordingSearchRequestBody struct {
	Type          Type          `json:"type,omitempty"`
	BodyID        string        `json:"bodyId,omitempty"`
	LevelOfDetail LevelOfDetail `json:"levelOfDetail,omitempty"`
	RecordingID   string        `json:"recordingId,omitempty"`
}

type RecordingSearchResponseBody struct {
	Type      Type            `json:"type,omitempty"`
	Status    StatusType      `json:"status,omitempty"`
	Message   string          `json:"message,omitempty"`
	IsBottom  bool            `json:"isBottom,omitempty"`
	IsTop     bool            `json:"isTop,omitempty"`
	Recording []RecordingItem `json:"recording,omitempty"`
}

type RecordingItem struct {
	BodyID          string         `json:"bodyId,omitempty"`
	RecordingID     string         `json:"recordingId,omitempty"`
	CollectionID    string         `json:"collectionId,omitempty"`
	CollectionType  CollectionType `json:"collectionType,omitempty"`
	Description     string         `json:"description,omitempty"`
	EpisodeNum      []int          `json:"episodeNum,omitempty"`
	Episodic        bool           `json:"episodic,omitempty"`
	IsEpisode       bool           `json:"isEpisode,omitempty"`
	OriginalAirDate string         `json:"originalAirdate,omitempty"`
	SeasonNumber    int            `json:"seasonNumber,omitempty"`
	ShortTitle      string         `json:"shortTitle,omitempty"`
	StartTime       jsontime.Time  `json:"startTime,omitempty" format:"2006-01-02 15:04:05"`
	Subtitle        string         `json:"subtitle,omitempty"`
	Title           string         `json:"title,omitempty"`
	MovieYear       int            `json:"movieYear,omitempty"`
}
