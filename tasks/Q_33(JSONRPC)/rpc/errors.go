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

func errParse() *RPCError {
	return &RPCError{Code: codeParseError, Message: "Ошибка парсера"}
}

func errInvalidRequest() *RPCError {
	return &RPCError{Code: codeInvalidRequest, Message: "Невалидный запрос"}
}

func errMethodNotFound() *RPCError {
	return &RPCError{Code: codeMethodNotFound, Message: "Метод не найден"}
}

func errInvalidParams(detail string) *RPCError {
	e := &RPCError{Code: codeInvalidParams, Message: "Невалидные парметры"}
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
