/*
 * https://github.com/couleur-tweak-tips/smoothie-rs/blob/main/target/scripts/weighting.py
 * Ported to Golang
 */

package weighting

import (
	"errors"
	"math"
)

var wizardry = false

func EnableWizardry() {
	wizardry = true
}

func Ascending(frames int) []float64 {
	weights := make([]float64, frames)
	for i := 1; i <= frames; i++ {
		weights[i-1] = float64(i)
	}
	return Normalise(weights)
}

func Descending(frames int) []float64 {
	weights := make([]float64, frames)
	for i := 0; i < frames; i++ {
		weights[i] = float64(frames - i)
	}
	return Normalise(weights)
}

func Equal(frames int) []float64 {
	weights := make([]float64, frames)
	for i := range weights {
		weights[i] = 1.0 / float64(frames)
	}
	return weights
}

func Pyramid(frames int) []float64 {
	half := float64(frames-1) / 2
	weights := make([]float64, frames)
	for i := 0; i < frames; i++ {
		weights[i] = half - math.Abs(float64(i)-half) + 1
	}
	return Normalise(weights)
}

func Vegas(input_fps int, output_fps int, blur_amount float64) []float64 {
	blendFactor := int(float64(input_fps) / float64(output_fps) * blur_amount)
	nWeights := blendFactor + (1 - (blendFactor % 2))

	weights := make([]float64, nWeights)

	if (blendFactor % 2) == 0 {
		for i := 1; i < nWeights; i++ {
			weights[i] = 2
		}
	} else {
		for i := 0; i < nWeights; i++ {
			weights[i] = 1
		}
	}

	return Normalise(weights)
}

func Divide(frames int, weights []float64) []float64 {
	stretched := make([]float64, frames)
	r := ScaleRange(frames, 0, float64(len(weights))-0.1)
	for i := 0; i < frames; i++ {
		stretched[i] = weights[int(r[i])]
	}
	return Normalise(stretched)
}

func warnBound(bound [2]float64, funcName string) error {
	if bound[0] == bound[1] {
		return errors.New(funcName + ": bound values must differ")
	}
	return nil
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}
