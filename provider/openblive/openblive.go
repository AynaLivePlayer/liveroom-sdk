package openblive

import (
	"context"
	openblive "github.com/aynakeya/open-bilibili-live"
	"liveroom"
	"strconv"
)

const ProviderName = "openblive"

type OpenBLiveClient struct {
	cfg             liveroom.LiveRoomConfig
	openbliveClient *openblive.BLiveClient
	conn            openblive.BLiveLongConnection
	onMessage       func(msg *liveroom.Message)
	onDisconnect    func(liveroom liveroom.LiveRoom)
	onStatusChange  func(connected bool)
}

func NewOpenBLiveClientProvider(apiServer string, appId int64) liveroom.LiveRoomProvider {
	return func(cfg liveroom.LiveRoomConfig) liveroom.LiveRoom {
		if cfg.Provider != ProviderName {
			return nil
		}
		return &OpenBLiveClient{
			cfg:             cfg,
			openbliveClient: openblive.NewBliveClient(appId, cfg.Room, newRemoteApiClient(apiServer)),
		}
	}
}

func guardLevelToPrivilege(level int) int {
	switch level {
	case 1:
		return liveroom.PrivilegeUltimate
	case 2:
		return liveroom.PrivilegeAdvanced
	case 3:
		return liveroom.PrivilegeBasic
	default:
		return liveroom.PrivilegeNone
	}
}

func (o *OpenBLiveClient) danmuHandler(data openblive.DanmakuData) {
	if o.onMessage == nil {
		return
	}
	o.onMessage(&liveroom.Message{
		User: liveroom.User{
			Uid:       data.OpenID,
			Username:  data.UName,
			Admin:     false, // not supported by open bilibili live
			Privilege: guardLevelToPrivilege(data.GuardLevel),
			Medal: liveroom.UserMedal{
				Name:   data.FansMedalName,
				Level:  data.FansMedalLevel,
				RoomID: strconv.Itoa(data.RoomID),
			},
		},
		Message: data.Msg,
	})
}

func (o *OpenBLiveClient) disconnectHandler(conn openblive.BLiveLongConnection) {
	if o.onStatusChange != nil {
		o.onStatusChange(false)
	}
	if o.onDisconnect != nil {
		o.onDisconnect(o)
	}
}

func (o *OpenBLiveClient) GetName() string {
	return ProviderName
}

func (o *OpenBLiveClient) Config() *liveroom.LiveRoomConfig {
	return &o.cfg
}

func (o *OpenBLiveClient) Connect() error {
	err := o.openbliveClient.Start()
	if err != nil {
		return err
	}
	o.conn = o.openbliveClient.GetLongConn()
	o.conn.OnDanmu(o.danmuHandler)
	o.conn.OnDisconnect(o.disconnectHandler)
	e := o.conn.EstablishConnection(context.Background())
	if e == nil {
		if o.onStatusChange != nil {
			o.onStatusChange(true)
		}
	}
	return e
}

func (o *OpenBLiveClient) Disconnect() error {
	_ = o.conn.CloseConnection()
	if o.onStatusChange != nil {
		o.onStatusChange(false)
	}
	return o.openbliveClient.End()
}

func (o *OpenBLiveClient) OnStatusChange(f func(connected bool)) {
	o.onStatusChange = f
}

func (o *OpenBLiveClient) OnDisconnect(f func(liveroom liveroom.LiveRoom)) {
	o.onDisconnect = f
}

func (o *OpenBLiveClient) OnMessage(f func(msg *liveroom.Message)) {
	o.onMessage = f
}
