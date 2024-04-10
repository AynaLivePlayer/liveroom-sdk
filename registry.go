package liveroom

var _providers map[string]ILiveRoomProvider = make(map[string]ILiveRoomProvider)

func RegisterProvider(name string, provider ILiveRoomProvider) {
	if _, ok := _providers[name]; ok {
		panic("provider " + name + " already exists")
		return
	}
	_providers[name] = provider
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
