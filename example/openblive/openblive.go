package main

import (
	"fmt"
	"github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/liveroom-sdk/provider/openblive"
	"os"
	"os/signal"
	"time"
)

const apiServer = "http://0.0.0.0:9090"

func main() {
	provider := openblive.NewOpenBLiveClientProvider(apiServer, 1661006726438)
	room, _ := provider.CreateLiveRoom(liveroom.LiveRoom{
		Room:     "CSMVI59S9HC02",
		Provider: openblive.ProviderName,
	})
	room.OnMessage(func(msg *liveroom.Message) {
		fmt.Println(msg.User.Username, msg.User.Uid, msg.User.Medal.Name, msg.Message)
	})
	room.OnStatusChange(
		func(connected bool) {
			if connected {
				fmt.Println("Connected")
			} else {
				fmt.Println("Disconnected")
			}
		})
	room.OnDisconnect(
		func(liveroom liveroom.ILiveRoom) {
			fmt.Println("Disconnected AAAAAAAa")
		})
	fmt.Println("connect", room.Connect())
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	fmt.Println("disconnect", room.Disconnect() == nil)
	time.Sleep(3 * time.Second)
	fmt.Println("Bye")
}
