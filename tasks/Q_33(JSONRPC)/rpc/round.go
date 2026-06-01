package rpc

import "math"

func round(val float64) float64 {
	mu.Lock()
	p := precision
	mu.Unlock()
	pow := math.Pow(10, float64(p))
	return math.Round(val*pow) / pow
}
