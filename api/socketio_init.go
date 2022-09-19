package api

// func (api *APIServer) InitSocketIO(r *mux.Router) {
// 	if api.apiSocketServer == nil {
// 		/**
// 		*
// 		* BEGIN SOCKET IO
// 		*
// 		 */
// 		pt := polling.Default
// 		wt := websocket.Default
// 		wt.CheckOrigin = func(req *http.Request) bool {
// 			return true
// 		}

// 		server := socketio.NewServer(&engineio.Options{
// 			Transports: []transport.Transport{
// 				pt,
// 				wt,
// 			},
// 		})

// 		go server.Serve()

// 		api.apiSocketServer = server

// 		api.apiSocketServer.OnConnect("/", func(s socketio.Conn) error {
// 			s.SetContext("")
// 			fmt.Println("[WS]connected:", s.ID())

// 			// api.WebSocketClients.Add(s.ID(), "UNKNOWN", s)
// 			// log.Println(api.WebSocketClients)
// 			return nil
// 		})

// 		api.apiSocketServer.OnEvent("/", "join", func(s socketio.Conn, msg string) {
// 			s.Join(msg)
// 			fmt.Println("[WS]Join:", s.ID(), "->", msg)
// 			api.WebSocketClients.Add(s.ID(), msg, s)
// 			// s.Emit("pong", "joined ip2l")
// 			s.Emit("pong", "have "+msg)

// 		})

// 		api.apiSocketServer.OnDisconnect("/", func(s socketio.Conn, msg string) {
// 			fmt.Println("[WS]closed!!!!", msg)
// 			api.WebSocketClients.Remove(s.ID())
// 			return
// 		})

// 	}

// 	r.Handle("/socket.io/", socketioMiddleware(api.apiSocketServer))
// }
