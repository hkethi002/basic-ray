package render

import (
	geometry "basic-ray/pkg/geometry"
)

type Camera interface {
	GetRay(i, j int, jitter geometry.Point2D) *geometry.Ray
	GetRays(i, j, samples int) []*geometry.Ray
	GetPixels() *[][]Color
	SetPixel(i, j int, color Color)
	GetViewPlane() ViewPlane
}

type OrthoCamera struct {
	ViewPlane ViewPlane
	Pixels    *[][]Color
	TopLeft   geometry.Point
	Direction geometry.Vector
}

type ViewPlane struct {
	HorizontalResolution int
	VerticalResolution   int
	PixelSize            float64
	Gamma                float64
	InverseGamma         float64
	Sampler              Sampler
}

func (camera *OrthoCamera) GetViewPlane() ViewPlane {
	return camera.ViewPlane
}

func (camera *OrthoCamera) GetRays(i, j, samples int) []*geometry.Ray {
	rays := make([]*geometry.Ray, samples)
	jitter := geometry.Point2D{0, 0}
	for s := 0; s < samples; s++ {
		if s > 0 {
			jitter = camera.ViewPlane.Sampler.SampleUnitSquare()
		}
		rays[s] = camera.GetRay(i, j, jitter)
	}

	return rays
}

// TODO: Figure out actual projection mapping
func (camera *OrthoCamera) GetRay(i, j int, jitter geometry.Point2D) *geometry.Ray {
	origin := geometry.Point{
		camera.TopLeft[0] + camera.ViewPlane.PixelSize*((float64)(i)-(0.5*(float64)(camera.ViewPlane.HorizontalResolution-1))+jitter[0]),
		camera.TopLeft[1] + camera.ViewPlane.PixelSize*((float64)(j)-(0.5*(float64)(camera.ViewPlane.VerticalResolution-1))+jitter[1]),
		camera.TopLeft[2],
	}
	return &geometry.Ray{Origin: origin, Vector: geometry.Vector{camera.Direction[0], camera.Direction[1], camera.Direction[2]}}
}

func (camera *OrthoCamera) GetPixels() *[][]Color {
	return camera.Pixels
}

func (camera *OrthoCamera) SetPixel(i, j int, color Color) {
	if camera.ViewPlane.Gamma != 1 {
		color = Pow(color, camera.ViewPlane.InverseGamma)
	}
	(*camera.Pixels)[i][j] = color
}

/*
func GetPoint(camera *Camera, pixelI int, pixelJ int) geometry.Point {
	v := geometry.Add(
		geometry.ScalarProduct(camera.unitX, float64(pixelI)),
		geometry.ScalarProduct(camera.unitY, float64(pixelJ)),
	)
	return geometry.Point{v[0] + camera.origin[0], v[1] + camera.origin[1], v[2] + camera.origin[2]}
}

func MakeCamera(bottomLeftCorner, bottomRightCorner, topLeftCorner geometry.Point, pixelWidth, pixelHeight int) *Camera {
	pixels := make([][]Color, pixelWidth)
	for i := range pixels {
		pixels[i] = make([]Color, pixelHeight)
	}
	return &Camera{
		origin: bottomLeftCorner,
		unitX:  geometry.ScalarProduct(geometry.CreateVector(bottomRightCorner, bottomLeftCorner), 1.0/float64(pixelWidth)),
		unitY:  geometry.ScalarProduct(geometry.CreateVector(topLeftCorner, bottomLeftCorner), 1.0/float64(pixelHeight)),
		Pixels: &pixels,
	}
}
*/
