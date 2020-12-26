package render

import (
	geometry "basic-ray/pkg/geometry"
)

type Color [3]float64

type Photon struct {
	vector geometry.Vector
	rgb    Color
}
