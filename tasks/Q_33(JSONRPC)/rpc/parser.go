package rpc
import (
	"encoding/json"
	"errors"
)

func parseXY(params json.RawMessage) (float64, float64, error) {
	var arr []float64
	if err := json.Unmarshal(params, &arr); err == nil && len(arr) == 2 {
		return arr[0], arr[1], nil
	}

	var obj XYMap
	if err := json.Unmarshal(params, &obj); err == nil {
		return obj.X, obj.Y, nil
	}

	return 0, 0, errors.New("неверный формат параметров")
}

func parseN(params json.RawMessage) (int, error) {
	var obj NMap
	if err := json.Unmarshal(params, &obj); err != nil {
		return 0, errors.New("неверный формат параметров")
	}
	return obj.N, nil
}