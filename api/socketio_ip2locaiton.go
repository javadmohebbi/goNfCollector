package api

import (
	"fmt"
	"time"
)

// sends to room every 1 seconds
func (api *APIServer) _ip2loc() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// log.Println("Room Len. (ip2l):", api.apiSocketServer.RoomLen("/", "ip2l"))
				api.apiSocketServer.BroadcastToRoom(
					"/",
					"ip2l",
					"pong",
					fmt.Sprintf("%d", time.Now().Unix()),
				)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}
