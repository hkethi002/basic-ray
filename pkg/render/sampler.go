package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
	"math/rand"
)

type Sampler interface {
	GenerateSamples()
	SampleUnitSquare() geometry.Point2D
}

type BaseSampler struct {
	NumberOfSamples int
	NumberOfSets    int
	Samples         []geometry.Point2D
	shuffledIndexes []int
	Count           int
	jump            int
}

func (sampler *BaseSampler) SampleUnitSquare() geometry.Point2D {
	sampler.Count++
	return sampler.Samples[sampler.Count%(sampler.NumberOfSamples*sampler.NumberOfSets)]
}

type JitteredPointSampler struct {
	BaseSampler
}

func (sampler *JitteredPointSampler) GenerateSamples() {
	n := (int)(math.Sqrt((float64)(sampler.NumberOfSamples)))
	sampler.Samples = make([]geometry.Point2D, sampler.NumberOfSamples*sampler.NumberOfSets)
	count := 0
	for p := 0; p < sampler.NumberOfSets; p++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				x := ((float64)(k) + rand.Float64()) / (float64)(n)
				y := ((float64)(j) + rand.Float64()) / (float64)(n)
				point := geometry.Point2D{x, y}
				sampler.Samples[count] = point
				count++
			}
		}
	}
}
