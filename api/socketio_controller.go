package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/gorilla/mux"
)

func (api *APIServer) IP2LocationUpdate(r *mux.Router) {

	if api.apiSocketServer == nil {
		/**
		*
		* BEGIN SOCKET IO
		*
		 */
		pt := polling.Default
		wt := websocket.Default
		wt.CheckOrigin = func(req *http.Request) bool {
			return true
		}

		server := socketio.NewServer(&engineio.Options{
			Transports: []transport.Transport{
				pt,
				wt,
			},
		})

		server.OnConnect("/", func(s socketio.Conn) error {
			s.SetContext("")
			fmt.Println("connected:", s.ID())
			return nil
		})

		server.OnEvent("/", "getIP2L", func(s socketio.Conn, msg string) {
			ss := fmt.Sprintf("onEvent getIP2L: %v %v", msg, time.Now().Unix())
			log.Println(ss)
			s.Emit("pong", ss)
		})

		server.OnDisconnect("/", func(s socketio.Conn, msg string) {
			fmt.Println("closed!!!!", msg)
		})

		go server.Serve()

		api.apiSocketServer = server
		// defer server.Close() close it in =>  api.go

		/**
		*
		* E N D SOCKET IO
		*
		 */
	}

	r.Handle("/socket.io/", socketioMiddleware(api.apiSocketServer))

}
