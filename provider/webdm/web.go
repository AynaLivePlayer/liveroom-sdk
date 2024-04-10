package webdm

import (
	"errors"
	"github.com/AynaLivePlayer/blivedm-go/client"
	"github.com/AynaLivePlayer/blivedm-go/message"
	"github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/liveroom-sdk/utils"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

const ProviderName = "biliweb"

type WebDanmuClient struct {
	cfg            liveroom.LiveRoomConfig
	webDmClient    *client.Client
	onMessage      func(msg *liveroom.Message)
	onDisconnect   func(liveroom liveroom.LiveRoom)
	onStatusChange func(connected bool)
}

func NewWebDanmuClientProvider(apiServer string) liveroom.ILiveRoomProvider {
	return &liveroom.LiveRoomProvider{
		Name:        ProviderName,
		Description: "default web protocol. enter room id to connect.",
		Func: func(cfg liveroom.LiveRoomConfig) (liveroom.LiveRoom, error) {
			if cfg.Provider != ProviderName {
				return nil, errors.New("invalid provider name")
			}
			roomId, err := cast.ToIntE(cfg.Room)
			if err != nil {
				return nil, errors.New("invalid room id, should be integer")
			}
			room := &WebDanmuClient{
				cfg:         cfg,
				webDmClient: client.NewClientWithApi(roomId, &remoteApi{client: resty.New().SetBaseURL(apiServer)}),
			}
			room.webDmClient.OnDanmaku(room.danmuHandler)
			return room, nil
		},
	}
}

func (w *WebDanmuClient) GetName() string {
	return ProviderName
}

func (w *WebDanmuClient) Config() *liveroom.LiveRoomConfig {
	return &w.cfg
}

func (w *WebDanmuClient) danmuHandler(data *message.Danmaku) {
	if w.onMessage == nil {
		return
	}
	w.onMessage(&liveroom.Message{
		User: liveroom.User{
			Uid:       cast.ToString(data.Sender.Uid),
			Username:  data.Sender.Uname,
			Admin:     data.Sender.Admin,
			Privilege: data.Sender.GuardLevel,
			Medal: liveroom.UserMedal{
				Name:   data.Sender.Medal.Name,
				Level:  utils.BilibiliGuardLevelToPrivilege(data.Sender.GuardLevel),
				RoomID: cast.ToString(data.Sender.Medal.UpRoomId),
			},
		},
		Message: data.Content,
	})
}

func (w *WebDanmuClient) Connect() error {
	err := w.webDmClient.Start()
	if err == nil && w.onStatusChange != nil {
		w.onStatusChange(true)
	}
	return err
}

func (w *WebDanmuClient) Disconnect() error {
	w.webDmClient.Stop()
	if w.onStatusChange != nil {
		w.onStatusChange(false)
	}
	return nil
}

func (w *WebDanmuClient) OnDisconnect(f func(liveroom liveroom.LiveRoom)) {
	w.onDisconnect = f
}

func (w *WebDanmuClient) OnStatusChange(f func(connected bool)) {
	w.onStatusChange = f
}

func (w *WebDanmuClient) OnMessage(f func(msg *liveroom.Message)) {
	w.onMessage = f
}
