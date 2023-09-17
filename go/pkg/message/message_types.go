package message

type Type string

const (
	TypeRequest  Type = "request"
	TypeResponse Type = "response"
	TypeError    Type = "error"

	TypeBodyAuthenticate          Type = "bodyAuthenticate"
	TypeCategorySearch            Type = "categorySearch"
	TypeChannelSearch             Type = "channelSearch"
	TypeClipMetadataSearch        Type = "clipMetadataSearch"
	TypeClipMetadata              Type = "clipMetadata"
	TypeClipMetadataList          Type = "clipMetadatList"
	TypeClipSegment               Type = "clipSegment"
	TypeClipSyncMark              Type = "clipSyncMark"
	TypeCollectionList            Type = "collectionList"
	TypeCollectionSearch          Type = "collectionSearch"
	TypeContentSearch             Type = "contentSearch"
	TypeIdSearch                  Type = "idSearch"
	TypeIdSet                     Type = "idSet"
	TypeOfferSearch               Type = "offerSearch"
	TypeRecordingFolderItemList   Type = "recordingFolderItemList"
	TypeRecordingFolderItemSearch Type = "recordingFolderItemSearch"
	TypeRecordingList             Type = "recordingList"
	TypeRecordingSearch           Type = "recordingSearch"
	TypeWhatsOnSearch             Type = "whatsOnSearch"
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

type IdNamespace string

const (
	IdNamespaceMFS IdNamespace = "mfs"
)

type IdType string

type SegmentType string

const (
	SegmentTypeAdSkip = "adSkip"
)
