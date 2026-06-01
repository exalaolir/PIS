package rpc

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
)

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func errParse() *RPCError {
	return &RPCError{Code: codeParseError, Message: "Parse error"}
}

func errInvalidRequest() *RPCError {
	return &RPCError{Code: codeInvalidRequest, Message: "Invalid Request"}
}

func errMethodNotFound() *RPCError {
	return &RPCError{Code: codeMethodNotFound, Message: "Method not found"}
}

func errInvalidParams(detail string) *RPCError {
	e := &RPCError{Code: codeInvalidParams, Message: "Invalid params"}
	if detail != "" {
		e.Data = detail
	}
	return e
}

func writeRPC(w http.ResponseWriter, resp *RPCResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Ошибка кодирования ответа:", err)
	}
}

func errorResponse(id interface{}, rpcErr *RPCError) *RPCResponse {
	return &RPCResponse{
		Jsonrpc: "2.0",
		Error:   rpcErr,
		ID:      id,
	}
}

func successResponse(id interface{}, result interface{}) *RPCResponse {
	return &RPCResponse{
		Jsonrpc: "2.0",
		Result:  result,
		ID:      id,
	}
}
