package main

import "github.com/jumper2017/melody/service_monitor"

func main() {
	serverAddr := []string{"localhost:2380", "localhost:2381", "localhost:2382"}
	//监控逻辑
	go service_monitor.WatchWithoutPassword(serverAddr, "/server/game/")

	//注册逻辑
	go service_monitor.RegisterWithoutPassword(serverAddr, "/server/game/hymahjong/", "hymahjong_1", "127.0.0.1:9999")
	go service_monitor.RegisterWithoutPassword(serverAddr, "/server/game/hymahjong/", "hymahjong_2", "127.0.0.1:9999")
	go service_monitor.RegisterWithoutPassword(serverAddr, "/server/game/wlmahjong/", "wlmahjong_1", "127.0.0.1:9999")
	go service_monitor.RegisterWithoutPassword(serverAddr, "/server/game/wlmahjong/", "wlmahjong_3", "127.0.0.1:9999")

	select {}
}
