package liveroom

var _providers map[string]LiveRoomProvider = make(map[string]LiveRoomProvider)

func RegisterProvider(name string, provider LiveRoomProvider) {
	if _, ok := _providers[name]; ok {
		panic("provider " + name + " already exists")
		return
	}
	_providers[name] = provider
}

func GetProvider(name string) (LiveRoomProvider, bool) {
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
