package fwsock

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/goNfCollector/configurations"
	"github.com/sirupsen/logrus"
)

// make a  socket listener
func (fws *FwSock) MakeSocketListener() (configurations.ErrorCodes, error) {
	err := os.RemoveAll(fws.SocketAddr)
	if err != nil {
		fws.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)", configurations.ERROR_CAN_T_REMOVE_SOCKET_ADDRESS.Int(), configurations.ERROR_CAN_T_REMOVE_SOCKET_ADDRESS, err), logrus.ErrorLevel)
		return configurations.ERROR_CAN_T_REMOVE_SOCKET_ADDRESS, err
	}

	ls, err := net.Listen("unix", fws.SocketAddr)
	if err != nil {
		fws.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)", configurations.ERROR_CAN_T_LISTEN_ON_LINUX_SOCKET_ADDRESS.Int(), configurations.ERROR_CAN_T_LISTEN_ON_LINUX_SOCKET_ADDRESS, err), logrus.ErrorLevel)
		return configurations.ERROR_CAN_T_LISTEN_ON_LINUX_SOCKET_ADDRESS, err
	}

	fws.lstn = ls

	return configurations.NO_ERROR, nil
}

func (fws *FwSock) Close() {
	fws.lstn.Close()
}

func (fws *FwSock) SetChann(ch chan os.Signal) {
	fws.Ch = ch
}

// Listen and accept connections
// it should be called via goroutines inside
// collector Listen function
func (fws *FwSock) Accept() {

	if fws.socketClientList == nil {
		fws.socketClientList = &socketClientList{}
	}

	for {
		conn, err := fws.lstn.Accept()
		if err != nil {
			select {
			case <-fws.Ch:
				log.Println("singnal closed recvd!")
				fws.lstn.Close()
				return
			default:
				log.Println("err in connection: ", fws.lstn.Addr().String())
				continue
			}
		} else {
			// handle requests
			go func() {
				_c := &netSockClient{
					Conn: conn,
				}
				// _c.IdentifySockClient()
				fws.socketClientList.AddSockClient(_c)

				_c.HandleSockConnection(fws.socketClientList)
			}()
		}
	}
}
