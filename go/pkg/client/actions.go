package client

import (
	"context"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/apicontext"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/shows"
	"go.uber.org/zap"
)

const maxLimitValue = 50

func (t *TivoClient) Authenticate(ctx context.Context) error {
	authRequest := message.NewTivoMessage().WithAuthRequest(config.Values.MediaAccessKey)
	err := t.Send(ctx, authRequest)
	if err != nil {
		return err
	}

	authResponseBody := &message.AuthResponseBody{}
	authResponse := message.NewTivoMessage().WithBody(authResponseBody)
	err = t.Receive(context.Background(), authResponse)
	if err != nil {
		return err
	}
	if authResponseBody.Status != message.StatusTypeSuccess {
		return fmt.Errorf("%w: %s", errorz.ErrAuthenticationFailed, authResponseBody.Message)
	}

	return nil
}

func (t *TivoClient) GetShows(ctx context.Context) ([]model.Show, error) {
	var result []model.Show
	ctx = apicontext.Wrap(ctx).
		WithShowOffset(0).
		WithShowLimit(maxLimitValue)

	for {
		shows, nextOffset, err := t.getShowsPage(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, shows...)
		if nextOffset < 0 {
			break
		}
		ctx = apicontext.Wrap(ctx).WithShowOffset(nextOffset)
	}
	result = shows.MergeEpisodes(result)

	return result, nil
}

func (t *TivoClient) getShowsPage(ctx context.Context) (shows []model.Show, nextOffset int, err error) {
	request := message.NewTivoMessage().WithGetRecordingListRequest(ctx, t.BodyID())
	err = t.Send(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	responseBody := &message.RecordingFolderItemSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, 0, err
	}
	if responseBody.Type != message.TypeRecordingFolderItemList {
		logz.Logger.Warn("tivo error response", zap.Any("request", request), zap.Any("response", responseBody))
		return nil, 0, fmt.Errorf("%w: unexpected response type: %s", errorz.ErrUnexpectedResponse, responseBody.Type)
	}
	for _, recording := range responseBody.RecordingFolderItem {
		show, err := t.getShowDetails(ctx, recording)
		if err != nil {
			return nil, 0, err
		}
		shows = append(shows, show)
	}
	showCount := len(shows)
	if showCount < apicontext.ShowLimit(ctx) {
		nextOffset = -1
	} else {
		nextOffset = apicontext.ShowOffset(ctx) + showCount
	}

	return
}

func (t *TivoClient) getShowDetails(ctx context.Context, recordingFolderItem message.RecordingFolderItem) (model.Show, error) {
	recordingDetails, err := t.getRecordingDetails(ctx, recordingFolderItem)
	if err != nil {
		return nil, err
	}
	collectionDetails, err := t.getCollectionDetails(ctx, []string{recordingFolderItem.CollectionID})
	if err != nil {
		return nil, err
	}
	objectID, err := t.getObjectID(ctx, recordingDetails.RecordingID)
	if err != nil {
		return nil, err
	}

	show := shows.New(t.tivo, objectID, recordingDetails, collectionDetails)

	return show, nil
}

func (t *TivoClient) getRecordingDetails(ctx context.Context, recordingFolderItem message.RecordingFolderItem) (*message.RecordingItem, error) {
	request := message.NewTivoMessage().WithGetRecordingRequest(ctx, t.BodyID(), recordingFolderItem.ChildRecordingID)
	err := t.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	responseBody := &message.RecordingSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, err
	}
	if responseBody.Type != message.TypeRecordingList {
		logz.Logger.Error("tivo error response", zap.Any("responseBody", responseBody))
		return nil, fmt.Errorf("%w: unexpected response type: %s", errorz.ErrUnexpectedResponse, responseBody.Type)
	}
	recordingCount := len(responseBody.Recording)
	if recordingCount != 1 {
		return nil, fmt.Errorf("%w: unexpected number of recordings in response: %d", errorz.ErrUnexpectedResponse, recordingCount)
	}
	recording := responseBody.Recording[0]
	recording.RecordingID = recordingFolderItem.ChildRecordingID
	recording.CollectionID = recordingFolderItem.CollectionID

	return &recording, nil
}

func (t *TivoClient) getCollectionDetails(ctx context.Context, collectionIDs []string) (*message.CollectionItem, error) {
	request := message.NewTivoMessage().WithGetCollectionRequest(ctx, collectionIDs)
	err := t.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	responseBody := &message.CollectionSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, err
	}
	if responseBody.Type != message.TypeCollectionList {
		logz.Logger.Error("tivo error response", zap.Any("responseBody", responseBody))
		return nil, fmt.Errorf("%w: unexpected response type: %s", errorz.ErrUnexpectedResponse, responseBody.Type)
	}
	collectionCount := len(responseBody.Collection)
	if collectionCount != len(collectionIDs) {
		return nil, fmt.Errorf("%w: unexpected number of collection items in response: %d", errorz.ErrUnexpectedResponse, collectionCount)
	}

	return &responseBody.Collection[0], nil
}

func (t *TivoClient) getObjectID(ctx context.Context, searchID string) (string, error) {
	request := message.NewTivoMessage().WithIdSearchRequest(ctx, t.BodyID(), searchID)
	err := t.Send(ctx, request)
	if err != nil {
		return "", err
	}

	responseBody := &message.IdSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return "", err
	}
	if responseBody.Type != message.TypeIdSet {
		logz.Logger.Warn("tivo error response", zap.Any("responseBody", responseBody))
		return "", fmt.Errorf("%w: unexpected response type: %s", errorz.ErrUnexpectedResponse, responseBody.Type)
	}
	objectIdCount := len(responseBody.ObjectID)
	if objectIdCount != 1 {
		return "", fmt.Errorf("%w: unexpected number of objectID items in response: %d", errorz.ErrUnexpectedResponse, objectIdCount)
	}

	return responseBody.ObjectID[0], nil
}
