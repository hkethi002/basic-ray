package cmd

import (
	geometry "basic-ray/pkg/geometry"
	mesh "basic-ray/pkg/mesh"
	sceneIo "basic-ray/pkg/scene"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	var (
		refinement int
		shading    bool
		output     string
	)

	var createSphereCmd = &cobra.Command{
		Use:   "sphere [X] [Y] [Z] [RADIUS]",
		Short: "Create a sphere mesh",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			origin := geometry.Point{}
			for i := 0; i < 3; i++ {
				coord, err := strconv.ParseFloat(args[i], 64)
				check(err)
				origin[i] = coord
			}

			radius, err := strconv.ParseFloat(args[3], 64)
			check(err)
			CreateSphereMesh(origin, radius, refinement, shading, output)
		},
	}

	createSphereCmd.Flags().IntVarP(&refinement, "refinement", "r", 0, "Iterations to refine mesh, increases # of triangles by 4")
	createSphereCmd.Flags().BoolVar(&shading, "gouraud", true, "Use Gouraud Shading instead of Flat")
	createSphereCmd.Flags().StringVarP(&output, "output", "o", "Sphere", "Sphere name")
	rootCmd.AddCommand(createSphereCmd)
}

func CreateSphereMesh(origin geometry.Point, radius float64, refinement int, shading bool, output string) {
	sphere := mesh.Sphere{
		Radius: radius,
		Origin: origin,
	}
	sphereMesh := sphere.CreateMesh(refinement, shading)
	textures := []geometry.TextureProperties{
		geometry.TextureProperties{
			DiffuseAlbedo: [3]float64{0.18, 0, 0.18},
			MaterialType:  3,
		},
	}
	object := mesh.CreateObject(sphereMesh, textures, make([]int, 0))

	filename := fmt.Sprintf("%s.json", output)
	sceneIo.WriteObject(filename, object)
}
