package openblive

import (
	"context"
	"errors"
	"github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/liveroom-sdk/utils"
	openblive "github.com/aynakeya/open-bilibili-live"
	"strconv"
)

const ProviderName = "openblive"

type OpenBLiveClient struct {
	cfg             liveroom.LiveRoom
	admins          map[string]int
	openbliveClient *openblive.BLiveClient
	conn            openblive.BLiveLongConnection
	onMessage       func(msg *liveroom.Message)
	onDisconnect    func(liveroom liveroom.ILiveRoom)
	onStatusChange  func(connected bool)
}

func NewOpenBLiveClientProvider(apiServer string, appId int64) liveroom.ILiveRoomProvider {
	return &liveroom.LiveRoomProvider{
		Name:        ProviderName,
		Description: "open bilibili live protocol. enter client key to connect.",
		Func: func(cfg liveroom.LiveRoom) (liveroom.ILiveRoom, error) {
			if cfg.Provider != ProviderName {
				return nil, errors.New("invalid provider name")
			}
			return &OpenBLiveClient{
				cfg:             cfg,
				openbliveClient: openblive.NewBliveClient(appId, cfg.Room, newRemoteApiClient(apiServer)),
			}, nil
		},
	}
}

func (o *OpenBLiveClient) danmuHandler(data openblive.DanmakuData) {
	msg := o.onMessage
	if msg == nil {
		return
	}
	roomId := strconv.Itoa(data.RoomID)
	if data.FansMedalName == "" {
		roomId = ""
	}
	_, isAdmin := o.admins[data.UName]
	msg(&liveroom.Message{
		User: liveroom.User{
			Uid:       data.OpenID,
			Username:  data.UName,
			Admin:     isAdmin, // not supported by open bilibili live
			Privilege: utils.BilibiliGuardLevelToPrivilege(data.GuardLevel),
			Medal: liveroom.UserMedal{
				Name:   data.FansMedalName,
				Level:  data.FansMedalLevel,
				RoomID: roomId,
			},
		},
		Message: data.Msg,
	})
}

func (o *OpenBLiveClient) disconnectHandler(conn openblive.BLiveLongConnection) {
	if x := o.onStatusChange; x != nil {
		x(false)
	}
	if x := o.onDisconnect; x != nil {
		x(o)
	}
}

func (o *OpenBLiveClient) GetName() string {
	return ProviderName
}

func (o *OpenBLiveClient) Config() *liveroom.LiveRoom {
	return &o.cfg
}

func (o *OpenBLiveClient) Connect() error {
	// if still running, return
	if o.openbliveClient.Status() {
		return nil
	}
	err := o.openbliveClient.Start()
	if err != nil {
		return err
	}
	// get admin list
	adminNames := getAdmins(o.openbliveClient.AppInfo.AnchorInfo.RoomID)
	o.admins = make(map[string]int)
	for _, name := range adminNames {
		o.admins[name] = 1
	}
	o.conn = o.openbliveClient.GetLongConn()
	o.conn.OnDanmu(o.danmuHandler)
	o.conn.OnDisconnect(o.disconnectHandler)
	e := o.conn.EstablishConnection(context.Background())
	if e == nil {
		if handler := o.onStatusChange; handler != nil {
			handler(true)
		}
	}
	return e
}

func (o *OpenBLiveClient) Disconnect() error {
	if o.conn == nil {
		return nil
	}
	_ = o.conn.CloseConnection()
	o.conn = nil
	if handler := o.onStatusChange; handler != nil {
		handler(false)
	}
	e := o.openbliveClient.End()
	return e
}

func (o *OpenBLiveClient) OnStatusChange(f func(connected bool)) {
	o.onStatusChange = f
}

func (o *OpenBLiveClient) OnDisconnect(f func(liveroom liveroom.ILiveRoom)) {
	o.onDisconnect = f
}

func (o *OpenBLiveClient) OnMessage(f func(msg *liveroom.Message)) {
	o.onMessage = f
}
