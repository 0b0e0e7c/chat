package common

import "golang.org/x/exp/constraints"

func LowHigh[E constraints.Ordered](a, b E) (lower, higher E) {
	if a > b {
		return b, a
	}
	return a, b
}
