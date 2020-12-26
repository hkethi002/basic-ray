package render

import (
	geometry "basic-ray/pkg/geometry"
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

func DiffuesShader(receiveVector geometry.Vector, lightVectors []Photon, triangle *geometry.Triangle) Photon {
	_ = triangle.GetNormal()
	totalReflected := Photon{vector: receiveVector}
	for _, lightVector := range lightVectors {
		totalReflected.rgb[0] += lightVector.rgb[0] * triangle.DiffuseAlbedo[0]
		totalReflected.rgb[1] += lightVector.rgb[1] * triangle.DiffuseAlbedo[1]
		totalReflected.rgb[2] += lightVector.rgb[2] * triangle.DiffuseAlbedo[2]
	}
	return totalReflected
}

// Refraction for transparent
