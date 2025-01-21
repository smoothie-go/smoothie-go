package weighting

import "math"

func Gaussian(frames int, mean, stdDev float64, bound [2]float64) ([]float64, error) {
	if err := warnBound(bound, "Gaussian"); err != nil {
		return nil, err
	}
	xAxis := ScaleRange(frames, bound[0], bound[1])
	weights := make([]float64, frames)
	for i, x := range xAxis {
		weights[i] = math.Exp(-math.Pow(x-mean, 2) / (2 * math.Pow(stdDev, 2)))
	}
	return Normalise(weights), nil
}

func GaussianSym(frames int, stdDev float64, bound [2]float64) ([]float64, error) {
	if err := warnBound(bound, "GaussianSym"); err != nil {
		return nil, err
	}

	maxAbs := math.Max(math.Abs(bound[0]), math.Abs(bound[1]))
	return Gaussian(frames, 0, stdDev, [2]float64{-maxAbs, maxAbs})
}

func ScaleRange(n int, start, end float64) []float64 {
	if n <= 1 {
		return []float64{start}
	}
	step := (end - start) / float64(n-1)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = start + step*float64(i)
	}
	return result
}
