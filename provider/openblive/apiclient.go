package openblive

import (
	"encoding/json"
	"errors"
	"github.com/aynakeya/open-bilibili-live"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

type remoteApiClient struct {
	client *resty.Client
}

type openbliveResponse struct {
	Result *openblive.AppStartResult `json:"result"`
	Error  *openblive.PublicError    `json:"error"`
}

func newRemoteApiClient(remoteApi string) openblive.IApiClient {
	return &remoteApiClient{
		client: resty.New().SetBaseURL(remoteApi),
	}
}

type sceneAppResponse struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data openbliveResponse `json:"data"`
}

func parseApiResponse(resp *resty.Response, err error) (*openblive.AppStartResult, *openblive.PublicError) {
	if err != nil {
		return nil, openblive.ErrUnknown.WithDetail(err)
	}
	var sceneResp sceneAppResponse
	err = json.Unmarshal(resp.Body(), &sceneResp)
	if err != nil {
		return nil, openblive.ErrUnknown.WithDetail(err)
	}
	if sceneResp.Code != 0 {
		return nil, openblive.ErrUnknown.WithDetail(errors.New(sceneResp.Msg))
	}
	if sceneResp.Data.Error != nil && sceneResp.Data.Error.Code == 0 && sceneResp.Data.Error.Message == "" {
		return sceneResp.Data.Result, nil
	}
	return sceneResp.Data.Result, sceneResp.Data.Error

}

func (r *remoteApiClient) AppStart(code string, appId int64) (*openblive.AppStartResult, *openblive.PublicError) {
	return parseApiResponse(r.client.R().
		SetQueryParam("code", code).
		SetQueryParam("app_id", cast.ToString(appId)).
		Get("/api/blivedm/openblive/app_start"))
}

func (r *remoteApiClient) AppEnd(appId int64, gameId string) *openblive.PublicError {
	_, err := parseApiResponse(r.client.R().
		SetQueryParam("app_id", cast.ToString(appId)).
		SetQueryParam("game_id", gameId).
		Get("/api/blivedm/openblive/app_end"))
	return err
}

func (r *remoteApiClient) HearBeat(gameId string) *openblive.PublicError {
	_, err := parseApiResponse(r.client.R().
		SetQueryParam("game_id", gameId).
		Get("/api/blivedm/openblive/heartbeat"))
	return err
}
