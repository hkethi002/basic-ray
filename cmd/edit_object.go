package cmd

import (
	geometry "basic-ray/pkg/geometry"
	sceneIo "basic-ray/pkg/scene"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	var (
		magnitude float64
		output    string
	)

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit a mesh",
	}

	var translateCmd = &cobra.Command{
		Use:   "move [OBJECT] [X] [Y] [Z]",
		Short: "Translate an object along a vector",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			vector := geometry.Vector{}
			for i := 1; i < 4; i++ {
				component, err := strconv.ParseFloat(args[i], 64)
				check(err)
				vector[i-1] = component
			}

			if magnitude != 0 {
				vector = geometry.ScalarProduct(
					geometry.Normalize(vector), magnitude)
			}

			if output == "" {
				output = args[0]
			}

			translate(args[0], output, vector)
		},
	}

	translateCmd.Flags().Float64VarP(&magnitude, "distance", "d", 0, "Distance to move the object")
	translateCmd.Flags().StringVarP(&output, "output", "o", "", "output name, defaults to writing over input file")
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(translateCmd)
}

func translate(objectFile, output string, vector geometry.Vector) {
	object, err := sceneIo.ReadObject(objectFile)
	check(err)

	for i, vertex := range object.Vertexes {
		object.Vertexes[i] = geometry.Translate(vertex, vector)
	}

	err = sceneIo.WriteObject(output, object)
	check(err)
}
