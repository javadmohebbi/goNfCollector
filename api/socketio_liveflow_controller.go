package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/fwsock"
	"github.com/sirupsen/logrus"
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

		var reqResp fwsock.ClientServerReqResp
		var metrics []common.Metric
		err = json.Unmarshal([]byte(jsonStr), &reqResp)
		if err != nil {
			api.d.Verbose(fmt.Sprintf("error in unmarshaling json in socket-io live-flow: %s", err.Error()), logrus.ErrorLevel)
			// due to error, no filter will be applied
			// no matter its enabled or disabled
			continue
		} else {
			// convert ReqResp to metrics
			for _, p := range reqResp.Payload.([]interface{}) {
				b, err := json.Marshal(p)
				if err == nil {
					var _m common.Metric
					err := json.Unmarshal(b, &_m)
					if err == nil {
						metrics = append(metrics, _m)
					}
				}
			}
		}

		// filter metrics to send to UI
		metricsToSend := api.filteLiveFlow.FilterMetrics(&metrics)

		// convert filtered metrics to json string
		b, err := json.Marshal(metricsToSend)
		if err != nil {
			api.d.Verbose(fmt.Sprintf("error is marshalin struct to json in socket-io live-flow: %s", err.Error()), logrus.ErrorLevel)
			continue
		}

		if api.apiSocketServer != nil {
			api.apiSocketServer.BroadcastToRoom(
				"/",
				"liveflow",
				"live-flow",
				string(b),
				// jsonStr,
			)
		}
	}
}
