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

func (sampler *BaseSampler) GenerateSamples() {
	panic("oh no")
}

func (sampler *BaseSampler) ShuffleIndexes() {
	// Setup shuffled indexes
	indexes := make([]int, sampler.NumberOfSamples)
	for p := 0; p < sampler.NumberOfSets; p++ {
		for i := 0; i < sampler.NumberOfSamples; i++ {
			indexes[i] = i
		}
		rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })
		sampler.shuffledIndexes = append(sampler.shuffledIndexes, indexes...)
	}
}

func (sampler *BaseSampler) SampleUnitSquare() geometry.Point2D {
	if sampler.Count%sampler.NumberOfSamples == 0 {
		sampler.jump = (rand.Int() % sampler.NumberOfSets) * sampler.NumberOfSamples
	}
	sampler.Count++
	return sampler.Samples[sampler.jump+sampler.shuffledIndexes[sampler.jump+(sampler.Count%sampler.NumberOfSamples)]]
}

type JitteredPointSampler struct {
	BaseSampler
}

func (sampler *JitteredPointSampler) GenerateSamples() {
	// Setup jittered samples
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

type RegularSampler struct {
	BaseSampler
}

func (sampler *RegularSampler) GenerateSamples() {
	sampler.Samples = make([]geometry.Point2D, sampler.NumberOfSamples*sampler.NumberOfSets)
	for i, _ := range sampler.Samples {
		sampler.Samples[i] = geometry.Point2D{0.5, 0.5}
	}
}
