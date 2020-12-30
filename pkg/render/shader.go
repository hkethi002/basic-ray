package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
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

// func GetAngleFactor(triangle *geometry.Triangle, photon *Photon) float64 {
// 	if len(triangle.VertexNormals) == 0 {
// 		return math.Max(0, geometry.DotProduct(normalVector, geometry.Normalize(photon.vector))*-1)
// 	}
//
// 	return angleFactor
// }

func GetShadingNormal(vertexNormals []geometry.Vector, vertexes []geometry.Point, receivePoint geometry.Point) geometry.Vector {
	var u, v, w float64
	v0 := geometry.CreateVector(vertexes[1], vertexes[0])
	v1 := geometry.CreateVector(vertexes[2], vertexes[0])
	v2 := geometry.CreateVector(receivePoint, vertexes[0])

	d00 := geometry.DotProduct(v0, v0)
	d01 := geometry.DotProduct(v0, v1)
	d11 := geometry.DotProduct(v1, v1)
	d20 := geometry.DotProduct(v2, v0)
	d21 := geometry.DotProduct(v2, v1)

	denom := (d00 * d11) - (d01 * d01)
	v = (d11*d20 - d01*d21) / denom
	w = (d00*d21 - d01*d20) / denom
	u = 1.0 - v - w

	normal := geometry.Add(
		geometry.Add(
			geometry.ScalarProduct(vertexNormals[0], u),
			geometry.ScalarProduct(vertexNormals[1], v),
		),
		geometry.ScalarProduct(vertexNormals[2], w),
	)
	return geometry.Normalize(normal)
}

func GetShadingNormal2(vertexNormals []geometry.Vector, vertexes []geometry.Point, receivePoint geometry.Point) geometry.Vector {
	var subTriangle *geometry.Triangle

	subTriangle = &geometry.Triangle{Vertex0: vertexes[0], Vertex1: vertexes[1], Vertex2: receivePoint}
	ABParea := subTriangle.GetArea()

	subTriangle.Vertex0 = vertexes[2]
	CBParea := subTriangle.GetArea()

	subTriangle.Vertex1 = vertexes[0]
	CAParea := subTriangle.GetArea()
	total := ABParea + CBParea + CAParea

	normal := geometry.Add(
		geometry.Add(
			geometry.ScalarProduct(vertexNormals[0], 1-CBParea/total),
			geometry.ScalarProduct(vertexNormals[1], 1-CAParea/total),
		),
		geometry.ScalarProduct(vertexNormals[2], 1-ABParea/total),
	)
	normal = geometry.Normalize(normal)
	return normal
}

func GouraudShader(receiveVector geometry.Vector, receivePoint geometry.Point, photons []*Photon, triangle *geometry.Triangle) Photon {
	normalVector := GetShadingNormal(triangle.VertexNormals, []geometry.Point{triangle.Vertex0, triangle.Vertex1, triangle.Vertex2}, receivePoint)
	totalCollected := Photon{vector: receiveVector}
	for _, photon := range photons {
		angleFactor := math.Max(0, geometry.DotProduct(normalVector, geometry.Normalize(photon.vector))*-1)
		totalCollected.rgb[0] += photon.rgb[0] * angleFactor
		totalCollected.rgb[1] += photon.rgb[1] * angleFactor
		totalCollected.rgb[2] += photon.rgb[2] * angleFactor
	}
	pi := 2 * 3.1415926538 / 2
	totalCollected.rgb[0] *= triangle.DiffuseAlbedo[0] / pi
	totalCollected.rgb[1] *= triangle.DiffuseAlbedo[1] / pi
	totalCollected.rgb[2] *= triangle.DiffuseAlbedo[2] / pi
	return totalCollected
}

func DiffuseShader(receiveVector geometry.Vector, photons []*Photon, triangle *geometry.Triangle) Photon {
	normalVector := triangle.GetNormal() // Should already be normalized
	totalCollected := Photon{vector: receiveVector}
	for _, photon := range photons {
		angleFactor := math.Max(0, geometry.DotProduct(normalVector, geometry.Normalize(photon.vector))*-1)
		totalCollected.rgb[0] += photon.rgb[0] * angleFactor
		totalCollected.rgb[1] += photon.rgb[1] * angleFactor
		totalCollected.rgb[2] += photon.rgb[2] * angleFactor
	}
	pi := 2 * 3.1415926538 / 2
	totalCollected.rgb[0] *= triangle.DiffuseAlbedo[0] / pi
	totalCollected.rgb[1] *= triangle.DiffuseAlbedo[1] / pi
	totalCollected.rgb[2] *= triangle.DiffuseAlbedo[2] / pi
	totalReflected := Photon{
		rgb: Color{
			totalCollected.rgb[0],
			totalCollected.rgb[1],
			totalCollected.rgb[2],
		},
		vector: receiveVector,
	}
	return totalReflected
}

// Refraction for transparent
