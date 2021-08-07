package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
	"math/rand"
)

type Sampler interface {
	GenerateSamples()
	SampleUnitSquare() geometry.Point2D
	SampleUnitCircle() geometry.Point2D
}

type BaseSampler struct {
	NumberOfSamples int
	NumberOfSets    int
	Samples         []geometry.Point2D
	CircleSamples   []geometry.Point2D
	shuffledIndexes []int
	Count           int
	jump            int
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

func (sampler *BaseSampler) MapSamplesToCircle() {
	sampler.CircleSamples = make([]geometry.Point2D, len(sampler.Samples))
	var point geometry.Point2D
	var r, phi float64
	for i, _ := range sampler.CircleSamples {
		point = geometry.Point2D{2*sampler.Samples[i][0] - 1, 2*sampler.Samples[i][0] - 1}

		if point[0] > -1*point[1] {
			if point[0] > point[1] {
				r = point[0]
				phi = point[1] / point[0]
			} else {
				r = point[1]
				phi = 1 - point[0]/point[1]
			}
		} else {
			if point[0] > point[1] {
				r = -1 * point[0]
				phi = 4 + point[1]/point[0]
			} else {
				r = -1 * point[1]
				if point[1] != 0 {
					phi = 6 - point[0]/point[1]
				} else {
					phi = 0
				}
			}
		}

		sampler.CircleSamples[i][0] = r * math.Cos(phi)
		sampler.CircleSamples[i][1] = r * math.Sin(phi)
	}
}

func (sampler *BaseSampler) SampleUnitCircle() geometry.Point2D {
	if sampler.Count%sampler.NumberOfSamples == 0 {
		sampler.jump = (rand.Int() % sampler.NumberOfSets) * sampler.NumberOfSamples
	}
	sampler.Count++
	return sampler.CircleSamples[sampler.jump+sampler.shuffledIndexes[sampler.jump+(sampler.Count%sampler.NumberOfSamples)]]
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

func CreateJitteredSampler(numberOfSamples, numberOfSets int) *JitteredPointSampler {
	sampler := &JitteredPointSampler{BaseSampler: BaseSampler{NumberOfSamples: numberOfSamples, NumberOfSets: numberOfSets}}

	sampler.ShuffleIndexes()
	sampler.GenerateSamples()
	sampler.MapSamplesToCircle()
	return sampler
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
