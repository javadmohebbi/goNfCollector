package api

import (
	"fmt"
	"net/http"

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

		// server.OnEvent("/", "getUserInfo", func(s socketio.Conn, msg string) {

		// 	// Get URL to extract queries
		// 	u := s.URL()

		// 	// extract token from query
		// 	t := u.Query().Get("token")

		// 	// extract uid from query
		// 	uid := u.Query().Get("usrId")

		// 	// check if token is valid
		// 	// if common.CheckTokenInRedisDB(api.rdb, uid, t, strings.Split(s.RemoteAddr().String(), ":")[0]) {
		// 	// OK - user has access to get his information

		// 	// Prepare for TLS or non-TLS
		// 	isTLS, tlsConf, _ := auth.PrepareGRPCClientSideTLS(*api.conf)

		// 	// Dialing RPC
		// 	var conn *grpc.ClientConn
		// 	if isTLS {
		// 		conn, err = grpc.Dial(fmt.Sprintf("%v:%v", api.conf.Auth.Host, api.conf.Auth.Port), grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
		// 	} else {
		// 		conn, err = grpc.Dial(fmt.Sprintf("%v:%v", api.conf.Auth.Host, api.conf.Auth.Port), grpc.WithInsecure())
		// 	}
		// 	//Check for dial error
		// 	if err != nil {
		// 		api.debugger.Verbose(fmt.Sprintf("Can not connect to authentication server %v:%v (TLS:%v): error: %v", api.conf.Auth.Host, api.conf.Auth.Port, isTLS, err), logrus.ErrorLevel)
		// 		s.Emit("UserInfo", map[string]interface{}{
		// 			"forceLogout": true,
		// 		})
		// 	}
		// 	defer conn.Close()

		// 	iUID, _ := strconv.Atoi(uid)
		// 	// Prepare RPC request parameters
		// 	uir := auth.UserInfoRequest{
		// 		Id:        uint32(iUID), //userid
		// 		ReqUserId: uint32(iUID), // requestedUserID
		// 	}

		// 	// Create new instance of RPC service client
		// 	client := auth.NewAuthServiceClient(conn)

		// 	// Call RPC method
		// 	auth, err := client.UserInfo(context.Background(), &uir)

		// 	// connection error
		// 	if err != nil {
		// 		err, _ := common.ProtobufErrorHandler(err)
		// 		api.debugger.Verbose(fmt.Sprintf("authentication server %v:%v (TLS:%v) error: %v", api.conf.Auth.Host, api.conf.Auth.Port, isTLS, err), logrus.ErrorLevel)
		// 		s.Emit("UserInfo", map[string]interface{}{
		// 			"forceLogout": true,
		// 		})
		// 	}
		// 	jsMap := common.ProtobufToJSONMap(auth, false)
		// 	jsMap["adminMenu"] = common.AdminMenuBuilder(jsMap)
		// 	jsMap["forceLogout"] = common.CheckForceLogout(t, iUID)

		// 	// send back data to user
		// 	s.Emit("UserInfo", jsMap)

		// 	// } else {
		// 	// 	// Not valid
		// 	// 	s.Emit("UserInfo", map[string]interface{}{
		// 	// 		"forceLogout": true,
		// 	// 	})
		// 	// }

		// })

		// server.OnEvent("/getUserInfo", "msg", func(s socketio.Conn, msg string) string {
		// 	s.SetContext(msg)
		// 	return "recv " + msg
		// })

		// server.OnEvent("/", "bye", func(s socketio.Conn) string {
		// 	last := s.Context().(string)
		// 	s.Emit("bye", last)
		// 	s.Close()
		// 	return last
		// })
		// server.OnError("/", func(s socketio.Conn, e error) {
		// 	api.debugger.Verbose(fmt.Sprintf("socket error: %v", e), logrus.ErrorLevel)
		// })

		server.OnDisconnect("/", func(s socketio.Conn, msg string) {
			fmt.Println("closed", msg)
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

	r.Handle("/socket.io/", api.apiSocketServer)

}
