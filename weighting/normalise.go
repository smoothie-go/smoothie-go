package weighting

import "math"

func min(values []float64) float64 {
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func Normalise(weights []float64) []float64 {
	// Remove negative weights if wizardry is not enabled
	if !wizardry {
		minWeight := min(weights)
		if minWeight < 0 {
			for i := range weights {
				weights[i] += math.Abs(minWeight) + 1
			}
		}
	}
	sum := sum(weights)
	for i := range weights {
		weights[i] /= sum
	}
	return weights
}
