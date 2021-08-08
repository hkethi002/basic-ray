package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
	"math/rand"
)

const PI = 3.14159265358

type Sampler interface {
	GenerateSamples()
	SampleUnitSquare() geometry.Point2D
	SampleUnitCircle() geometry.Point2D
}

type BaseSampler struct {
	NumberOfSamples   int
	NumberOfSets      int
	Samples           []geometry.Point2D
	CircleSamples     []geometry.Point2D
	HemisphereSamples []geometry.Point
	shuffledIndexes   []int
	Count             int
	jump              int
	CircleCount       int
	circleJump        int
	HemisphereCount   int
	HemisphereJump    int
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

func (sampler *BaseSampler) MapSamplesToHemisphere(e float64) {
	sampler.HemisphereSamples = make([]geometry.Point, len(sampler.Samples))
	for i, _ := range sampler.HemisphereSamples {
		cos_phi := math.Cos(2.0 * PI * sampler.Samples[i][0])
		sin_phi := math.Sin(2.0 * PI * sampler.Samples[i][0])
		cos_theta := math.Pow((1.0 - sampler.Samples[i][1]), 1/(e+1.0))
		sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)
		sampler.HemisphereSamples[i][0] = sin_theta * cos_phi
		sampler.HemisphereSamples[i][1] = sin_theta * sin_phi
		sampler.HemisphereSamples[i][2] = cos_theta
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
	if sampler.CircleCount%sampler.NumberOfSamples == 0 {
		sampler.jump = (rand.Int() % sampler.NumberOfSets) * sampler.NumberOfSamples
	}
	sampler.CircleCount++
	return sampler.CircleSamples[sampler.jump+sampler.shuffledIndexes[sampler.jump+(sampler.CircleCount%sampler.NumberOfSamples)]]
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

func CreateJitteredSampler(numberOfSamples, numberOfSets int, e float64) *JitteredPointSampler {
	sampler := &JitteredPointSampler{BaseSampler: BaseSampler{NumberOfSamples: numberOfSamples, NumberOfSets: numberOfSets}}

	sampler.ShuffleIndexes()
	sampler.GenerateSamples()
	sampler.MapSamplesToCircle()
	sampler.MapSamplesToHemisphere(e)
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
