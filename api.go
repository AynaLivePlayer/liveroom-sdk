package liveroom

import "errors"

func CreateLiveRoom(config LiveRoom) (ILiveRoom, error) {
	provider, ok := GetProvider(config.Provider)
	if !ok {
		return nil, errors.New("provider not found")
	}
	return provider.CreateLiveRoom(config)
}
