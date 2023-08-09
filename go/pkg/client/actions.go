package client

import (
	"context"
	"fmt"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/errorz"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/message"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"go.uber.org/zap"
)

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
		return errorz.ErrNotAuthenticated(authResponseBody.Message)
	}

	return nil
}

func (t *TivoClient) GetShows(ctx context.Context) ([]model.Show, error) {

	request := message.NewTivoMessage().WithGetRecordingListRequest(ctx, t.BodyID())
	err := t.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	responseBody := &message.RecordingFolderItemSearchResponseBody{}
	response := message.NewTivoMessage().WithBody(responseBody)
	err = t.Receive(ctx, response)
	if err != nil {
		return nil, err
	}
	if responseBody.Type != message.TypeRecordingFolderItemList {
		logz.Logger.Warn("tivo error response", zap.Any("request", request), zap.Any("response", responseBody))
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected response type: %s", responseBody.Type))
	}
	var result []model.Show
	for _, recording := range responseBody.RecordingFolderItem {
		show, err := t.GetShowDetails(ctx, recording)
		if err != nil {
			return nil, err
		}
		result = append(result, show)
	}
	result = model.MergeEpisodes(result)

	return result, nil
}

func (t *TivoClient) GetShowDetails(ctx context.Context, recordingFolderItem message.RecordingFolderItem) (model.Show, error) {

	recordingDetails, err := t.GetRecordingDetails(ctx, recordingFolderItem)
	if err != nil {
		return nil, err
	}
	collectionDetails, err := t.GetCollectionDetails(ctx, []string{recordingFolderItem.CollectionID})
	if err != nil {
		return nil, err
	}

	show, err := model.NewShow(recordingDetails, &collectionDetails[0])
	if err != nil {
		return nil, err
	}

	return show, nil
}

func (t *TivoClient) GetRecordingDetails(ctx context.Context, recordingFolderItem message.RecordingFolderItem) (*message.RecordingItem, error) {

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
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected response type: %s", responseBody.Type))
	}
	recordingCount := len(responseBody.Recording)
	if recordingCount != 1 {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected number of recordings in response: %d", recordingCount))
	}
	recording := responseBody.Recording[0]
	recording.RecordingID = recordingFolderItem.ChildRecordingID
	recording.CollectionID = recordingFolderItem.CollectionID

	return &recording, nil
}

func (t *TivoClient) GetCollectionDetails(ctx context.Context, collectionIDs []string) ([]message.CollectionItem, error) {

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
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected response type: %s", responseBody.Type))
	}
	collectionCount := len(responseBody.Collection)
	if collectionCount != len(collectionIDs) {
		return nil, errorz.ErrResponse(fmt.Sprintf("unexpected number of collection items in response: %d", collectionCount))
	}

	return responseBody.Collection, nil
}
