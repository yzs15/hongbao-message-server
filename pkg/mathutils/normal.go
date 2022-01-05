package mathutils

import (
	"math"
)

func NormalFunc(exp int64, var_ float64, x int64) float64 {
	return math.Exp(-1*math.Pow(float64(x-exp), 2.0)/(2.0*var_)) / math.Sqrt(2.0*math.Pi*var_)
}

func CalVariance(exp int64, tarP float64) float64 {
	left := 0.0
	right := 1e10

	var var_ float64
	for {
		var_ = (left + right) / 2
		curP := NormalFunc(exp, var_, exp)

		if math.Abs(curP-tarP) < 1e-6 {
			break
		}

		if curP < tarP {
			right = var_
		} else {
			left = var_
		}
	}

	return var_
}
