package liveroom

type LiveRoomConfig struct {
	Provider string `json:"provider"` // Provider is the name of the live room provider
	Room     string `json:"room"`     // RoomID is the unique identifier of the live room
}

func (l *LiveRoomConfig) Identifier() string {
	return l.Provider + "_" + l.Room
}

type ILiveRoomProvider interface {
	GetName() string
	GetDescription() string
	CreateLiveRoom(cfg LiveRoomConfig) (LiveRoom, error)
}

type LiveRoomProvider struct {
	Name        string
	Description string
	Func        func(cfg LiveRoomConfig) (LiveRoom, error)
}

func (l *LiveRoomProvider) GetName() string {
	return l.Name
}

func (l *LiveRoomProvider) GetDescription() string {
	return l.Name
}

func (l *LiveRoomProvider) CreateLiveRoom(cfg LiveRoomConfig) (LiveRoom, error) {
	return l.Func(cfg)
}

type UserMedal struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	RoomID string `json:"room_id"`
}

const (
	PrivilegeNone = iota
	PrivilegeBasic
	PrivilegeAdvanced
	PrivilegeUltimate
)

type User struct {
	Uid       string
	Username  string
	Admin     bool
	Privilege int
	Medal     UserMedal
}

type Message struct {
	User    User
	Message string
}

type LiveRoom interface {
	GetName() string         // should return the name of the provider
	Config() *LiveRoomConfig // should return mutable model (not a copy)
	Connect() error
	Disconnect() error
	OnDisconnect(func(liveroom LiveRoom))
	OnStatusChange(func(connected bool))
	OnMessage(func(msg *Message))
}
