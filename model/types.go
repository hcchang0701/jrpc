package model

// type SingleReqObj struct {
// 	Headers map[string][]string `yaml:"headers"`
// 	Data RpcPayload `yaml:"request"`
// }

type Request struct {
	Url    string              `yaml:"url"`
	Header map[string][]string `yaml:"header"`
	Body   []RpcPayload        `yaml:"body"`
}

type RpcPayload struct {
	Version string        `json:"jsonrpc" yaml:"jsonrpc"`
	Method  string        `json:"method" yaml:"method"`
	Params  []interface{} `json:"params" yaml:"params"`
	ID      interface{}   `json:"id,omitempty" yaml:"id"`
}

type RespPayload struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   ErrPayload  `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type ErrPayload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
