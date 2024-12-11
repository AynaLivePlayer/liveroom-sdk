package main

import (
	"fmt"
	"github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/liveroom-sdk/provider/webdm"
	"os"
	"os/signal"
	"time"
)

const apiServer = "https://ayliveplayer.scene.aynakeya.com:10443"

//const apiServer = "https://api.biliaudiobot.com/"

// const apiServer = "https://ayliveplayer.scene.aynakeya.com:10443"

func main() {
	provider := webdm.NewWebDanmuClientProvider(apiServer)
	room, err := provider.CreateLiveRoom(liveroom.LiveRoom{
		Room:     "3819533",
		Provider: webdm.ProviderName,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
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
	fmt.Println("disconnect", room.Disconnect())
	time.Sleep(3 * time.Second)
	fmt.Println("Bye")
}
