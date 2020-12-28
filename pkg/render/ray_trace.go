package render

import (
	geometry "basic-ray/pkg/geometry"
	_ "fmt"
	pb "github.com/cheggaaa/pb/v3"
	"math"
)

func Main(origin geometry.Point, lightSources []LightSource, camera *Camera, triangles []*geometry.Triangle) {
	var light Photon
	bar := pb.StartNew(len(*camera.Pixels) * len((*camera.Pixels)[0]))

	for i, row := range *camera.Pixels {
		for j, _ := range row {
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			light = Trace(&ray, triangles, lightSources, 0)
			(*camera.Pixels)[i][j] = light.rgb // GetWeightedColor()
			bar.Increment()
		}
	}
	bar.Finish()
}

func MultiThreadedMain(origin geometry.Point, lightSources []LightSource, camera *Camera, triangles []*geometry.Triangle) {
	progress := make(chan bool, 20)
	totalCount := len(*camera.Pixels) * len((*camera.Pixels)[0])
	bar := pb.StartNew(totalCount)
	go renderPixels(origin, triangles, lightSources, camera, progress)
	reportProgress(bar, progress, totalCount)
	close(progress)
	bar.Finish()

}

func reportProgress(bar *pb.ProgressBar, progress <-chan bool, totalCount int) {
	for i := 0; i < totalCount; i++ {
		<-progress
		bar.Increment()
	}
}

func renderPixels(origin geometry.Point, triangles []*geometry.Triangle, lightSources []LightSource, camera *Camera, progress chan<- bool) {
	jobs := make(chan bool, 5)
	for i, row := range *camera.Pixels {
		for j, _ := range row {
			jobs <- true
			ray := geometry.Ray{
				Origin: origin,
				Vector: geometry.Normalize(geometry.CreateVector(GetPoint(camera, i, j), origin)),
			}
			go renderPixel(&ray, triangles, lightSources, camera, i, j, progress, jobs)
		}
	}
}

func renderPixel(ray *geometry.Ray, triangles []*geometry.Triangle, lightSources []LightSource, camera *Camera, i, j int, progress chan<- bool, jobs <-chan bool) {
	light := Trace(ray, triangles, lightSources, 0)
	(*camera.Pixels)[i][j] = light.rgb // GetWeightedColor()
	// complete job
	<-jobs
	// report progress
	progress <- true
}

func Trace(ray *geometry.Ray, triangles []*geometry.Triangle, lightSources []LightSource, depth int) Photon {
	receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	photon := Photon{vector: receiveVector}
	closestPoint := float64(math.Inf(1))
	if depth >= 3 {
		return photon
	}
	for _, triangle := range triangles {
		intersects := geometry.GetIntersection(ray, triangle)
		if intersects == nil {
			continue
		}
		collision := *intersects
		distance := geometry.Distance(collision, ray.Origin)

		if distance < closestPoint {
			closestPoint = distance
		} else {
			continue
		}

		photon = GetColor(ray, collision, triangle, lightSources, triangles, depth)
	}

	return photon
}

func GetColor(
	ray *geometry.Ray,
	reflectionPoint geometry.Point,
	triangle *geometry.Triangle,
	lightSources []LightSource,
	triangles []*geometry.Triangle,
	depth int,
) Photon {
	receiveVector := geometry.Normalize(geometry.ScalarProduct(ray.Vector, -1))
	switch triangle.MaterialType {
	case geometry.REFLECTIVE:
		reflectionRay := &geometry.Ray{Origin: reflectionPoint, Vector: GetReflectiveVector(ray.Vector, triangle)}
		return Trace(reflectionRay, triangles, lightSources, depth+1)
	case geometry.DIFFUSE:
		photons := GetDirectLight(reflectionPoint, triangles, lightSources)
		// sampleRays := geometry.MakeSampleRays(reflectionPoint, triangle.GetNormal(), 16)
		// photon := Photon{vector: receiveVector}
		// for _, sampleRay := range sampleRays {
		// 	photon = Trace(sampleRay, triangles, lightSources, depth+1)
		// 	distance := geometry.Distance(ray.Origin, reflectionPoint)
		// 	d2 := distance * distance
		// 	photon.rgb[0] = photon.rgb[0] / d2
		// 	photon.rgb[1] = photon.rgb[1] / d2
		// 	photon.rgb[2] = photon.rgb[2] / d2
		// 	photons = append(photons, &photon)
		// }

		return DiffuseShader(receiveVector, photons, triangle)

	}

	return Photon{vector: receiveVector}
}
