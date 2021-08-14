package render

import (
	geometry "basic-ray/pkg/geometry"
	_ "math"
)

type Shader interface {
	BRDF(shadeRec *ShadeRec, incidentVector *geometry.Vector, outVector *geometry.Vector) Color
	SampleBRDF(shadeRec *ShadeRec, incidentVector *geometry.Vector, outVector *geometry.Vector) Color
	Rho(shadeRec *ShadeRec, outVector *geometry.Vector) Color
}

type LambertianShader struct {
	DiffuseReflectionCoefficient float64
	DiffuseColor                 Color
}

func (shader *LambertianShader) BRDF(shadeRec *ShadeRec, incidentVector *geometry.Vector, outVector *geometry.Vector) Color {
	return ScalarProduct(shader.DiffuseColor, shader.DiffuseReflectionCoefficient/PI)
}

func (shader *LambertianShader) Rho(shadeRec *ShadeRec, outVector *geometry.Vector) Color {
	return ScalarProduct(shader.DiffuseColor, shader.DiffuseReflectionCoefficient)
}

type Material interface {
	Shade(shadeRec *ShadeRec) Color
	// AreaLightShade(shadeRec *ShadeRec) Color
	// PathShade(shadeRec *ShadeRec) Color
}

type MatteMaterial struct {
	AmbientBRDF *LambertianShader
	DiffuseBRDF *LambertianShader
}

func (material *MatteMaterial) Shade(shadeRec *ShadeRec) Color {
	wo := geometry.ScalarProduct(shadeRec.Ray.Vector, -1)
	L := ElementwiseProduct(
		material.AmbientBRDF.Rho(shadeRec, &wo),
		shadeRec.World.AmbientLight.IncidentRadiance(shadeRec),
	)
	for _, light := range shadeRec.World.Lights {
		wi := light.GetDirection(shadeRec)
		incidentCos := geometry.DotProduct(shadeRec.Normal, wi)
		if incidentCos > 0.0 {
			L = Add(L, ScalarProduct(
				ElementwiseProduct(
					material.DiffuseBRDF.BRDF(shadeRec, &wo, &wi),
					light.IncidentRadiance(shadeRec),
				),
				incidentCos,
			))
		}
	}
	return L
}
