package webdm

import (
	"encoding/json"
	"errors"
	"github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type remoteApi struct {
	client *resty.Client
}

type apiData struct {
	UID       int           `json:"uid"`
	DanmuInfo api.DanmuInfo `json:"danmu_info"`
	Error     string        `json:"error"`
}

type apiResponse struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data apiData `json:"data"`
}

func (r remoteApi) GetDanmuInfo(roomID int) (int, *api.DanmuInfo, error) {
	resp, err := r.client.R().
		SetQueryParam("room_id", strconv.Itoa(roomID)).
		Get("/api/blivedm/web/dm_info")
	if err != nil {
		return 0, nil, err
	}
	var sceneResp apiResponse
	err = json.Unmarshal(resp.Body(), &sceneResp)
	if err != nil {
		return 0, nil, err
	}
	if sceneResp.Code != 0 {
		return 0, nil, errors.New(sceneResp.Msg)
	}
	if sceneResp.Data.Error != "" {
		return 0, nil, errors.New(sceneResp.Data.Error)
	}
	return sceneResp.Data.UID, &sceneResp.Data.DanmuInfo, nil
}

func (r remoteApi) GetRoomInfo(roomID int) (*api.RoomInfo, error) {
	return api.GetRoomInfo(roomID)
}
