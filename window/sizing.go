package window

import (
	"math"
)

type PaneSplitDir int

const (
	PaneSplitVertical PaneSplitDir = iota
	PaneSplitHorizontal
)

type Sizing func(w, h int) (int, int)

type SideSizing func(x int) int

func FixedSideSizing(x int) SideSizing {
	return func(int) int {
		return x
	}
}

func RatioSideSizing(r float64) SideSizing {
	return func(x int) int {
		return int(math.Round(float64(x) * r))
	}
}

func CompSizing(ws SideSizing, hs SideSizing) Sizing {
	return func(w, h int) (int, int) {
		return ws(w), hs(h)
	}
}

func HalfSizing(dir PaneSplitDir) Sizing {
	switch dir {
	case PaneSplitHorizontal:
		return CompSizing(RatioSideSizing(0.5), RatioSideSizing(1))
	case PaneSplitVertical:
		return CompSizing(RatioSideSizing(1), RatioSideSizing(0.5))
	}
	return nil
}
