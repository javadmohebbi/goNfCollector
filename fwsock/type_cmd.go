package fwsock

import (
	"encoding/json"
	"errors"
)

type Command int

const (

	// API/Collector first initialization
	CMD_INIT Command = iota + 0x1000

	// Collector exported data to socket
	CMD_EXPORTED
)

func StringToJSONClientServerReqResp(reqJson string) (ClientServerReqResp, error) {
	var jsonReq ClientServerReqResp
	err := json.Unmarshal([]byte(reqJson), &jsonReq)
	if err != nil {
		return ClientServerReqResp{}, err
	}

	if jsonReq.RequestID == "" {
		return ClientServerReqResp{}, errors.New("reqId could not be empty")
	}

	return jsonReq, nil
}

func (req ClientServerReqResp) JSONToStringClientServerReqResp() (string, error) {
	// var jsonReq ClientServerReqResp
	// err := json.Unmarshal([]byte(reqJson), &jsonReq)
	// if err != nil {
	// 	return ClientServerReqResp{}, err
	// }

	// return jsonReq, nil
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (req ClientServerReqResp) JSONToByteClientServerReqResp() ([]byte, error) {
	// var jsonReq ClientServerReqResp
	// err := json.Unmarshal([]byte(reqJson), &jsonReq)
	// if err != nil {
	// 	return ClientServerReqResp{}, err
	// }

	// return jsonReq, nil
	b, err := json.Marshal(req)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
