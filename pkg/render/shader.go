package render

import (
	geometry "basic-ray/pkg/geometry"
	_ "fmt"
)

func SpecularRatio(incedentVector, receiveVector, normalVector geometry.Vector, specularAlbedo [3]float64) [3]float64 {
	reflectionVector := geometry.Subtract(incedentVector, geometry.ScalarProduct(geometry.ScalarProduct(normalVector, geometry.DotProduct(incedentVector, normalVector)), 2))

	if geometry.Normalize(receiveVector) != geometry.Normalize(reflectionVector) {
		return [3]float64{}
	} else {
		return specularAlbedo
	}
}

func SpecularShader(lightVector, receiveVector geometry.Vector, triangle *geometry.Triangle) [3]float64 {
	normalVector := triangle.GetNormal()
	return SpecularRatio(lightVector, receiveVector, normalVector, triangle.SpecularAlbedo)
}

func DiffuseShader(receiveVector geometry.Vector, photons []*Photon, triangle *geometry.Triangle) Photon {
	normalVector := triangle.GetNormal() // Should already be normalized
	totalCollected := Photon{vector: receiveVector}
	// facingAngleFactor := geometry.DotProduct(normalVector, receiveVector)
	for _, photon := range photons {
		angleFactor := geometry.DotProduct(normalVector, photon.vector) * -1
		totalCollected.rgb[0] += photon.rgb[0] * angleFactor
		totalCollected.rgb[1] += photon.rgb[1] * angleFactor
		totalCollected.rgb[2] += photon.rgb[2] * angleFactor
	}
	pi := 3.1415926538
	totalCollected.rgb[0] *= triangle.DiffuseAlbedo[0] / pi
	totalCollected.rgb[1] *= triangle.DiffuseAlbedo[1] / pi
	totalCollected.rgb[2] *= triangle.DiffuseAlbedo[2] / pi
	totalReflected := Photon{
		rgb: Color{
			totalCollected.rgb[0],
			totalCollected.rgb[1],
			totalCollected.rgb[2],
		},
	}
	return totalReflected
}

// Refraction for transparent
