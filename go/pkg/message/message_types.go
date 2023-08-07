package message

type Type string

const (
	TypeRequest                   Type = "request"
	TypeResponse                  Type = "response"
	TypeBodyAuthenticate          Type = "bodyAuthenticate"
	TypeRecordingFolderItemSearch Type = "recordingFolderItemSearch"
	TypeRecordingFolderItemList   Type = "recordingFolderItemList"
	TypeRecordingSearch           Type = "recordingSearch"
	TypeRecordingList             Type = "recordingList"
	TypeCollectionSearch          Type = "collectionSearch"
	TypeCollectionList            Type = "collectionList"
	TypeError                     Type = "error"
)

type StatusType string

const (
	StatusTypeSuccess StatusType = "success"
	StatusTypeFailure StatusType = "failure"
)

type ResponseCount string

const (
	ResponseCountSingle   ResponseCount = "single"
	ResponseCountMultiple ResponseCount = "multiple"
)

type CredentialType string

const (
	CredentialTypeMak CredentialType = "makCredential"
)

type CollectionType string

const (
	CollectionTypeSeries  CollectionType = "series"
	CollectionTypeMovie   CollectionType = "movie"
	CollectionTypeSpecial CollectionType = "special"
)

type LevelOfDetail string

const (
	LevelOfDetailLow    LevelOfDetail = "low"
	LevelOfDetailMedium LevelOfDetail = "medium"
	LevelOfDetailHigh   LevelOfDetail = "high"
)
