package physics

import (
	geometry "basic-ray/pkg/geometry"
)

type PhysicObject interface {
	Simulate(ticks int)
}

type MoveableMesh interface {
	SetLocation(location geometry.Point)
}

type SimpleObject struct {
	Location geometry.Point
	Velocity geometry.Vector // units per tick
	Mesh     MoveableMesh
}

func (simpleObject *SimpleObject) Simulate(ticks int) {
	simpleObject.Location = geometry.Translate(
		simpleObject.Location,
		geometry.ScalarProduct(simpleObject.Velocity, float64(ticks)),
	)
	simpleObject.Mesh.SetLocation(simpleObject.Location)
}
