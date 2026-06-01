package rpc

import (
	"encoding/json"
	"sync"
)

var (
	precision = 0
	mu        sync.Mutex
)

func sum(params json.RawMessage) (interface{}, *RPCError) {
	x, y, err := parseXY(params)
	if err != nil {
		return nil, errInvalidParams(err.Error())
	}
	return round(x + y), nil
}

func sub(params json.RawMessage) (interface{}, *RPCError) {
	x, y, err := parseXY(params)
	if err != nil {
		return nil, errInvalidParams(err.Error())
	}
	return round(x - y), nil
}

func mul(params json.RawMessage) (interface{}, *RPCError) {
	x, y, err := parseXY(params)
	if err != nil {
		return nil, errInvalidParams(err.Error())
	}
	return round(x * y), nil
}

func div(params json.RawMessage) (interface{}, *RPCError) {
	x, y, err := parseXY(params)
	if err != nil {
		return nil, errInvalidParams(err.Error())
	}
	if y == 0 {
		return nil, errInvalidParams("деление на ноль")
	}
	return round(x / y), nil
}

func pre(params json.RawMessage) (interface{}, *RPCError) {
	n, err := parseN(params)
	if err != nil {
		return nil, errInvalidParams(err.Error())
	}
	mu.Lock()
	precision = n
	mu.Unlock()
	return "ok", nil
}
