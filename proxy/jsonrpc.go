package main

type JsonRPCRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Id      uint32      `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

type JsonRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JsonRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JsonRPCError `json:"error,omitempty"`
	Id      uint32        `json:"id"`
}
