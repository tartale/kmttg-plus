package message

type CollectionSearchRequestBody struct {
	Type          Type          `json:"type,omitempty"`
	BodyID        string        `json:"bodyId,omitempty"`
	CollectionIDs []string      `json:"collectionId,omitempty"`
	LevelOfDetail LevelOfDetail `json:"levelOfDetail,omitempty"`
	Offset        *int          `json:"offset,omitempty"`
	Count         *int          `json:"count,omitempty"`
}

type CollectionSearchResponseBody struct {
	Type       Type             `json:"type,omitempty"`
	Status     StatusType       `json:"status,omitempty"`
	Message    string           `json:"message,omitempty"`
	Collection []CollectionItem `json:"collection,omitempty"`
}

type CollectionItem struct {
	CollectionID   string            `json:"collectionId,omitempty"`
	CollectionType CollectionType    `json:"collectionType,omitempty"`
	Description    string            `json:"description,omitempty"`
	Episodic       bool              `json:"episodic,omitempty"`
	Title          string            `json:"title,omitempty"`
	TVRating       string            `json:"tvRating,omitempty"`
	Images         []CollectionImage `json:"image,omitempty"`
}

type CollectionImage struct {
	ImageURL  string              `json:"imageUrl,omitempty"`
	ImageType CollectionImageType `json:"imageType,omitempty"`
	Height    int                 `json:"height,omitempty"`
	Width     int                 `json:"width,omitempty"`
}

type CollectionImageType string

const (
	CollectionImageTypeShowcaseBanner = "showcaseBanner"
)
