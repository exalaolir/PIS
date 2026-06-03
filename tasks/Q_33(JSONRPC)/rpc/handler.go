package rpc

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleRPC(req RPCRequest) *RPCResponse {
	log.Printf("метод=%s параметры=%s id=%v\n", req.Method, string(req.Params), req.ID)

	if req.ID == nil {
		invokeMethod(req)
		return nil
	}

	if req.Jsonrpc != "2.0" {
		return errorResponse(req.ID, errInvalidRequest())
	}
	if req.Method == "" {
		return errorResponse(req.ID, errInvalidRequest())
	}

	result, rpcErr := invokeMethod(req)
	if rpcErr != nil {
		log.Printf("ошибка=%v\n", rpcErr)
		return errorResponse(req.ID, rpcErr)
	}

	log.Printf("результат=%v\n", result)
	return successResponse(req.ID, result)
}

func invokeMethod(req RPCRequest) (interface{}, *RPCError) {
	switch req.Method {
	case "sum":
		return sum(req.Params)
	case "sub":
		return sub(req.Params)
	case "mul":
		return mul(req.Params)
	case "div":
		return div(req.Params)
	case "pre":
		return pre(req.Params)
	default:
		return nil, errMethodNotFound()
	}
}

func RpcHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var raw json.RawMessage

	if err := decoder.Decode(&raw); err != nil {
		log.Println("Некорректный запрос:", err)
		writeRPC(w, errorResponse(nil, errParse()))
		r.Body.Close()
		return
	}

	if decoder.More() {
		log.Println("В запросе содержится несколько JSON-объектов")
		writeRPC(w, errorResponse(nil, errInvalidRequest()))
		r.Body.Close()
		return
	}

	if len(raw) > 0 && raw[0] == '[' {
		var reqs []RPCRequest
		if err := json.Unmarshal(raw, &reqs); err != nil {
			log.Println("Некорректный пакет:", err)
			writeRPC(w, errorResponse(nil, errInvalidRequest()))
			r.Body.Close()
			return
		}

		var res []RPCResponse

		for _, req := range reqs {
			rpcRes := handleRPC(req)
			if rpcRes != nil {
				res = append(res, *rpcRes)
			}
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			r.Body.Close()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println("Ошибка пакетной обработки при кодировании:", err)
		}
		r.Body.Close()
		return
	}

	var req RPCRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		log.Println("Некорректный единичный запрос:", err)
		writeRPC(w, errorResponse(nil, errInvalidRequest()))
		r.Body.Close()
		return
	}

	rpcRes := handleRPC(req)

	if rpcRes == nil {
		w.WriteHeader(http.StatusNoContent)
		r.Body.Close()
		return
	}

	writeRPC(w, rpcRes)
	r.Body.Close()
}
