package fwsock

// struct for making standard message between clients
// and server
type ClientServerReqResp struct {

	// for knowing it's himself or not
	ItSelf bool `json:"itSelf,omitempty"`

	// API - if it's a api client
	API bool `json:"api,omitempty"`

	// to find if it's a collector server
	Collector bool `json:"collector,omitempty"`

	// request id
	RequestID string `json:"reqId,omitempty"`

	// PayLoad
	Payload interface{} `json:"payload,omitempty"`

	// command between server and client
	Command Command `json:"cmd"`
}
