package liveroom

var _providers map[string]ILiveRoomProvider = make(map[string]ILiveRoomProvider)

func RegisterProvider(provider ILiveRoomProvider) {
	if _, ok := _providers[provider.GetName()]; ok {
		panic("provider " + provider.GetName() + " already exists")
		return
	}
	_providers[provider.GetName()] = provider
}

func GetProvider(name string) (ILiveRoomProvider, bool) {
	provider, ok := _providers[name]
	return provider, ok
}

func ListAvailableProviders() []string {
	var names []string
	for name := range _providers {
		names = append(names, name)
	}
	return names
}
