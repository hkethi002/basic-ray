package scene

import (
	geometry "basic-ray/pkg/geometry"
	render "basic-ray/pkg/render"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	DIRECTIONAL = iota
	DELTA       = iota
)

type CameraSpec struct {
	Eye         geometry.Point `json:"eye"`
	BottomLeft  geometry.Point `json:"bottomLeft"`
	BottomRight geometry.Point `json:"bottomRight"`
	TopLeft     geometry.Point `json:"topLeft"`
	Resolution  [2]int         `json:"resolution"`
}

type LightSpec struct {
	Type      string          `json:"type"`
	Direction geometry.Vector `json:"direction,omitempty"`
	Location  geometry.Point  `json:"location, omitempty"`
	RGB       render.Color    `json:"rgb"`
}

type Scene struct {
	Assets       []string    `json:"assets"`
	Camera       CameraSpec  `json:"camera"`
	LightSources []LightSpec `json:"light_sources"`
}

func (scene *Scene) GetComponents() (*render.Camera, []*geometry.Object, []render.LightSource, error) {
	var err error
	camera := render.MakeCamera(
		scene.Camera.BottomLeft,
		scene.Camera.BottomRight,
		scene.Camera.TopLeft,
		scene.Camera.Resolution[0],
		scene.Camera.Resolution[1],
	)

	objects := make([]*geometry.Object, len(scene.Assets))
	for i := 0; i < len(scene.Assets); i++ {
		object, err := ReadObject(scene.Assets[i])
		if err != nil {
			fmt.Printf("Error loading asset %s\n: error:%s\n", scene.Assets[i], err)
		}
		objects[i] = object
	}

	lightSources := make([]render.LightSource, len(scene.LightSources))
	for i := 0; i < len(scene.LightSources); i++ {
		switch scene.LightSources[i].Type {
		case "directional":
			lightSources[i] = &render.DirectionalLight{
				Direction: scene.LightSources[i].Direction,
				RGB:       scene.LightSources[i].RGB,
			}
		case "delta":
			lightSources[i] = &render.DeltaLight{
				Location: scene.LightSources[i].Location,
				RGB:      scene.LightSources[i].RGB,
			}
		}

	}
	return camera, objects, lightSources, err
}

func LoadScene(filename string) (*Scene, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// json data
	var scene Scene

	// unmarshall it
	err = json.Unmarshal(data, &scene)
	if err != nil {
		return nil, err
	}

	return &scene, nil
}
