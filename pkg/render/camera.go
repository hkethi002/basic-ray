package render

import (
	geometry "basic-ray/pkg/geometry"
	"math"
)

type Camera interface {
	GetRays(i, j, samples int) []*geometry.Ray
	GetPixels() *[][]Color
	SetPixel(i, j int, color Color)
	GetDisplayPixels() []uint8
	GetViewPlane() ViewPlane
}

type BaseCamera struct {
	ViewPlane ViewPlane
	Pixels    *[][]Color
	Display   bool
	DisplayPixels []uint8
}

type OrthoCamera struct {
	Center geometry.Point
	BaseCamera
	Direction geometry.Vector
}

type PerspectiveCamera struct {
	Center geometry.Point
	Eye    geometry.Point
	BaseCamera
}

type PinholeCamera struct {
	Eye                 geometry.Point
	Zoom                float64
	LookPoint           geometry.Point
	UpVector            geometry.Vector
	DistanceToViewPlane float64
	ExposureTime        float64
	u                   geometry.Vector
	v                   geometry.Vector
	w                   geometry.Vector
	BaseCamera
}

type ThinLensCamera struct {
	Eye                 geometry.Point
	Zoom                float64
	LookPoint           geometry.Point
	UpVector            geometry.Vector
	DistanceToViewPlane float64
	ExposureTime        float64
	u                   geometry.Vector
	v                   geometry.Vector
	w                   geometry.Vector
	Sampler             Sampler
	LensRadius          float64
	FocalDistance       float64
	BaseCamera
}

type ViewPlane struct {
	HorizontalResolution int
	VerticalResolution   int
	PixelSize            float64
	Gamma                float64
	InverseGamma         float64
	Sampler              Sampler
}

func (camera *BaseCamera) GetViewPlane() ViewPlane {
	return camera.ViewPlane
}

func (camera *BaseCamera) GetPixels() *[][]Color {
	return camera.Pixels
}

func (camera *BaseCamera) GetDisplayPixels() []uint8 {
	return camera.DisplayPixels
}

func (camera *BaseCamera) SetPixel(i, j int, color Color) {
	if camera.ViewPlane.Gamma != 1 {
		color = Pow(color, camera.ViewPlane.InverseGamma)
	}
	(*camera.Pixels)[i][j] = color
	if camera.Display {
		offset := ((camera.ViewPlane.HorizontalResolution - i - 1) + (camera.ViewPlane.HorizontalResolution * j) )* 4
		(camera.DisplayPixels)[offset] = uint8(math.Min(255, color[0] * 25))
		(camera.DisplayPixels)[offset+1] = uint8(math.Min(255, color[1] * 25))
		(camera.DisplayPixels)[offset+2] = uint8(math.Min(255, color[2] * 25))
		(camera.DisplayPixels)[offset+3] = 255
	}
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

func (camera *PerspectiveCamera) GetRays(i, j, samples int) []*geometry.Ray {
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

func (camera *PerspectiveCamera) GetRay(i, j int, jitter geometry.Point2D) *geometry.Ray {
	origin := geometry.Point{
		camera.Center[0] + camera.ViewPlane.PixelSize*((float64)(i)-(0.5*(float64)(camera.ViewPlane.HorizontalResolution-1))+jitter[0]),
		camera.Center[1] + camera.ViewPlane.PixelSize*((float64)(j)-(0.5*(float64)(camera.ViewPlane.VerticalResolution-1))+jitter[1]),
		camera.Center[2],
	}
	direction := geometry.Normalize(geometry.CreateVector(origin, camera.Eye))
	return &geometry.Ray{Origin: origin, Vector: direction}
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
		camera.Center[0] + camera.ViewPlane.PixelSize*((float64)(i)-(0.5*(float64)(camera.ViewPlane.HorizontalResolution-1))+jitter[0]),
		camera.Center[1] + camera.ViewPlane.PixelSize*((float64)(j)-(0.5*(float64)(camera.ViewPlane.VerticalResolution-1))+jitter[1]),
		camera.Center[2],
	}
	return &geometry.Ray{Origin: origin, Vector: geometry.Vector{camera.Direction[0], camera.Direction[1], camera.Direction[2]}}
}

func (camera *PinholeCamera) Initialize() {
	if (camera.UpVector == geometry.Vector{0, 0, 0}) {
		camera.UpVector[1] = 1
	}
	if camera.ExposureTime == 0 {
		camera.ExposureTime = 1.0
	}
	if camera.Zoom == 0 {
		camera.Zoom = 1.0
	}
	camera.ViewPlane.PixelSize = camera.ViewPlane.PixelSize / camera.Zoom
	camera.w = geometry.Normalize(geometry.CreateVector(camera.Eye, camera.LookPoint))
	camera.u = geometry.Normalize(geometry.CrossProduct(camera.UpVector, camera.w))
	camera.v = geometry.CrossProduct(camera.w, camera.u)
	if camera.Display {
		camera.DisplayPixels = make([]uint8, 4 * camera.ViewPlane.HorizontalResolution * camera.ViewPlane.VerticalResolution)
	}
}

func (camera *PinholeCamera) GetRay(i, j int, jitter geometry.Point2D) *geometry.Ray {
	x := camera.ViewPlane.PixelSize * ((float64)(i) - (0.5*float64(camera.ViewPlane.HorizontalResolution) + jitter[0]))
	y := camera.ViewPlane.PixelSize * ((float64)(j) - (0.5*float64(camera.ViewPlane.VerticalResolution) + jitter[1]))

	vector := geometry.Subtract(
		geometry.Add(
			geometry.ScalarProduct(camera.u, x),
			geometry.ScalarProduct(camera.v, y),
		),
		geometry.ScalarProduct(camera.w, camera.DistanceToViewPlane),
	)

	return &geometry.Ray{
		Origin: camera.Eye,
		Vector: geometry.Normalize(vector),
	}
}

func (camera *PinholeCamera) GetRays(i, j, samples int) []*geometry.Ray {
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

func (camera *ThinLensCamera) Initialize() {
	if (camera.UpVector == geometry.Vector{0, 0, 0}) {
		camera.UpVector[1] = 1
	}
	if camera.ExposureTime == 0 {
		camera.ExposureTime = 1.0
	}
	if camera.Zoom == 0 {
		camera.Zoom = 1.0
	}
	camera.ViewPlane.PixelSize = camera.ViewPlane.PixelSize / camera.Zoom
	camera.w = geometry.Normalize(geometry.CreateVector(camera.Eye, camera.LookPoint))
	camera.u = geometry.Normalize(geometry.CrossProduct(camera.UpVector, camera.w))
	camera.v = geometry.CrossProduct(camera.w, camera.u)
	if camera.Display {
		camera.DisplayPixels = make([]uint8, 4 * camera.ViewPlane.HorizontalResolution * camera.ViewPlane.VerticalResolution)
	}
}

func (camera *ThinLensCamera) GetRay(i, j int, jitter geometry.Point2D) *geometry.Ray {
	x := camera.ViewPlane.PixelSize * ((float64)(i) - (0.5*float64(camera.ViewPlane.HorizontalResolution) + jitter[0]))
	y := camera.ViewPlane.PixelSize * ((float64)(j) - (0.5*float64(camera.ViewPlane.VerticalResolution) + jitter[1]))

	lensSample := camera.Sampler.SampleUnitCircle()
	lensX := camera.LensRadius * lensSample[0]
	lensY := camera.LensRadius * lensSample[1]

	focalX := x * camera.FocalDistance / camera.DistanceToViewPlane
	focalY := y * camera.FocalDistance / camera.DistanceToViewPlane

	vector := geometry.Subtract(
		geometry.Add(
			geometry.ScalarProduct(camera.u, focalX-lensX),
			geometry.ScalarProduct(camera.v, focalY-lensY),
		),
		geometry.ScalarProduct(camera.w, camera.FocalDistance),
	)

	origin := geometry.Translate(
		camera.Eye,
		geometry.Add(
			geometry.ScalarProduct(camera.u, lensX),
			geometry.ScalarProduct(camera.v, lensY),
		),
	)

	return &geometry.Ray{
		Origin: origin,
		Vector: geometry.Normalize(vector),
	}
}

func (camera *ThinLensCamera) GetRays(i, j, samples int) []*geometry.Ray {
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
