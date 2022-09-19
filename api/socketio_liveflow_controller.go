package api

import (
	"io"
)

// read from unix socket file and send on livestream web sockets
func (api *APIServer) _readAndForward(r io.Reader) {
	buf := make([]byte, 100*1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			continue
		}
		// log.Println("Client got:", string(buf[0:n]))
		jsonStr := string(buf[0:n])
		// for k, v := range api.WebSocketClients.WSClient {
		// 	log.Println(">>>>>>>>>>>>>", k)
		// 	if v.Kind == "liveflow" {
		// 		v.Conn.Emit("live-flow", jsonStr)
		// 	}
		// }
		// fmt.Println(jsonStr)
		if api.apiSocketServer != nil {
			api.apiSocketServer.BroadcastToRoom(
				"/",
				"liveflow",
				"live-flow",
				jsonStr,
			)
		}
	}
}
