package scene

import (
	render "basic-ray/pkg/render"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func WritePNG(camera render.Camera, filepath string) {
	height := camera.GetViewPlane().VerticalResolution
	width := camera.GetViewPlane().HorizontalResolution
	pixels := *camera.GetPixels()
	weight := 255.0

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(width-x-1, height-y-1, color.NRGBA{
				R: (uint8)(pixels[x][y][0] * weight),
				G: (uint8)(pixels[x][y][1] * weight),
				B: (uint8)(pixels[x][y][2] * weight),
				A: 255,
			})
		}
	}
	f, err := os.Create(filepath)
	check(err)

	if err := png.Encode(f, img); err != nil {
		f.Close()
		check(err)
	}

	if err := f.Close(); err != nil {
		check(err)
	}
}

func WriteImage(camera render.Camera, filePath string) {
	image := transpose(*camera.GetPixels())
	fmt.Println(image[0][0])
	fileObject, err := os.Create(filePath)
	check(err)
	defer fileObject.Close()

	writer := bufio.NewWriter(fileObject)

	writer.WriteString("P3\n")
	writer.WriteString(fmt.Sprintf("%d %d\n", len(image[0]), len(image)))

	writer.WriteString("255\n")
	weight := 255.0 // / findMax(image)

	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[i]); j++ {
			color := image[i][j]
			writer.WriteString(fmt.Sprintf("%d %d %d ", int64(math.Round(color[0]*weight)), int64(math.Round(color[1]*weight)), int64(math.Round(color[2]*weight))))
		}
		writer.WriteString("\n")
	}
	writer.Flush()
}

func transpose(slice [][]render.Color) [][]render.Color {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]render.Color, xl)
	for i := range result {
		result[i] = make([]render.Color, yl)
	}
	for i := 0; i < xl; i++ {
		for j := yl - 1; j >= 0; j-- {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func flip(image [][]render.Color) [][]render.Color {
	rows := len(image)
	result := make([][]render.Color, rows)
	for i, row := range image {
		result[rows-(i+1)] = row
	}
	return result
}

func findMaxVal(color render.Color) float64 {
	return math.Max(math.Max(color[0], color[1]), color[2])
}

func findMaxCol(col []render.Color) float64 {
	max := 0.0
	for _, c := range col {
		max = math.Max(max, findMaxVal(c))
	}
	return max
}

func findMax(image [][]render.Color) float64 {
	max := 0.0
	for _, col := range image {
		max = math.Max(max, findMaxCol(col))
	}
	return max
}
