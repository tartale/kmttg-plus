package message

type ClipMetadata struct {
	Type           Type                   `json:"type,omitempty"`
	ClipMetadataID string                 `json:"clipMetadataId,omitempty"`
	ContentID      string                 `json:"contentId,omitempty"`
	SegmentType    SegmentType            `json:"segmentType,omitempty"`
	Segment        []ClipMetadataSegment  `json:"segment,omitempty"`
	SyncMark       []ClipMetadataSyncMark `json:"syncMark,omitempty"`
}

type ClipMetadataSegment struct {
	Type        Type     `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	Keyword     []string `json:"keyword,omitempty"`
	StartOffset string   `json:"startOffset,omitempty"`
	EndOffset   string   `json:"endOffset,omitempty"`
}

type ClipMetadataSyncMark struct {
	Type      Type   `json:"type,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type ClipMetadataSearchRequest struct {
	Type           Type     `json:"type,omitempty"`
	ClipMetadataID []string `json:"clipMetadataId,omitempty"`
	ContentID      string   `json:"contentId,omitempty"`
}

type ClipMetadataSearchResponse struct {
	Type         Type           `json:"type,omitempty"`
	ClipMetadata []ClipMetadata `json:"clipMetadata,omitempty"`
}
