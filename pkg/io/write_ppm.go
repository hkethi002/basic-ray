package io

import (
	"bufio"
	"fmt"
	"os"
)

type Image [][][3]int

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Write(image Image, filePath string) {
	fileObject, err := os.Create(filePath)
	check(err)
	defer fileObject.Close()

	writer := bufio.NewWriter(fileObject)

	writer.WriteString("P3\n")
	writer.WriteString(fmt.Sprintf("%d %d\n", len(image[0]), len(image)))
	writer.WriteString("255\n")
	for _, i := range image {
		for _, j := range i {
			writer.WriteString(fmt.Sprintf("%d %d %d ", j[0], j[1], j[2]))
		}
		writer.WriteString("\n")
	}
	writer.Flush()
}

func main() {
	image := [][][3]int{
		[][3]int{
			[3]int{255, 0, 0},
			[3]int{255, 0, 0},
			[3]int{255, 0, 0},
		},
		[][3]int{
			[3]int{0, 255, 0},
			[3]int{0, 255, 0},
			[3]int{0, 255, 0},
		},
		[][3]int{
			[3]int{0, 0, 255},
			[3]int{0, 0, 255},
			[3]int{0, 0, 255},
		},
		[][3]int{
			[3]int{255, 0, 255},
			[3]int{90, 90, 25},
			[3]int{80, 10, 55},
		},
	}
	Write(image, "image.ppm")
}
