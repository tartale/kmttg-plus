package message

import "github.com/tartale/go/pkg/formattedtime"

type RecordingFolderItemSearchRequestBody struct {
	Type    Type   `json:"type,omitempty"`
	BodyID  string `json:"bodyId,omitempty"`
	Offset  *int   `json:"offset,omitempty"`
	Count   *int   `json:"count,omitempty"`
	Flatten *bool  `json:"flatten,omitempty"`
}

type RecordingFolderItemSearchResponseBody struct {
	Type                Type                  `json:"type,omitempty"`
	RecordingFolderItem []RecordingFolderItem `json:"recordingFolderItem,omitempty"`
}

type RecordingFolderItem struct {
	ChildRecordingID      string             `json:"childRecordingId,omitempty"`
	RecordingFolderItemID string             `json:"recordingFolderItemId,omitempty"`
	StartTime             formattedtime.Time `json:"start_time,omitempty"`
	Title                 string             `json:"title,omitempty"`
	CollectionType        CollectionType     `json:"collectionType,omitempty"`
	PercentWatched        int                `json:"percentWatched,omitempty"`
}
