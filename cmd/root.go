package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "basicray",
		Short: "Basic Ray tracing render",
		Long: `A Slow and Rigid Ray tracigin 3D renderer built by Harsha in Go.
	Complete documentation is available at https://github.com/hkethi002/basic-ray`,
		Version: "0.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
