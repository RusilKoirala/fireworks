package main

import "math"

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func hash(x, y int) int {
	h := uint32(x)*73856093 ^ uint32(y)*19349663
	h ^= h >> 13
	h *= 0x5bd1e995
	return int(h >> 8)
}
