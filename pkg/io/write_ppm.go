package io

import (
	render "basic-ray/pkg/render"
	"bufio"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Write(camera *render.Camera, filePath string) {
	image := *camera.Pixels
	fileObject, err := os.Create(filePath)
	check(err)
	defer fileObject.Close()

	writer := bufio.NewWriter(fileObject)

	writer.WriteString("P3\n")
	writer.WriteString(fmt.Sprintf("%d %d\n", len(image[0]), len(image)))
	writer.WriteString("255\n")
	for _, i := range image {
		for _, j := range i {
			writer.WriteString(fmt.Sprintf("%d %d %d ", int64(j[0]), int64(j[1]), int64(j[2])))
		}
		writer.WriteString("\n")
	}
	writer.Flush()
}
