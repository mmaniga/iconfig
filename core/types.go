package core

//ConfigInfo :
type ConfigInfo struct {
	Key   string      `json:"key"`
	Value string      `json:"value"`
	MType RequestType `json:"mType"`
}

//ConfigRequest :
type ConfigRequest struct {
	ClientID string      `json:"clientID"`
	Key      string      `json:"key"`
	MType    RequestType `json:"mType"`
}

//ConfigMap :
type ConfigMap struct {
	CMap map[string]string `json:"cmap"`
}

//RequestType :"Identifies different request types"
type RequestType int

const (
	//MessageSuccess :
	MessageSuccess RequestType = iota //0

	//RequestConfig :"Client sends this message for getting request"
	RequestConfig

	//ConfigSend :
	ConfigSend
)
