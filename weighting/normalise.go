package weighting

import "math"

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func Normalise(weights []float64) []float64 {
	if len(weights) == 0 {
		return weights
	}
	// Remove negative weights if wizardry is not enabled
	if !wizardry {
		minWeight := min(weights)
		if minWeight < 0 {
			for i := range weights {
				weights[i] += math.Abs(minWeight) + 1
			}
		}
	}
	sumVal := sum(weights)
	if sumVal == 0 {
		return weights
	}
	for i := range weights {
		weights[i] /= sumVal
	}
	return weights
}
