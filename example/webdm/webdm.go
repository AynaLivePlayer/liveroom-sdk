package main

import (
	"fmt"
	"liveroom"
	"liveroom/provider/webdm"
	"os"
	"os/signal"
	"time"
)

const apiServer = "http://0.0.0.0:9090"

func main() {
	provider := webdm.NewWebDanmuClientProvider(apiServer)
	room, _ := provider(liveroom.LiveRoomConfig{
		Room:     "7777",
		Provider: webdm.ProviderName,
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
		func(liveroom liveroom.LiveRoom) {
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
