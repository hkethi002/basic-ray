package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
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

type GlossySpecular struct {
	SpecularReflectionCoefficient float64
	SpecularColor                 Color
	Exp                           float64
}

func (shader *GlossySpecular) BRDF(shadeRec *ShadeRec, incidentVector *geometry.Vector, outVector *geometry.Vector) Color {
	incidentCos := geometry.DotProduct(shadeRec.Normal, *incidentVector)
	reflectionVector := geometry.Subtract(geometry.ScalarProduct(shadeRec.Normal, 2*incidentCos), *incidentVector)
	outCos := geometry.DotProduct(reflectionVector, *outVector)

	if outCos > 0.0 {
		color := ScalarProduct(
			shader.SpecularColor,
			math.Pow(outCos, shader.Exp)*shader.SpecularReflectionCoefficient,
		)
		return color
	}
	return BLACK
}

func (shader *GlossySpecular) Rho(shadeRec *ShadeRec, outVector *geometry.Vector) Color {
	return BLACK
}

type Material interface {
	Shade(shadeRec *ShadeRec) Color
	AreaLightShade(shadeRec *ShadeRec) Color
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
			inShadow := false

			if light.CastsShadows() {
				shadowRay := geometry.Ray{Origin: shadeRec.HitPoint, Vector: wi}
				inShadow = light.InShadow(&shadowRay, shadeRec)
			}

			if !inShadow {
				L = Add(L, ScalarProduct(
					ElementwiseProduct(
						material.DiffuseBRDF.BRDF(shadeRec, &wi, &wo),
						light.IncidentRadiance(shadeRec),
					),
					incidentCos,
				))
			}
		}
	}
	return L
}

func (material *MatteMaterial) AreaLightShade(shadeRec *ShadeRec) Color {
	wo := geometry.ScalarProduct(shadeRec.Ray.Vector, -1)
	L := ElementwiseProduct(
		material.AmbientBRDF.Rho(shadeRec, &wo),
		shadeRec.World.AmbientLight.IncidentRadiance(shadeRec),
	)
	for _, light := range shadeRec.World.Lights {
		wi := light.GetDirection(shadeRec)
		incidentCos := geometry.DotProduct(shadeRec.Normal, wi)
		if incidentCos > 0.0 {
			inShadow := false

			if light.CastsShadows() {
				shadowRay := geometry.Ray{Origin: shadeRec.HitPoint, Vector: wi}
				inShadow = light.InShadow(&shadowRay, shadeRec)
			}

			if !inShadow {
				L = Add(L, ScalarProduct(
					ElementwiseProduct(
						material.DiffuseBRDF.BRDF(shadeRec, &wi, &wo),
						light.IncidentRadiance(shadeRec),
					),
					incidentCos*light.GeometricFactor(shadeRec)/light.PDF(shadeRec),
				))
			}
		}
	}
	return L
}

type PhongMaterial struct {
	AmbientBRDF *LambertianShader
	DiffuseBRDF *LambertianShader
	GlossyBRDF  *GlossySpecular
}

func (material *PhongMaterial) Shade(shadeRec *ShadeRec) Color {
	wo := geometry.ScalarProduct(shadeRec.Ray.Vector, -1)
	L := ElementwiseProduct(
		material.AmbientBRDF.Rho(shadeRec, &wo),
		shadeRec.World.AmbientLight.IncidentRadiance(shadeRec),
	)
	for _, light := range shadeRec.World.Lights {
		wi := light.GetDirection(shadeRec)
		incidentCos := geometry.DotProduct(shadeRec.Normal, wi)
		if incidentCos > 0.0 {
			inShadow := false

			if light.CastsShadows() {
				shadowRay := geometry.Ray{Origin: shadeRec.HitPoint, Vector: wi}
				inShadow = light.InShadow(&shadowRay, shadeRec)
			}

			if !inShadow {
				L = Add(L, ScalarProduct(
					ElementwiseProduct(
						Add(
							material.DiffuseBRDF.BRDF(shadeRec, &wi, &wo),
							material.GlossyBRDF.BRDF(shadeRec, &wi, &wo),
						),
						light.IncidentRadiance(shadeRec),
					),
					incidentCos,
				))
			}
		}
	}
	return L
}

func (material *PhongMaterial) AreaLightShade(shadeRec *ShadeRec) Color {
	wo := geometry.ScalarProduct(shadeRec.Ray.Vector, -1)
	L := ElementwiseProduct(
		material.AmbientBRDF.Rho(shadeRec, &wo),
		shadeRec.World.AmbientLight.IncidentRadiance(shadeRec),
	)
	for _, light := range shadeRec.World.Lights {
		wi := light.GetDirection(shadeRec)
		incidentCos := geometry.DotProduct(shadeRec.Normal, wi)
		if incidentCos > 0.0 {
			inShadow := false

			if light.CastsShadows() {
				shadowRay := geometry.Ray{Origin: shadeRec.HitPoint, Vector: wi}
				inShadow = light.InShadow(&shadowRay, shadeRec)
			}

			if !inShadow {
				L = Add(L, ScalarProduct(
					ElementwiseProduct(
						Add(
							material.DiffuseBRDF.BRDF(shadeRec, &wi, &wo),
							material.GlossyBRDF.BRDF(shadeRec, &wi, &wo),
						),
						light.IncidentRadiance(shadeRec),
					),
					incidentCos*light.GeometricFactor(shadeRec)/light.PDF(shadeRec),
				))
			}
		}
	}
	return L
}

type EmmisiveMaterial struct {
	RadianceScalingFactor float64
	Color                 Color
}

func (material *EmmisiveMaterial) AreaLightShade(shadeRec *ShadeRec) Color {
	outCos := geometry.DotProduct(geometry.ScalarProduct(shadeRec.Normal, -1), shadeRec.Ray.Vector)
	if outCos > 0.0 {
		return ScalarProduct(material.Color, material.RadianceScalingFactor)
	} else {
		return BLACK
	}
}

func (material *EmmisiveMaterial) Shade(shadeRec *ShadeRec) Color {
	return BLACK
}

func (material *EmmisiveMaterial) GetEmmittedRadiance(shadeRec *ShadeRec) Color {
	return ScalarProduct(material.Color, material.RadianceScalingFactor)
}
