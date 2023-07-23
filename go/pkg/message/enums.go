package message

type Type string

const (
	TypeBodyAuthenticate          Type = "bodyAuthenticate"
	TypeRecordingFolderItemSearch Type = "recordingFolderItemSearch"
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
	CollectionTypeSeries CollectionType = "series"
)
