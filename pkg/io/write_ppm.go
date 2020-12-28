package io

import (
	render "basic-ray/pkg/render"
	"bufio"
	"fmt"
	"math"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Write(camera *render.Camera, filePath string) {
	image := flip(transpose(*camera.Pixels))
	fmt.Println(image[0][0])
	fileObject, err := os.Create(filePath)
	check(err)
	defer fileObject.Close()

	writer := bufio.NewWriter(fileObject)

	writer.WriteString("P3\n")
	writer.WriteString(fmt.Sprintf("%d %d\n", len(image[0]), len(image)))

	writer.WriteString("255\n")
	weight := 255.0 / findMax(image)
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
